package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type SettingOpt func(*Settings, *slog.Logger) error

// Settings sets the application's configurations.
type Settings struct {
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
func Socket(appSet *Settings, log *slog.Logger) error {
	socket, err := configEnv("SOCKET", log)

	if err != nil {
		return err
	}
	appSet.Socket = socket

	return nil
}

func NewSettings(log *slog.Logger, opts ...SettingOpt) (Settings, error) {
	appSet := Settings{}
	err := godotenv.Load("../.env")

	if err != nil {
		envErr := fmt.Errorf("error while loading the .env file (check it and try again): %v", err)
		log.Error(envErr.Error())
		return Settings{}, envErr
	}

	for _, opt := range opts {
		if err := opt(&appSet, log); err != nil {
			return Settings{}, err
		}
	}

	return appSet, nil
}
