package memory

import (
	"sync"
	"xrate/common/store"
)

func init() {
	mStore := memoryStore{
		mu:   sync.RWMutex{},
		item: make(map[string]interface{}),
	}

	store.Register("memory", &mStore)
}
