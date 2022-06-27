package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Configuration struct {
	SPOTIFY_ID     string
	SPOTIFY_SECRET string
	AUTH_URL       string
}

func SetupConfig() *Configuration {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")
	var configuration Configuration
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("error reading config file, %s", err)
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	return &configuration
}
