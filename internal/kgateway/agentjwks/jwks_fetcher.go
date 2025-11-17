package agentjwks

import (
	"container/heap"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"sync"
	"time"

	"github.com/go-jose/go-jose/v4"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// JwksFetcher is used for fetching and periodic updates of jwks
// Internally fetched jwks are stored as key-values: jwks-url:jose.JSONWebKeySet
// When a jwks is updated, registered subscribers are sent the current state of
// jwks store.
type JwksFetcher struct {
	mu            sync.Mutex
	Cache         *jwksCache
	Client        JwksHttpClient
	KeysetSources map[string]*JwksSource
	Schedule      FetchingSchedule
	Subscribers   []chan map[string]string
}

type jwksCache struct {
	Jwks map[string]jose.JSONWebKeySet
}

type FetchingSchedule []fetchAt

//go:generate go tool mockgen -destination mocks/mock_jwks_http_client.go -package mocks -source ./jwks_fetcher.go
type JwksHttpClient interface {
	FetchJwks(ctx context.Context, jwksURL string) (jose.JSONWebKeySet, error)
}

type JwksSource struct {
	JwksURL string
	Ttl     time.Duration
	Deleted bool
}

type fetchAt struct {
	at           time.Time
	keysetSource *JwksSource
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
		KeysetSources: make(map[string]*JwksSource),
		Schedule:      make([]fetchAt, 0),
		Subscribers:   make([]chan map[string]string, 0),
	}
	heap.Init(&toret.Schedule)

	return toret
}

// this function must be called when holding the `mu` lock
func (c *jwksCache) toJson() (map[string]string, error) {
	toret := make(map[string]string)
	for k, v := range c.Jwks {
		serializedJwks, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		toret[k] = string(serializedJwks)
	}
	return toret, nil
}

func LoadJwksfromJson(storedJwks map[string]string) (map[string]jose.JSONWebKeySet, error) {
	toret := make(map[string]jose.JSONWebKeySet)
	var errs []error

	for url, serializedJwks := range storedJwks {
		jwks := jose.JSONWebKeySet{}
		if err := json.Unmarshal([]byte(serializedJwks), &jwks); err != nil {
			errs = append(errs, err)
		}
		toret[url] = jwks
	}

	return toret, errors.Join(errs...)
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
	return &s[0]
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

		fetch := heap.Pop(&f.Schedule).(fetchAt)
		if fetch.keysetSource.Deleted {
			continue
		}
		jwks, err := f.Client.FetchJwks(ctx, fetch.keysetSource.JwksURL)
		if err != nil {
			log.Error(err, "error fetching jwks from ", fetch.keysetSource.JwksURL)
			if fetch.retryAttempt < 5 { // backoff by 5s * retry attempt number
				heap.Push(&f.Schedule, fetchAt{at: now.Add(time.Duration(5*(fetch.retryAttempt+1)) * time.Second), keysetSource: fetch.keysetSource, retryAttempt: fetch.retryAttempt + 1})
			} else {
				// give up retrying and schedule an update at a later time
				heap.Push(&f.Schedule, fetchAt{at: now.Add(fetch.keysetSource.Ttl), keysetSource: fetch.keysetSource})
			}
			continue
		}

		if !reflect.DeepEqual(jwks, f.Cache.Jwks[fetch.keysetSource.JwksURL]) {
			f.Cache.Jwks[fetch.keysetSource.JwksURL] = jwks
			haveUpdates = true
		}

		heap.Push(&f.Schedule, fetchAt{at: now.Add(fetch.keysetSource.Ttl), keysetSource: fetch.keysetSource})
	}

	if haveUpdates {
		update, err := f.Cache.toJson()
		if err != nil {
			log.Error(err, "error serializing jwks store")
			return
		}
		for _, s := range f.Subscribers {
			s <- update
		}
	}
}

func (f *JwksFetcher) UpdateJwksSources(ctx context.Context, updates []JwksSource) error {
	log := log.FromContext(ctx)

	var errs []error
	maybeUpdates := make(map[string]JwksSource)
	for _, s := range updates {
		maybeUpdates[s.JwksURL] = s
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	todelete := make([]string, 0)
	for s := range f.KeysetSources {
		if _, ok := maybeUpdates[s]; !ok {
			todelete = append(todelete, s)
		}
	}

	for _, s := range updates {
		if _, ok := f.KeysetSources[s.JwksURL]; !ok {
			if err := f.addKeyset(s.JwksURL, s.Ttl); err != nil {
				errs = append(errs, err)
			}
			continue
		}
		if *f.KeysetSources[s.JwksURL] != s {
			if err := f.updateKeyset(s.JwksURL, s.Ttl); err != nil {
				errs = append(errs, err)
			}
		}
	}

	haveUpdates := false
	for _, s := range todelete {
		haveUpdates = haveUpdates || f.removeKeyset(s)
	}

	if haveUpdates {
		update, err := f.Cache.toJson()
		if err != nil {
			log.Error(err, "error serializing jwks store")
			errs = append(errs, err)
			return errors.Join(errs...)
		}
		for _, s := range f.Subscribers {
			s <- update
		}
	}

	return errors.Join(errs...)
}

func (f *JwksFetcher) addKeyset(jwksUrl string, ttl time.Duration) error {
	if _, err := url.Parse(jwksUrl); err != nil {
		return fmt.Errorf("error parsing jwks url %w", err)
	}

	keysetSource := &JwksSource{JwksURL: jwksUrl, Ttl: ttl, Deleted: false}
	f.KeysetSources[jwksUrl] = keysetSource
	heap.Push(&f.Schedule, fetchAt{at: time.Now(), keysetSource: keysetSource}) // schedule an immediate fetch

	return nil
}

func (f *JwksFetcher) removeKeyset(jwksUrl string) bool {
	if keysetSource, ok := f.KeysetSources[jwksUrl]; ok {
		delete(f.KeysetSources, jwksUrl)
		delete(f.Cache.Jwks, jwksUrl)
		keysetSource.Deleted = true
		return true
	}
	return false
}

func (f *JwksFetcher) updateKeyset(jwksUrl string, ttl time.Duration) error {
	f.removeKeyset(jwksUrl)
	return f.addKeyset(jwksUrl, ttl)
}

func (f *JwksFetcher) SubscribeToUpdates() chan map[string]string {
	f.mu.Lock()
	defer f.mu.Unlock()

	subscriber := make(chan map[string]string)
	f.Subscribers = append(f.Subscribers, subscriber)

	return subscriber
}

func (c *jwksHttpClientImpl) FetchJwks(ctx context.Context, jwksURL string) (jose.JSONWebKeySet, error) {
	log := log.FromContext(ctx)
	log.Info("fetching jwks", "url", jwksURL)

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

	if response.StatusCode != 200 {
		return jose.JSONWebKeySet{}, fmt.Errorf("unexpected status code from jwks endpoint at %s: %d", jwksURL, response.StatusCode)
	}

	var jwks jose.JSONWebKeySet
	if err := json.NewDecoder(response.Body).Decode(&jwks); err != nil {
		return jose.JSONWebKeySet{}, fmt.Errorf("could not decode jwks: %w", err)
	}

	return jwks, nil
}
