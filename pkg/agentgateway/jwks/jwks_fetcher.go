package jwks

import (
	"container/heap"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"sync"
	"time"

	"github.com/go-jose/go-jose/v4"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type jwksCache struct {
	Jwks map[string]jose.JSONWebKeySet
}

type FetchingSchedule []fetchAt

type JwksFetcher struct {
	mu            sync.Mutex
	Cache         *jwksCache
	Client        JwksHttpClient
	KeysetSources map[string]*keysetSource
	Schedule      FetchingSchedule
	Subscribers   []chan string
}

//go:generate go tool mockgen -destination mocks/mock_jwks_http_client.go -package mocks -source ./jwks_fetcher.go
type JwksHttpClient interface {
	FetchJwks(ctx context.Context, jwksURL string) (jose.JSONWebKeySet, error)
}

type keysetSource struct {
	JwksURL string
	Ttl     time.Duration
	Deleted bool
}

type fetchAt struct {
	at           time.Time
	keysetSource *keysetSource
	retryAttempt int
}

type jwksHttpClientImpl struct {
	Client *http.Client
}

func NewJwksCache() *jwksCache {
	return &jwksCache{
		Jwks: make(map[string]jose.JSONWebKeySet),
	}
}

func NewJwksFetcher(cache *jwksCache) *JwksFetcher {
	toret := &JwksFetcher{
		Cache: cache,
		Client: &jwksHttpClientImpl{
			Client: &http.Client{Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			}}},
		KeysetSources: make(map[string]*keysetSource),
		Schedule:      make([]fetchAt, 0),
		Subscribers:   make([]chan string, 0),
	}
	heap.Init(&toret.Schedule)

	return toret
}

func (c *jwksCache) toJson() ([]byte, error) {
	return json.Marshal(c.Jwks)
}

func (c *jwksCache) LoadfromJson(storedJwks string) error {
	stored := make(map[string]jose.JSONWebKeySet)

	if storedJwks == "" {
		return nil
	}

	if err := json.Unmarshal([]byte(storedJwks), &stored); err != nil {
		return err
	}
	c.Jwks = stored

	return nil
}

// heap implementation
func (s FetchingSchedule) Len() int           { return len(s) }
func (s FetchingSchedule) Less(i, j int) bool { return s[i].at.Before(s[j].at) }
func (s FetchingSchedule) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s *FetchingSchedule) Push(x any) {
	*s = append(*s, x.(fetchAt))
}
func (s *FetchingSchedule) Pop() any {
	old := *s
	n := len(old)
	x := old[n-1]
	*s = old[0 : n-1]
	return x
}
func (s FetchingSchedule) Peek() *fetchAt {
	if len(s) == 0 {
		return nil
	}
	return &s[len(s)-1]
}

func (f *JwksFetcher) Run(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			f.maybeFetchJwks(ctx)
		}
	}
}

func (f *JwksFetcher) maybeFetchJwks(ctx context.Context) {
	log := log.FromContext(ctx)
	haveUpdates := false

	f.mu.Lock()
	defer f.mu.Unlock()

	now := time.Now()
	for {
		maybeFetch := f.Schedule.Peek()
		if maybeFetch == nil || maybeFetch.at.After(now) {
			break
		}

		fetch := f.Schedule.Pop().(fetchAt)
		if fetch.keysetSource.Deleted {
			continue
		}
		jwks, err := f.Client.FetchJwks(ctx, fetch.keysetSource.JwksURL)
		if err != nil {
			log.Error(err, "error fetching jwks from ", fetch.keysetSource.JwksURL)
			if fetch.retryAttempt < 5 { // backoff by 5s * retry attempt number
				f.Schedule.Push(fetchAt{at: now.Add(time.Duration(5*(fetch.retryAttempt+1)) * time.Second), keysetSource: fetch.keysetSource, retryAttempt: fetch.retryAttempt + 1})
			} else {
				// give up retrying and schedule an update at a later time
				f.Schedule.Push(fetchAt{at: now.Add(fetch.keysetSource.Ttl), keysetSource: fetch.keysetSource})
			}
			continue
		}

		if !reflect.DeepEqual(jwks, f.Cache.Jwks[fetch.keysetSource.JwksURL]) {
			f.Cache.Jwks[fetch.keysetSource.JwksURL] = jwks
			haveUpdates = true
		}

		f.Schedule.Push(fetchAt{at: now.Add(fetch.keysetSource.Ttl), keysetSource: fetch.keysetSource})
	}

	if haveUpdates {
		update, err := f.Cache.toJson()
		if err != nil {
			log.Error(err, "error serializing jwks store")
			return
		}
		for _, s := range f.Subscribers {
			s <- string(update)
		}
	}
}

func (f *JwksFetcher) AddKeyset(jwksUrl string, ttl time.Duration) error {
	if _, err := url.Parse(jwksUrl); err != nil {
		return fmt.Errorf("error parsing jwks url %w", err)
	}

	keysetSource := &keysetSource{JwksURL: jwksUrl, Ttl: ttl, Deleted: false}

	f.mu.Lock()
	defer f.mu.Unlock()

	f.KeysetSources[jwksUrl] = keysetSource
	heap.Push(&f.Schedule, fetchAt{at: time.Now(), keysetSource: keysetSource}) // schedule an immediate fetch

	return nil
}

func (f *JwksFetcher) RemoveKeyset(jwksUrl string) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if keysetSource, ok := f.KeysetSources[jwksUrl]; ok {
		delete(f.KeysetSources, jwksUrl)
		delete(f.Cache.Jwks, jwksUrl)
		keysetSource.Deleted = true
		// TODO signal updates
	}
}

func (f *JwksFetcher) SubscribeToUpdates() chan string {
	f.mu.Lock()
	defer f.mu.Unlock()

	subscriber := make(chan string)
	f.Subscribers = append(f.Subscribers, subscriber)

	return subscriber
}

func (c *jwksHttpClientImpl) FetchJwks(ctx context.Context, jwksURL string) (jose.JSONWebKeySet, error) {
	log := log.FromContext(ctx)
	log.Info("fetching jwks", jwksURL)

	request, err := http.NewRequest(http.MethodGet, jwksURL, nil)
	if err != nil {
		return jose.JSONWebKeySet{}, fmt.Errorf("could not build request to get JWKS: %w", err)
	}

	// TODO (dmitri-d) control the size here maybe?
	response, err := c.Client.Do(request)
	if err != nil {
		return jose.JSONWebKeySet{}, err
	}
	defer response.Body.Close() //nolint:errcheck

	var jwks jose.JSONWebKeySet
	if err := json.NewDecoder(response.Body).Decode(&jwks); err != nil {
		return jose.JSONWebKeySet{}, fmt.Errorf("could not decode jwks: %w", err)
	}

	return jwks, nil
}
