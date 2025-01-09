package chttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/MaKcm14/best-price-service/price-service/internal/services"
	"github.com/MaKcm14/best-price-service/price-service/internal/services/filter"
)

// Controller handles the clients' requests.
type Controller struct {
	contr  *echo.Echo
	logger *slog.Logger
	filter filter.Filter
	valid  validator
}

func NewController(contr *echo.Echo, logger *slog.Logger, filter filter.Filter) Controller {
	return Controller{
		contr:  contr,
		logger: logger,
		filter: filter,
	}
}

// Run configures and starts the http-server.
func (с Controller) Run(socket string) {
	defer с.logger.Info("the http-server was stopped")
	defer с.contr.Close()

	с.logger.Info("configuring and starting the http-server begun")

	с.configController()
	err := с.contr.Start(socket)

	if err != nil {
		serverErr := fmt.Errorf("http-server wasn't started: %v", err)
		с.logger.Error(serverErr.Error())
		panic(serverErr)
	}
}

// configController configurates the controller by setting the middlewares and paths.
func (c *Controller) configController() {
	c.configMW()
	c.configPaths()
}

// configPaths configurates the controller's paths and methods.
func (c *Controller) configPaths() {
	c.contr.GET("/products/filter/price/price-range", c.handlePriceRangeRequest)
	c.contr.GET("/products/filter/price/best-price", c.handleBestPriceRequest)
	c.contr.GET("/products/filter/price/exact-price", c.handleExactPriceRequest)
	c.contr.GET("/products/filter/markets", c.handleMarketsRequest)
}

// configMW configurates the controller's middleware.
func (c *Controller) configMW() {
	c.contr.Use(middleware.BodyLimit("600K"))
	c.contr.Use(middleware.Gzip())

	c.contr.HTTPErrorHandler = func(err error, ctx echo.Context) {
		if errHttp, flagCheck := err.(*echo.HTTPError); flagCheck {
			if errHttp.Code == http.StatusNotFound {
				c.logger.Warn("the wrong request path was got")
				ctx.JSON(http.StatusNotFound, ResponseErr{ErrRequestPath.Error()})
			} else {
				c.logger.Warn(fmt.Sprintf("%v", err))
				ctx.JSON(http.StatusBadRequest, ResponseErr{ErrRequest.Error()})
			}
			return
		}

		if err != nil {
			c.logger.Warn(fmt.Sprintf("%v", err))
			ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrServerHandling.Error()})
		}
	}
}

// setBasicHeaders sets the headers that are common for every successfull response.
func (c *Controller) setBasicHeaders(ctx echo.Context, buf []byte) {
	ctx.Response().Header().Add("Connection", "keep-alive")
	ctx.Response().Header().Add("Content-Language", "en, ru")
	ctx.Response().Header().Add("Content-Length", fmt.Sprintf("%d", len(buf)+1))
}

// filterByPriceUpDown defines the logic of the handling the filter-by-price-down-up requests.
func (c *Controller) handlePriceRangeRequest(ctx echo.Context) error {
	const filterType = "price-range-filter"

	requestInfo, err := c.valid.validProductRequest(ctx,
		c.valid.validQuery,
		c.valid.validMarkets,
		c.valid.validAmount,
		c.valid.validSample,
		c.valid.validSort,
		c.valid.validNoImage,
	)

	if err != nil {
		c.logger.Warn(fmt.Sprintf("error of the %v: %v", filterType, err))
		return ctx.JSON(http.StatusBadRequest, ResponseErr{ErrRequestInfo.Error()})
	}
	priceDown, _ := strconv.Atoi(ctx.QueryParam("price_down"))
	priceUp, _ := strconv.Atoi(ctx.QueryParam("price_up"))

	if priceDown < 0 || priceUp <= 0 || priceUp < priceDown {
		c.logger.Warn(fmt.Sprintf("error of the %v: %v", filterType, ErrRequestInfo))
		return ctx.JSON(http.StatusBadRequest, ResponseErr{ErrRequestInfo.Error()})
	}

	products, err := c.filter.FilterByPriceRange(ctx, requestInfo, priceDown, priceUp)

	if err != nil {
		c.logger.Warn(fmt.Sprintf("error of the %v: %v", filterType, err))

		if errors.Is(err, services.ErrGettingProducts) {
			return ctx.JSON(http.StatusBadGateway, ResponseErr{ErrExternalServer.Error()})
		}

		return ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrServerHandling.Error()})
	}

	response := NewProductResponse(products)
	buf, _ := json.Marshal(response)

	c.setBasicHeaders(ctx, buf)

	ctx.Response().Header().Add("Content-Type", "application/json; charset=utf-8")
	ctx.Response().Header().Add("Cache-Control", "public, max-age=43200")

	return ctx.JSON(http.StatusOK, response)
}

// filterByBestPrice defines the logic of the handling the filter-by-minimal-price requests.
func (c *Controller) handleBestPriceRequest(ctx echo.Context) error {
	const filterType = "best-price-filter"

	requestInfo, err := c.valid.validProductRequest(ctx,
		c.valid.validQuery,
		c.valid.validMarkets,
		c.valid.validAmount,
		c.valid.validSample,
		c.valid.validNoImage,
	)

	if err != nil {
		c.logger.Warn(fmt.Sprintf("error of the %v: %v", filterType, err))
		return ctx.JSON(http.StatusBadRequest, ResponseErr{ErrRequestInfo.Error()})
	}

	products, err := c.filter.FilterByBestPrice(ctx, requestInfo)

	if err != nil {
		c.logger.Warn(fmt.Sprintf("error of the %v: %v", filterType, err))

		if errors.Is(err, services.ErrGettingProducts) {
			return ctx.JSON(http.StatusBadGateway, ResponseErr{ErrExternalServer.Error()})
		}

		return ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrServerHandling.Error()})
	}

	response := NewProductResponse(products)
	buf, _ := json.Marshal(response)

	c.setBasicHeaders(ctx, buf)

	ctx.Response().Header().Add("Content-Type", "application/json; charset=utf-8")
	ctx.Response().Header().Add("Cache-Control", "public, max-age=43200")

	return ctx.JSON(http.StatusOK, response)
}

// filterByExactPrice defines the logic of the handling the filter-by-set-price requests.
func (c *Controller) handleExactPriceRequest(ctx echo.Context) error {
	const filterType = "exact-price-filter"

	requestInfo, err := c.valid.validProductRequest(ctx,
		c.valid.validQuery,
		c.valid.validMarkets,
		c.valid.validAmount,
		c.valid.validSample,
		c.valid.validSort,
		c.valid.validNoImage,
	)

	if err != nil {
		c.logger.Warn(fmt.Sprintf("error of the %v: %v", filterType, err))
		return ctx.JSON(http.StatusBadRequest, ResponseErr{ErrRequestInfo.Error()})
	}

	exactPrice, _ := strconv.Atoi(ctx.QueryParam("price"))

	if exactPrice <= 0 {
		c.logger.Warn(fmt.Sprintf("error of the %v: %v", filterType, ErrRequestInfo))
		return ctx.JSON(http.StatusBadRequest, ResponseErr{ErrRequestInfo.Error()})
	}

	products, err := c.filter.FilterByExactPrice(ctx, requestInfo, exactPrice)

	if err != nil {
		c.logger.Warn(fmt.Sprintf("error of the %v: %v", filterType, err))

		if errors.Is(err, services.ErrGettingProducts) {
			return ctx.JSON(http.StatusBadGateway, ResponseErr{ErrExternalServer.Error()})
		}

		return ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrServerHandling.Error()})
	}
	response := NewProductResponse(products)
	buf, _ := json.Marshal(response)

	c.setBasicHeaders(ctx, buf)

	ctx.Response().Header().Add("Content-Type", "application/json; charset=utf-8")
	ctx.Response().Header().Add("Cache-Control", "public, max-age=43200")

	return ctx.JSON(http.StatusOK, response)
}

// filterByMarkets defines the logic of the handling the filter-by-markets requests.
func (c *Controller) handleMarketsRequest(ctx echo.Context) error {
	const filterType = "markets-filter"

	requestInfo, err := c.valid.validProductRequest(ctx,
		c.valid.validQuery,
		c.valid.validMarkets,
		c.valid.validAmount,
		c.valid.validSample,
		c.valid.validSort,
		c.valid.validNoImage,
	)

	if err != nil {
		c.logger.Warn(fmt.Sprintf("error of the %v: %v", filterType, err))
		return ctx.JSON(http.StatusBadRequest, ResponseErr{err.Error()})
	}

	products, err := c.filter.FilterByMarkets(ctx, requestInfo)

	if err != nil {
		c.logger.Warn(fmt.Sprintf("error of the %v: %v", filterType, ErrServerHandling))

		if errors.Is(err, services.ErrGettingProducts) {
			return ctx.JSON(http.StatusBadGateway, ResponseErr{ErrExternalServer.Error()})
		}

		return ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrServerHandling.Error()})
	}

	response := NewProductResponse(products)
	buf, _ := json.Marshal(response)

	c.setBasicHeaders(ctx, buf)

	ctx.Response().Header().Add("Content-Type", "application/json; charset=utf-8")
	ctx.Response().Header().Add("Cache-Control", "public, max-age=43200")

	return ctx.JSON(http.StatusOK, response)
}
