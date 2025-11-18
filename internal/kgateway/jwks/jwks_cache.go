package jwks

import (
	"encoding/json"
	"hash/fnv"
	"maps"

	"github.com/go-jose/go-jose/v4"
)

const MAX_JWKS_STORE_SIZE = 35 * 1024 // 1024*1024 + 400*1024 // 1.4MiB

// jwks cache is an ordered collection of stores, each containing jwks whose
// total size doesn't exceed MAX_JWKS_STORE_SIZE
// a consistent jwks uri hash (see `UriHash()`) is used to determine the store
// for a given jwks
type jwksCache struct {
	stores []jwksStore
}

type jwksStore struct {
	jwks map[string]string // jwks uri -> jwks
	size int               // this is an approximate size, see `jwksStore.Size()``
}

func NewJwksCache() *jwksCache {
	c := jwksCache{}
	c.stores = append(c.stores, newJwksStore())
	return &c
}

func newJwksStore() jwksStore {
	return jwksStore{
		jwks: make(map[string]string),
	}
}

// Re-create jwks cache from state persisted in ConfigMaps
func (c *jwksCache) LoadJwksFromStores(storedJwks []map[string]string) error {
	newCache := NewJwksCache()

	for _, store := range storedJwks {
		for uri, serializedJwks := range store {
			jwks := jose.JSONWebKeySet{}
			if err := json.Unmarshal([]byte(serializedJwks), &jwks); err != nil {
				return err
			}
			newCache.compareAndAddJwks(uri, jwks)
		}
	}

	c.stores = newCache.stores
	return nil
}

// Add a jwks to cache. If an exact same jwks is already present in the cache, the result is a nop.
// If the store size exceeds `MAX_JWKS_STORE_SIZE` as a result of adding a jwks, a new store is added
// and all jwks are re-balanced between all stores.
func (c *jwksCache) compareAndAddJwks(uri string, jwks jose.JSONWebKeySet) (bool, error) {
	serializedJwks, err := json.Marshal(jwks)
	if err != nil {
		return false, err
	}

	idx := 0
	if l := len(c.stores); l > 1 {
		idx = int(UriHash(uri) % uint64(l))
	}

	if j, ok := c.stores[idx].jwks[uri]; ok {
		if j == string(serializedJwks) {
			return false, nil
		}
	}

	c.stores[idx].jwks[uri] = string(serializedJwks)
	c.stores[idx].size += len(uri) + len(c.stores[idx].jwks[uri])

	if c.stores[idx].serializedStoreSize() > MAX_JWKS_STORE_SIZE {
		c.addStore()
	}

	return true, nil
}

// Remove jwks from cache. If store size is reduced to zero as a result of jwks removal
// the store is deleted, and all jwks are re-balanced between remaining stores.
// If the last store size is reduced to zero, that store is not deleted.
// TODO (dmitri-d maybe we should delete it?)
func (c *jwksCache) deleteJwks(uri string) {
	idx := 0
	if l := len(c.stores); l > 1 {
		idx = int(UriHash(uri) % uint64(l))
	}

	if jwks, ok := c.stores[idx].jwks[uri]; ok {
		delete(c.stores[idx].jwks, uri)
		c.stores[idx].size -= (len(uri) + len(jwks))
	}

	if c.stores[idx].size == 0 {
		c.deleteStore()
	}
}

func (c *jwksCache) toJson() []map[string]string {
	copy := make([]map[string]string, len(c.stores))
	for i, store := range c.stores {
		copy[i] = maps.Clone(store.jwks)
	}
	return copy
}

// add another store, rebalance keysets
func (c *jwksCache) addStore() {
	newCache := jwksCache{
		stores: make([]jwksStore, len(c.stores)+1),
	}
	for i := range len(c.stores) + 1 {
		newCache.stores[i] = newJwksStore()
	}

	for _, store := range c.stores {
		for k, v := range store.jwks {
			newCache.copyJwks(k, v)
		}
	}

	c.stores = newCache.stores
}

// delete an empty store, rebalance keysets
func (c *jwksCache) deleteStore() {
	if len(c.stores) == 1 {
		return
	}

	newCache := jwksCache{
		stores: make([]jwksStore, len(c.stores)-1),
	}

	for i := range len(c.stores) - 1 {
		newCache.stores[i] = newJwksStore()
	}

	for _, store := range c.stores {
		for k, v := range store.jwks {
			newCache.copyJwks(k, v)
		}
	}

	c.stores = newCache.stores
}

func (c *jwksCache) copyJwks(uri string, jwks string) {
	idx := 0
	if l := len(c.stores); l > 1 {
		idx = int(UriHash(uri) % uint64(l))
	}
	c.stores[idx].jwks[uri] = string(jwks)
	c.stores[idx].size += len(uri) + len(c.stores[idx].jwks[uri])
}

// returns the size of serialized store (as it's stored in a ConfigMap)
// weird formula is based on comparison of differences in sizes between
// the internal store representation and store state persisted in a ConfigMap
func (s jwksStore) serializedStoreSize() int {
	return s.size + 7*len(s.jwks) + 2
}

func UriHash(uri string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(uri))
	return h.Sum64()
}
