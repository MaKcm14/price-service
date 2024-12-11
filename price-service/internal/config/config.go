package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DSN string
}

func configEnv(key string, log *slog.Logger) string {
	env := os.Getenv(key)

	if len(env) == 0 {
		err := fmt.Errorf("error while parsing the .env file: check the %s var is set correctly", key)
		log.Error(err.Error())
		panic(err)
	}

	return env
}

func NewConfig(log *slog.Logger) Config {
	config := Config{}
	err := godotenv.Load("../../.env")

	if err != nil {
		envErr := fmt.Errorf("error while loading the .env file (check it and try again): %v", err)
		log.Error(envErr.Error())
		panic(envErr)
	}

	config.DSN = configEnv("DSN", log)

	return config
}
