package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

// Config sets the application's configurations.
type Config struct {
	DSN string
}

// configEnv gets ENV var. It returns the error if var is unset or unexisting.
func configEnv(key string, log *slog.Logger) (string, error) {
	env := os.Getenv(key)

	if len(env) == 0 {
		err := fmt.Errorf("error while parsing the .env file: check the %s var is set correctly", key)
		log.Error(err.Error())
		return "", err
	}

	return env, nil
}

func NewConfig(log *slog.Logger) (Config, error) {
	config := Config{}
	err := godotenv.Load("../../.env")

	if err != nil {
		envErr := fmt.Errorf("error while loading the .env file (check it and try again): %v", err)
		log.Error(envErr.Error())
		return Config{}, envErr
	}

	config.DSN, err = configEnv("DSN", log)

	if err != nil {
		return Config{}, err
	}

	return config, nil
}
