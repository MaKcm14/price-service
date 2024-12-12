package controller

import (
	"fmt"
	"log/slog"

	"github.com/labstack/echo/v4"
)

// HttpController handles the clients requests.
type HttpController struct {
	contr  *echo.Echo
	logger *slog.Logger
}

func NewHttpController(contr *echo.Echo, logger *slog.Logger) *HttpController {
	return &HttpController{
		contr:  contr,
		logger: logger,
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
	httpContr.contr.GET("/products/filter/:product_name", httpContr.filterByPriceUpDown)
	httpContr.contr.GET("/products/filter/best-price/:product_name", httpContr.filterByBestPrice)
	httpContr.contr.GET("/products/filter/exact-price/:product_name", httpContr.filterByExactPrice)
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
