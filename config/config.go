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
	viper.AddConfigPath("$HOME/.config/spotify-search/")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")
	viper.BindEnv("SPOTIFY_ID")
	viper.BindEnv("SPOTIFY_SECRET")
	viper.BindEnv("AUTH_URL")
	var configuration Configuration
	if err := viper.ReadInConfig(); err != nil {
		switch err.(type) {
		default:
			panic(fmt.Errorf("Fatal error loading config file: %s \n", err))
		case viper.ConfigFileNotFoundError:
			fmt.Println("No config file found. Using defaults and environment variables")
		}
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	return &configuration
}
