package store

import (
	"fmt"
	"sync"
	"xrate/common/store/driver"
)

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]driver.Driver)
)

//Register register the driver
func Register(name string, driver driver.Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if driver == nil {
		panic("storage: Register driver is nil")
	}
	drivers[name] = driver
}

//Open open the driver by given name
func Open(name string) (driver.Driver, error) {
	driversMu.RLock()
	driveri, ok := drivers[name]
	driversMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("unknown driver %q (forgotten import?)", name)
	}
	return driveri, nil
}
