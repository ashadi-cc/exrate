package config

import "os"

type Client struct {
	ApiKey      string
	ApiProvider string
}

func GetClient() Client {
	return Client{
		ApiKey:      os.Getenv("API_KEY"),
		ApiProvider: os.Getenv("API_PROVIDER"),
	}
}
