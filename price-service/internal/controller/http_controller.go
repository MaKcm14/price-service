package controller

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/MaKcm14/best-price-service/price-service/internal/services"
)

// HttpController handles the clients' requests.
type HttpController struct {
	contr  *echo.Echo
	logger *slog.Logger
	filter services.Filter
	socket string
	valid  validator
}

func NewHttpController(contr *echo.Echo, logger *slog.Logger, filter services.Filter, socket string) HttpController {
	return HttpController{
		contr:  contr,
		logger: logger,
		filter: filter,
		socket: socket,
	}
}

// Run configures and starts the http-server.
func (httpContr HttpController) Run() {
	defer httpContr.logger.Info("the http-server was stopped")
	defer httpContr.contr.Close()

	httpContr.logger.Info("configuring and starting the http-server begun")

	httpContr.configPath()
	err := httpContr.contr.Start(httpContr.socket)

	if err != nil {
		serverErr := fmt.Errorf("http-server wasn't started: %v", err)
		httpContr.logger.Error(serverErr.Error())
		panic(serverErr)
	}
}

func (httpContr *HttpController) configPath() {
	httpContr.contr.GET("/products/filter/price/price-range", httpContr.handlePriceRangeRequest)
	httpContr.contr.GET("/products/filter/price/best-price", httpContr.handleBestPriceRequest)
	httpContr.contr.GET("/products/filter/price/exact-price", httpContr.handleExactPriceRequest)
	httpContr.contr.GET("/products/filter/markets", httpContr.handleMarketsRequest)

	httpContr.contr.HTTPErrorHandler = func(err error, cont echo.Context) {
		if httpErr, flagCheck := err.(*echo.HTTPError); flagCheck {
			if httpErr.Code == http.StatusNotFound {
				httpContr.logger.Warn("the wrong request path was got")
				cont.JSON(http.StatusNotFound, ResponseErr{ErrRequestPath.Error()})
			}
		}
	}
}

// filterByPriceUpDown defines the logic of the handling the filter-by-price-down-up requests.
func (httpContr *HttpController) handlePriceRangeRequest(ctx echo.Context) error {
	const filterType = "price-range-filter"

	requestInfo, err := httpContr.valid.validProductRequest(ctx,
		httpContr.valid.validQuery,
		httpContr.valid.validMarkets,
		httpContr.valid.validAmount,
		httpContr.valid.validSample,
	)

	if err != nil {
		httpContr.logger.Warn(fmt.Sprintf("error of the %v: %v", filterType, err))
		return ctx.JSON(http.StatusBadRequest, ResponseErr{ErrRequestInfo.Error()})
	}
	priceDown, _ := strconv.Atoi(ctx.QueryParam("price_down"))
	priceUp, _ := strconv.Atoi(ctx.QueryParam("price_up"))

	if priceDown < 0 || priceUp <= 0 || priceUp < priceDown {
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
func (httpContr *HttpController) handleBestPriceRequest(ctx echo.Context) error {
	const filterType = "best-price-filter"

	requestInfo, err := httpContr.valid.validProductRequest(ctx,
		httpContr.valid.validQuery,
		httpContr.valid.validMarkets,
		httpContr.valid.validAmount,
		httpContr.valid.validSample,
	)

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
func (httpContr *HttpController) handleExactPriceRequest(ctx echo.Context) error {
	const filterType = "exact-price-filter"

	requestInfo, err := httpContr.valid.validProductRequest(ctx,
		httpContr.valid.validQuery,
		httpContr.valid.validMarkets,
		httpContr.valid.validAmount,
		httpContr.valid.validSample,
	)

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
func (httpContr *HttpController) handleMarketsRequest(ctx echo.Context) error {
	const filterType = "markets-filter"

	requestInfo, err := httpContr.valid.validProductRequest(ctx,
		httpContr.valid.validQuery,
		httpContr.valid.validMarkets,
		httpContr.valid.validAmount,
		httpContr.valid.validSample,
	)

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
