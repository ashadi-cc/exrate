package memory

import (
	"fmt"
	"sync"
)

type memoryStore struct {
	mu   sync.RWMutex
	item map[string]string
}

//Set implementing driver.Set
func (store *memoryStore) Set(key, value string) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	store.item[key] = value
	return nil
}

//Get implementing driver.Get
func (store *memoryStore) Get(key string) (string, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	val, ok := store.item[key]
	if !ok {
		return "", fmt.Errorf("key: %s not found", key)
	}
	return val, nil
}
