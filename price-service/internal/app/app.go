package app

import (
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"

	"github.com/MaKcm14/best-price-service/price-service/internal/config"
	"github.com/MaKcm14/best-price-service/price-service/internal/controller"
	"github.com/MaKcm14/best-price-service/price-service/internal/services"
)

// App unions every parts of the application.
type App struct {
	appContr Runner
	logger   *slog.Logger
	logFile  *os.File
}

type Runner interface {
	Run()
}

func NewApp() App {
	logFile, _ := os.Create("log.txt")
	log := slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{Level: slog.LevelInfo}))

	log.Info("main application's configuring begun")

	conf, err := config.NewConfig(log, config.Socket)

	if err != nil {
		logFile.Close()
		panic(err)
	}

	return App{
		appContr: controller.NewHttpController(echo.New(), log, services.NewProductsFilter(log), conf.Socket),
		logger:   log,
		logFile:  logFile,
	}
}

// Run starts the configured application.
func (app App) Run() {
	defer app.logFile.Close()
	defer app.logger.Info("the app was STOPPED")

	app.logger.Info("the app was STARTED")
	app.appContr.Run()
}
