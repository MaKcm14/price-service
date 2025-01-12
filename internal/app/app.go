package app

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"

	"github.com/MaKcm14/best-price-service/price-service/internal/config"
	"github.com/MaKcm14/best-price-service/price-service/internal/controller/chttp"
	"github.com/MaKcm14/best-price-service/price-service/internal/entities"
	"github.com/MaKcm14/best-price-service/price-service/internal/repository/api"
	"github.com/MaKcm14/best-price-service/price-service/internal/repository/api/mmega"
	"github.com/MaKcm14/best-price-service/price-service/internal/repository/api/wildb"
	"github.com/MaKcm14/best-price-service/price-service/internal/services"
	"github.com/MaKcm14/best-price-service/price-service/internal/services/filter"
)

// Service unions every parts of the application.
type Service struct {
	appContr    chttp.Controller
	chrome      services.Driver
	logger      *slog.Logger
	mainLogFile *os.File
	appSet      config.Settings
}

func NewService() Service {
	mainLogFile, err := os.Create("../logs/price-service-main-logs.txt")

	if err != nil {
		panic(fmt.Sprintf("error of creating the main-log-file: %v", err))
	}

	log := slog.New(slog.NewTextHandler(mainLogFile, &slog.HandlerOptions{Level: slog.LevelInfo}))

	log.Info("main application's configuring begun")

	appSet, err := config.NewSettings(log, config.Socket)

	if err != nil {
		mainLogFile.Close()
		panic(err)
	}

	chrome := api.NewChromePull()

	return Service{
		appContr: chttp.NewController(echo.New(), log,
			filter.New(
				log,
				map[entities.Market]services.ApiInteractor{
					entities.Wildberries: wildb.NewWildberriesAPI(chrome.NewContext(), log, 1),
					entities.MegaMarket:  mmega.NewMegaMarketAPI(chrome.NewContext(), log),
				})),
		logger:      log,
		mainLogFile: mainLogFile,
		chrome:      chrome,
		appSet:      appSet,
	}
}

// Run starts the configured application.
func (s Service) Run() {
	defer s.chrome.Close()
	defer s.mainLogFile.Close()
	defer s.logger.Info("the app was STOPPED")

	s.logger.Info("the app was STARTED")
	s.appContr.Run(s.appSet.Socket)
}
