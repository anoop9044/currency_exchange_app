package config

import (
	"log"

	"github.com/spf13/viper"
)

type Configuration struct {
	RedisAddr string
	JwtSecret string
}

var Config Configuration

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	err = viper.Unmarshal(&Config)
	if err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}
}