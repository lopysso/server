package config

import (
	"sync"

	config2 "github.com/lopysso/server/config"
)

var configService config2.Config

var once sync.Once

func NewConfig() config2.Config {

	return config2.Load()
}

func ConfigInstance() config2.Config {
	once.Do(func() {
		configService = NewConfig()
	})

	return configService
}
