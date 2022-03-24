package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	AppConfig
	ServerConfig
	StorageConfig
}

type AppConfig struct {
	SizeOfLRUCache int
}

type ServerConfig struct {
	Port string
}

type StorageConfig struct {
	StorageFolder string
}

func NewConfig(configFolder, storageFolder string) *Config {
	viper.SetConfigType("yml")
	viper.AddConfigPath(configFolder)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalln(err)
	}

	config.StorageConfig.StorageFolder = storageFolder

	return &config
}
