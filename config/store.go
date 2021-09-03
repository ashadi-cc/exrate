package config

import "os"

type Store struct {
	Driver string
}

func GetStore() Store {
	return Store{
		Driver: os.Getenv("STORAGE_DRIVER"),
	}
}
