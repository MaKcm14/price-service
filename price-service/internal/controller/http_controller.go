package controller

import (
	"fmt"
	"log/slog"

	"github.com/MaKcm14/best-price-service/price-service/internal/services"
	"github.com/labstack/echo/v4"
)

// HttpController handles the clients requests.
type HttpController struct {
	contr  *echo.Echo
	logger *slog.Logger
	filter services.Filter
}

func NewHttpController(contr *echo.Echo, logger *slog.Logger, filter services.Filter) *HttpController {
	return &HttpController{
		contr:  contr,
		logger: logger,
		filter: filter,
	}
}

// Run configures and starts the http-server.
func (httpContr *HttpController) Run() {
	httpContr.logger.Info("configuring and starting the http-server begun")

	httpContr.configPath()
	err := httpContr.contr.Start("localhost:8080")

	if err != nil {
		serverErr := fmt.Errorf("http-server wasn't started: %v", err)
		httpContr.logger.Error(serverErr.Error())
		panic(serverErr)
	}
}

func (httpContr *HttpController) configPath() {
	httpContr.contr.GET("/products/filter/price/price-range/:product_name", httpContr.filterByPriceUpDown)
	httpContr.contr.GET("/products/filter/price/best-price/:product_name", httpContr.filterByBestPrice)
	httpContr.contr.GET("/products/filter/price/exact-price/:product_name", httpContr.filterByExactPrice)
	httpContr.contr.GET("/products/filter/markets/:product_name", httpContr.filterByMarkets)
}

// filterByPriceUpDown defines the logic of the handling the filter-by-price-down-up requests.
func (httpContr *HttpController) filterByPriceUpDown(cont echo.Context) error {
	return nil
}

// filterByBestPrice defines the logic of the handling the filter-by-minimal-price requests.
func (httpContr *HttpController) filterByBestPrice(cont echo.Context) error {
	return nil
}

// filterByExactPrice defines the logic of the handling the filter-by-set-price requests.
func (httpContr *HttpController) filterByExactPrice(cont echo.Context) error {
	return nil
}

// filterByMarkets defines the logic of the handling the filter-by-markets requests.
func (httpContr *HttpController) filterByMarkets(cont echo.Context) error {
	return nil
}
