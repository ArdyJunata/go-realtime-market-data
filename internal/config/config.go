package config

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig(file string) {
	godotenv.Load(file)
}

func GetString(key ConfigKey) string {
	val := os.Getenv(string(key))
	if val == "" {
		return ""
	}

	return val
}
