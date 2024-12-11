package app

import (
	"log/slog"
	"os"

	"github.com/MaKcm14/best-price-service/price-service/internal/config"
)

type App struct {
	logger *slog.Logger
}

func NewApp() *App {
	logFile, _ := os.Create("log.txt")
	log := slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{Level: slog.LevelInfo}))

	log.Info("configuring begun")

	conf := config.NewConfig(log)

	_ = conf

	return &App{
		logger: log,
	}
}

func (app *App) Run() {
	defer app.logger.Info("the app was STOPPED")

	app.logger.Info("the app was STARTED")
}
