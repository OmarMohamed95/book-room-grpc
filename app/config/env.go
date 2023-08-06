package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error while reading env file %s", err)
	}
}

func Get(key string) string {
	value, valueOk := viper.Get(key).(string)

	if !valueOk {
		logMessage := fmt.Sprintf("%s must be string", key)
		panic(logMessage)
	}

	return value
}
