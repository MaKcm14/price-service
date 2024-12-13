package controller

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/MaKcm14/best-price-service/price-service/internal/entities"
	"github.com/MaKcm14/best-price-service/price-service/internal/services"
	"github.com/labstack/echo/v4"
)

type (
	// HttpController handles the clients requests.
	HttpController struct {
		contr  *echo.Echo
		logger *slog.Logger
		filter services.Filter
	}

	ResponseErr struct {
		Error string `json:"error"`
	}
)

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

func (httpContr *HttpController) validProductRequest(cont echo.Context) (services.ProductRequest, error) {
	product := services.NewProductRequest()

	product.ProductName = cont.Param("product_name")

	sample, _ := strconv.Atoi(cont.QueryParam("sample"))

	if sample < 0 {
		return services.ProductRequest{}, ErrRequestInfo
	}
	product.Sample = sample

	markets := cont.QueryParam("markets")

	for _, market := range strings.Split(markets, " ") {
		if market == "avito" {
			product.Markets = append(product.Markets, entities.Avito)
		} else if market == "wildberries" {
			product.Markets = append(product.Markets, entities.Wildberries)
		} else if market == "yandex_market" {
			product.Markets = append(product.Markets, entities.YandexMarket)
		} else if market == "ozon" {
			product.Markets = append(product.Markets, entities.Ozon)
		}
	}

	if len(markets) == 0 {
		return services.ProductRequest{}, ErrRequestInfo
	}

	return product, nil
}

// filterByPriceUpDown defines the logic of the handling the filter-by-price-down-up requests.
func (httpContr *HttpController) filterByPriceUpDown(cont echo.Context) error {
	const filterType = "price-range-filter"

	requestInfo, err := httpContr.validProductRequest(cont)

	if err != nil {
		httpContr.logger.Warn(fmt.Errorf("error of the %v: %v", filterType, err).Error())
		return cont.JSON(http.StatusBadRequest, ResponseErr{err.Error()})
	}
	priceDown, _ := strconv.Atoi(cont.QueryParam("price_down"))
	priceUp, _ := strconv.Atoi(cont.QueryParam("price_up"))

	if priceDown < 0 || priceUp < 0 || priceUp < priceDown {
		httpContr.logger.Warn(fmt.Errorf("error of the %v: %v", filterType, ErrRequestInfo).Error())
		return cont.JSON(http.StatusBadRequest, ResponseErr{ErrRequestInfo.Error()})
	}

	products, err := httpContr.filter.FilterByPriceRange(requestInfo, priceDown, priceUp)

	if err != nil {
		httpContr.logger.Warn(fmt.Errorf("error of the %v: %v", filterType, err).Error())
		return cont.JSON(http.StatusInternalServerError, ResponseErr{ErrServerHandling.Error()})
	}

	return cont.JSON(http.StatusOK, products)
}

// filterByBestPrice defines the logic of the handling the filter-by-minimal-price requests.
func (httpContr *HttpController) filterByBestPrice(cont echo.Context) error {
	const filterType = "best-price-filter"

	requestInfo, err := httpContr.validProductRequest(cont)

	if err != nil {
		httpContr.logger.Warn(fmt.Errorf("error of the %v: %v", filterType, err).Error())
		return cont.JSON(http.StatusBadRequest, ResponseErr{err.Error()})
	}

	products, err := httpContr.filter.FilterByBestPrice(requestInfo)

	if err != nil {
		httpContr.logger.Warn(fmt.Errorf("error of the %v: %v", filterType, err).Error())
		return cont.JSON(http.StatusInternalServerError, ResponseErr{ErrServerHandling.Error()})
	}

	return cont.JSON(http.StatusOK, products)
}

// filterByExactPrice defines the logic of the handling the filter-by-set-price requests.
func (httpContr *HttpController) filterByExactPrice(cont echo.Context) error {
	const filterType = "exact-price-filter"

	requestInfo, err := httpContr.validProductRequest(cont)

	if err != nil {
		httpContr.logger.Warn(fmt.Errorf("error of the %v: %v", filterType, err).Error())
		return cont.JSON(http.StatusBadRequest, ResponseErr{err.Error()})
	}

	exactPrice, _ := strconv.Atoi(cont.QueryParam("price"))

	if exactPrice <= 0 {
		httpContr.logger.Warn(fmt.Errorf("error of the %v: %v", filterType, ErrRequestInfo).Error())
		return cont.JSON(http.StatusBadRequest, ResponseErr{ErrRequestInfo.Error()})
	}

	products, err := httpContr.filter.FilterByExactPrice(requestInfo, exactPrice)

	if err != nil {
		httpContr.logger.Warn(fmt.Errorf("error of the %v: %v", filterType, err).Error())
		return cont.JSON(http.StatusInternalServerError, ResponseErr{ErrServerHandling.Error()})
	}

	return cont.JSON(http.StatusOK, products)
}

// filterByMarkets defines the logic of the handling the filter-by-markets requests.
func (httpContr *HttpController) filterByMarkets(cont echo.Context) error {
	const filterType = "markets-filter"

	requestInfo, err := httpContr.validProductRequest(cont)

	if err != nil {
		httpContr.logger.Warn(fmt.Errorf("error of the %v: %v", filterType, err).Error())
		return cont.JSON(http.StatusBadRequest, ResponseErr{ErrRequestInfo.Error()})
	}

	products, err := httpContr.filter.FilterByMarkets(requestInfo)

	if err != nil {
		httpContr.logger.Warn(fmt.Errorf("error of the %v: %v", filterType, ErrServerHandling).Error())
		return cont.JSON(http.StatusInternalServerError, ResponseErr{ErrServerHandling.Error()})
	}

	return cont.JSON(http.StatusOK, products)
}
