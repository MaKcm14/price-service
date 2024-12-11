package app

import (
	"log/slog"
	"os"

	"github.com/MaKcm14/best-price-service/price-service/internal/config"
)

type App struct {
	logger  *slog.Logger
	logFile *os.File
}

func NewApp() *App {
	logFile, _ := os.Create("log.txt")
	log := slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{Level: slog.LevelInfo}))

	log.Info("configuring begun")

	conf, err := config.NewConfig(log)

	if err != nil {
		logFile.Close()
		panic(err)
	}

	_ = conf

	return &App{
		logger:  log,
		logFile: logFile,
	}
}

func (app *App) Run() {
	defer app.logFile.Close()
	defer app.logger.Info("the app was STOPPED")

	app.logger.Info("the app was STARTED")
}
