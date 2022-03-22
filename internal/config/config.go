package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	AppConfig
	ServerConfig
}

type AppConfig struct {
	SizeOfLRUCache int
}

type ServerConfig struct {
	Port string
}

func NewConfig(configFolder string) *Config {
	viper.SetConfigType("yml")
	viper.AddConfigPath(configFolder)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalln(err)
	}

	return &config
}
