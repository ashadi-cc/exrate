package crate

import (
	"fmt"
	"sync"
	"xrate/common/crate/provider"
)

var (
	clientMu sync.RWMutex
	clients  = make(map[string]provider.Client)
)

//Register register the client
func Register(name string, client provider.Client) {
	clientMu.Lock()
	defer clientMu.Unlock()
	if client == nil {
		panic("rate: Register driver is nil")
	}
	clients[name] = client
}

//Open open the client by given name
func Open(name string) (provider.Client, error) {
	clientMu.RLock()
	clienti, ok := clients[name]
	clientMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("unknown client %q (forgotten import?)", name)
	}
	return clienti, nil
}
