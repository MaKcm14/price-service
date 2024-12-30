package chttp

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	filter "github.com/MaKcm14/best-price-service/price-service/internal/services/filter"
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

	с.configPaths()
	err := с.contr.Start(socket)

	if err != nil {
		serverErr := fmt.Errorf("http-server wasn't started: %v", err)
		с.logger.Error(serverErr.Error())
		panic(serverErr)
	}
}

func (c *Controller) configPaths() {
	c.contr.GET("/products/filter/price/price-range", c.handlePriceRangeRequest)
	c.contr.GET("/products/filter/price/best-price", c.handleBestPriceRequest)
	c.contr.GET("/products/filter/price/exact-price", c.handleExactPriceRequest)
	c.contr.GET("/products/filter/markets", c.handleMarketsRequest)

	c.contr.HTTPErrorHandler = func(err error, cont echo.Context) {
		if err, flagCheck := err.(*echo.HTTPError); flagCheck {
			if err.Code == http.StatusNotFound {
				c.logger.Warn("the wrong request path was got")
				cont.JSON(http.StatusNotFound, ResponseErr{ErrRequestPath.Error()})
			}
		}
	}
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
		return ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrServerHandling.Error()})
	}

	ctx.Response().Header().Add("Cache-Control", "public,max-age=43200")

	return ctx.JSON(http.StatusOK, products)
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
		return ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrServerHandling.Error()})
	}

	ctx.Response().Header().Add("Cache-Control", "public,max-age=43200")

	return ctx.JSON(http.StatusOK, products)
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
		return ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrServerHandling.Error()})
	}

	ctx.Response().Header().Add("Cache-Control", "public,max-age=43200")

	return ctx.JSON(http.StatusOK, products)
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
		return ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrServerHandling.Error()})
	}

	ctx.Response().Header().Add("Cache-Control", "public,max-age=43200")

	return ctx.JSON(http.StatusOK, products)
}
