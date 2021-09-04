package config

import "time"

//Server to hold server configuration values
type Server struct {
	Address      string
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
	IddleTimeout time.Duration
}

func GetServer() Server {
	return Server{
		Address:      "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IddleTimeout: time.Second * 60,
	}
}
