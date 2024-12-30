package app

import (
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"

	"github.com/MaKcm14/best-price-service/price-service/internal/config"
	"github.com/MaKcm14/best-price-service/price-service/internal/controller/chttp"
	"github.com/MaKcm14/best-price-service/price-service/internal/entities"
	"github.com/MaKcm14/best-price-service/price-service/internal/repository/api"
	"github.com/MaKcm14/best-price-service/price-service/internal/repository/api/wildb"
	"github.com/MaKcm14/best-price-service/price-service/internal/services"
	"github.com/MaKcm14/best-price-service/price-service/internal/services/strainer"
)

// App unions every parts of the application.
type Service struct {
	appContr chttp.Controller
	chrome   services.Driver
	logger   *slog.Logger
	logFile  *os.File
	appSet   config.Settings
}

func NewService() Service {
	logFile, _ := os.Create("log.txt")
	log := slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{Level: slog.LevelInfo}))

	log.Info("main application's configuring begun")

	appSet, err := config.NewSettings(log, config.Socket)

	if err != nil {
		logFile.Close()
		panic(err)
	}

	chrome := api.NewChromePull()

	return Service{
		appContr: chttp.NewController(echo.New(), log,
			strainer.NewProductsFilter(
				log,
				map[entities.Market]services.ApiInteractor{
					entities.Wildberries: wildb.NewWildberriesAPI(log, 1.2, chrome.NewContext()),
				})),
		logger:  log,
		logFile: logFile,
		chrome:  chrome,
		appSet:  appSet,
	}
}

// Run starts the configured application.
func (a Service) Run() {
	defer a.chrome.Close()
	defer a.logFile.Close()
	defer a.logger.Info("the app was STOPPED")

	a.logger.Info("the app was STARTED")
	a.appContr.Run(a.appSet.Socket)
}
