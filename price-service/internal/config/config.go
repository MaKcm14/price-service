package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type ConfigOpts func(*Config, *slog.Logger) error

// Config sets the application's configurations.
type Config struct {
	DSN    string
	Socket string
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

// Socket configs the Socket ENV.
func Socket(config *Config, log *slog.Logger) error {
	socket, err := configEnv("SOCKET", log)

	if err != nil {
		return err
	}
	config.Socket = socket

	return nil
}

func NewConfig(log *slog.Logger, opts ...ConfigOpts) (Config, error) {
	config := Config{}
	err := godotenv.Load("../../.env")

	if err != nil {
		envErr := fmt.Errorf("error while loading the .env file (check it and try again): %v", err)
		log.Error(envErr.Error())
		return Config{}, envErr
	}

	for _, opt := range opts {
		if err := opt(&config, log); err != nil {
			return Config{}, err
		}
	}

	return config, nil
}
