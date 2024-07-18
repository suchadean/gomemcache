package gomemcache

import (
	"fmt"
	"sync"
	"time"
)

// CacheInterface identifies cache operators
type CacheInterface interface {
	// GetValue get the value of the specified key if it exists. Otherwise an error is returned.
	GetValue(key string) ([]byte, error)
	// SetValue set the specified key value pair with the given time to live (ttl)
	// If the ttl is defined as 0, the key value pair is held indefinitely
	SetValue(key string, value []byte, ttl time.Duration)

	// KeyExists returns true if the specified key exists in cache otherwise false
	KeyExists(key string) bool

	// DeleteKey removes the entry regarding specified key from the cache
	// If the specified key does not exist DeleteKey is a no-op
	DeleteKey(key string)
}

// MemCache simple in-memory cache
// Stores value as byte slice in a map with string keys
type MemCache struct {
	lock sync.RWMutex
	data map[string][]byte
}

// New instantiates a new MemCache
func New() *MemCache {
	return &MemCache{
		data: make(map[string][]byte),
	}
}

// GetValue get the value of the specified key if it exists. Otherwise an error is returned.
func (m *MemCache) GetValue(key string) ([]byte, error) {
	// Acquire read lock to ensure concurrent retrieval safety
	m.lock.RLock()
	defer m.lock.RUnlock()
	// Retrieve value throw error if none exists
	value, ok := m.data[key]
	if !ok {
		return nil, fmt.Errorf("key %s does not exist", key)
	}

	return value, nil
}

// SetValue set the specified key value pair with the given time to live (ttl)
// If the ttl is defined as 0, the key value pair is held indefinitely
func (m *MemCache) SetValue(key string, value []byte, ttl time.Duration) {
	// Acquire write lock to ensure concurrent write safety
	m.lock.Lock()
	defer m.lock.Unlock()

	// Set value on key
	m.data[key] = value

	// If ttl is greater than zero, start a concurrent function to clean up after ttl duration
	// If ttl is smaller than or zero, keep key value pair in cache indefinitely
	if ttl > 0 {
		go func() {
			<-time.After(ttl)
			m.lock.Lock()
			defer m.lock.Unlock()
			delete(m.data, key)
		}()
	}
}

// KeyExists returns true if the specified key exists in cache otherwise false
func (m *MemCache) KeyExists(key string) bool {
	// Acquire read lock to ensure concurrent retrieval safety
	m.lock.RLock()
	defer m.lock.RUnlock()

	// Check if Key exists in map
	_, ok := m.data[key]

	return ok
}

// DeleteKey removes the entry regarding specified key from the cache
// If the specified key does not exist DeleteKey is a no-op
func (m *MemCache) DeleteKey(key string) {
	// Acquire write lock to ensure concurrent write safety
	m.lock.Lock()
	defer m.lock.Unlock()

	delete(m.data, key)
}
