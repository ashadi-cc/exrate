package helper

import (
	"xrate/common/store"
	"xrate/common/store/driver"
	"xrate/config"
)

//CurrentStore returns selected store from os.environment
func CurrentStore() (driver.Driver, error) {
	store, err := store.Open(config.GetStore().Driver)
	if err != nil {
		return nil, err
	}

	return store, nil
}
