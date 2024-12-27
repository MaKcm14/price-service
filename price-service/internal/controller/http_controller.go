package controller

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/MaKcm14/best-price-service/price-service/internal/entities"
	"github.com/MaKcm14/best-price-service/price-service/internal/services"
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
	defer httpContr.logger.Info("the http-server was stopped")
	defer httpContr.contr.Close()

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
	httpContr.contr.GET("/products/filter/price/price-range", httpContr.filterByPriceUpDown)
	httpContr.contr.GET("/products/filter/price/best-price", httpContr.filterByBestPrice)
	httpContr.contr.GET("/products/filter/price/exact-price", httpContr.filterByExactPrice)
	httpContr.contr.GET("/products/filter/markets", httpContr.filterByMarkets)

	httpContr.contr.HTTPErrorHandler = func(err error, cont echo.Context) {
		if httpErr, flagCheck := err.(*echo.HTTPError); flagCheck {
			if httpErr.Code == http.StatusNotFound {
				httpContr.logger.Warn("the wrong request path was got")
				cont.JSON(http.StatusNotFound, ResponseErr{ErrRequestPath.Error()})
			}
		}
	}
}

func (httpContr *HttpController) isDataSafe(data string) bool {
	data = strings.ToLower(data)

	for _, elem := range data {
		if string(elem) == "=" {
			return false
		}
	}

	if strings.Contains(data, "drop ") || strings.Contains(data, "union ") || strings.Contains(data, "--") {
		return false
	}

	return true
}

func (httpContr *HttpController) validProductRequest(ctx echo.Context) (entities.ProductRequest, error) {
	product := entities.NewProductRequest()

	if query := ctx.QueryParam("query"); httpContr.isDataSafe(query) {
		product.ProductName, _ = url.QueryUnescape(query)
	} else {
		return entities.ProductRequest{}, ErrRequestInfo
	}

	sample, err := strconv.Atoi(ctx.QueryParam("sample"))

	if sample < 0 {
		return entities.ProductRequest{}, ErrRequestInfo
	} else if err != nil {
		return entities.ProductRequest{}, err
	}
	product.Sample = sample

	markets := ctx.QueryParam("markets")

	for _, market := range strings.Split(markets, " ") {
		if market == "wildberries" {
			product.Markets = append(product.Markets, entities.Wildberries)
		} else if market == "ozon" {
			product.Markets = append(product.Markets, entities.Ozon)
		} else if market == "megamarket" {
			product.Markets = append(product.Markets, entities.MegaMarket)
		}
	}
	amt := ctx.QueryParam("amount")

	if len(product.Markets) == 0 {
		return entities.ProductRequest{}, ErrRequestInfo
	}

	if amt != "max" && amt != "min" {
		amt = "min"
	}
	product.Amount = amt

	return product, nil
}

// filterByPriceUpDown defines the logic of the handling the filter-by-price-down-up requests.
func (httpContr *HttpController) filterByPriceUpDown(ctx echo.Context) error {
	const filterType = "price-range-filter"

	requestInfo, err := httpContr.validProductRequest(ctx)

	if err != nil {
		httpContr.logger.Warn(fmt.Sprintf("error of the %v: %v", filterType, err))
		return ctx.JSON(http.StatusBadRequest, ResponseErr{ErrRequestInfo.Error()})
	}
	priceDown, _ := strconv.Atoi(ctx.QueryParam("price_down"))
	priceUp, _ := strconv.Atoi(ctx.QueryParam("price_up"))

	if priceDown < 0 || priceUp < 0 || priceUp < priceDown {
		httpContr.logger.Warn(fmt.Sprintf("error of the %v: %v", filterType, ErrRequestInfo))
		return ctx.JSON(http.StatusBadRequest, ResponseErr{ErrRequestInfo.Error()})
	}

	products, err := httpContr.filter.FilterByPriceRange(ctx, requestInfo, priceDown, priceUp)

	if err != nil {
		httpContr.logger.Warn(fmt.Sprintf("error of the %v: %v", filterType, err))
		return ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrServerHandling.Error()})
	}

	ctx.Response().Header().Add("Cache-Control", "public,max-age=43200")

	return ctx.JSON(http.StatusOK, products)
}

// filterByBestPrice defines the logic of the handling the filter-by-minimal-price requests.
func (httpContr *HttpController) filterByBestPrice(ctx echo.Context) error {
	const filterType = "best-price-filter"

	requestInfo, err := httpContr.validProductRequest(ctx)

	if err != nil {
		httpContr.logger.Warn(fmt.Sprintf("error of the %v: %v", filterType, err))
		return ctx.JSON(http.StatusBadRequest, ResponseErr{ErrRequestInfo.Error()})
	}

	products, err := httpContr.filter.FilterByBestPrice(ctx, requestInfo)

	if err != nil {
		httpContr.logger.Warn(fmt.Sprintf("error of the %v: %v", filterType, err))
		return ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrServerHandling.Error()})
	}

	ctx.Response().Header().Add("Cache-Control", "public,max-age=43200")

	return ctx.JSON(http.StatusOK, products)
}

// filterByExactPrice defines the logic of the handling the filter-by-set-price requests.
func (httpContr *HttpController) filterByExactPrice(ctx echo.Context) error {
	const filterType = "exact-price-filter"

	requestInfo, err := httpContr.validProductRequest(ctx)

	if err != nil {
		httpContr.logger.Warn(fmt.Sprintf("error of the %v: %v", filterType, err))
		return ctx.JSON(http.StatusBadRequest, ResponseErr{ErrRequestInfo.Error()})
	}

	exactPrice, _ := strconv.Atoi(ctx.QueryParam("price"))

	if exactPrice <= 0 {
		httpContr.logger.Warn(fmt.Sprintf("error of the %v: %v", filterType, ErrRequestInfo))
		return ctx.JSON(http.StatusBadRequest, ResponseErr{ErrRequestInfo.Error()})
	}

	products, err := httpContr.filter.FilterByExactPrice(ctx, requestInfo, exactPrice)

	if err != nil {
		httpContr.logger.Warn(fmt.Sprintf("error of the %v: %v", filterType, err))
		return ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrServerHandling.Error()})
	}

	ctx.Response().Header().Add("Cache-Control", "public,max-age=43200")

	return ctx.JSON(http.StatusOK, products)
}

// filterByMarkets defines the logic of the handling the filter-by-markets requests.
func (httpContr *HttpController) filterByMarkets(ctx echo.Context) error {
	const filterType = "markets-filter"

	requestInfo, err := httpContr.validProductRequest(ctx)

	if err != nil {
		httpContr.logger.Warn(fmt.Sprintf("error of the %v: %v", filterType, err))
		return ctx.JSON(http.StatusBadRequest, ResponseErr{err.Error()})
	}

	products, err := httpContr.filter.FilterByMarkets(ctx, requestInfo)

	if err != nil {
		httpContr.logger.Warn(fmt.Sprintf("error of the %v: %v", filterType, ErrServerHandling))
		return ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrServerHandling.Error()})
	}

	ctx.Response().Header().Add("Cache-Control", "public,max-age=43200")

	return ctx.JSON(http.StatusOK, products)
}
