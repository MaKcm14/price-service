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
func (h HttpController) Run() {
	defer h.logger.Info("the http-server was stopped")
	defer h.contr.Close()

	h.logger.Info("configuring and starting the http-server begun")

	h.configPaths()
	err := h.contr.Start(h.socket)

	if err != nil {
		serverErr := fmt.Errorf("http-server wasn't started: %v", err)
		h.logger.Error(serverErr.Error())
		panic(serverErr)
	}
}

func (h *HttpController) configPaths() {
	h.contr.GET("/products/filter/price/price-range", h.handlePriceRangeRequest)
	h.contr.GET("/products/filter/price/best-price", h.handleBestPriceRequest)
	h.contr.GET("/products/filter/price/exact-price", h.handleExactPriceRequest)
	h.contr.GET("/products/filter/markets", h.handleMarketsRequest)

	h.contr.HTTPErrorHandler = func(err error, cont echo.Context) {
		if err, flagCheck := err.(*echo.HTTPError); flagCheck {
			if err.Code == http.StatusNotFound {
				h.logger.Warn("the wrong request path was got")
				cont.JSON(http.StatusNotFound, ResponseErr{ErrRequestPath.Error()})
			}
		}
	}
}

// filterByPriceUpDown defines the logic of the handling the filter-by-price-down-up requests.
func (h *HttpController) handlePriceRangeRequest(ctx echo.Context) error {
	const filterType = "price-range-filter"

	requestInfo, err := h.valid.validProductRequest(ctx,
		h.valid.validQuery,
		h.valid.validMarkets,
		h.valid.validAmount,
		h.valid.validSample,
		h.valid.validSort,
		h.valid.validNoImage,
	)

	if err != nil {
		h.logger.Warn(fmt.Sprintf("error of the %v: %v", filterType, err))
		return ctx.JSON(http.StatusBadRequest, ResponseErr{ErrRequestInfo.Error()})
	}
	priceDown, _ := strconv.Atoi(ctx.QueryParam("price_down"))
	priceUp, _ := strconv.Atoi(ctx.QueryParam("price_up"))

	if priceDown < 0 || priceUp <= 0 || priceUp < priceDown {
		h.logger.Warn(fmt.Sprintf("error of the %v: %v", filterType, ErrRequestInfo))
		return ctx.JSON(http.StatusBadRequest, ResponseErr{ErrRequestInfo.Error()})
	}

	products, err := h.filter.FilterByPriceRange(ctx, requestInfo, priceDown, priceUp)

	if err != nil {
		h.logger.Warn(fmt.Sprintf("error of the %v: %v", filterType, err))
		return ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrServerHandling.Error()})
	}

	ctx.Response().Header().Add("Cache-Control", "public,max-age=43200")

	return ctx.JSON(http.StatusOK, products)
}

// filterByBestPrice defines the logic of the handling the filter-by-minimal-price requests.
func (h *HttpController) handleBestPriceRequest(ctx echo.Context) error {
	const filterType = "best-price-filter"

	requestInfo, err := h.valid.validProductRequest(ctx,
		h.valid.validQuery,
		h.valid.validMarkets,
		h.valid.validAmount,
		h.valid.validSample,
		h.valid.validNoImage,
	)

	if err != nil {
		h.logger.Warn(fmt.Sprintf("error of the %v: %v", filterType, err))
		return ctx.JSON(http.StatusBadRequest, ResponseErr{ErrRequestInfo.Error()})
	}

	products, err := h.filter.FilterByBestPrice(ctx, requestInfo)

	if err != nil {
		h.logger.Warn(fmt.Sprintf("error of the %v: %v", filterType, err))
		return ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrServerHandling.Error()})
	}

	ctx.Response().Header().Add("Cache-Control", "public,max-age=43200")

	return ctx.JSON(http.StatusOK, products)
}

// filterByExactPrice defines the logic of the handling the filter-by-set-price requests.
func (h *HttpController) handleExactPriceRequest(ctx echo.Context) error {
	const filterType = "exact-price-filter"

	requestInfo, err := h.valid.validProductRequest(ctx,
		h.valid.validQuery,
		h.valid.validMarkets,
		h.valid.validAmount,
		h.valid.validSample,
		h.valid.validSort,
		h.valid.validNoImage,
	)

	if err != nil {
		h.logger.Warn(fmt.Sprintf("error of the %v: %v", filterType, err))
		return ctx.JSON(http.StatusBadRequest, ResponseErr{ErrRequestInfo.Error()})
	}

	exactPrice, _ := strconv.Atoi(ctx.QueryParam("price"))

	if exactPrice <= 0 {
		h.logger.Warn(fmt.Sprintf("error of the %v: %v", filterType, ErrRequestInfo))
		return ctx.JSON(http.StatusBadRequest, ResponseErr{ErrRequestInfo.Error()})
	}

	products, err := h.filter.FilterByExactPrice(ctx, requestInfo, exactPrice)

	if err != nil {
		h.logger.Warn(fmt.Sprintf("error of the %v: %v", filterType, err))
		return ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrServerHandling.Error()})
	}

	ctx.Response().Header().Add("Cache-Control", "public,max-age=43200")

	return ctx.JSON(http.StatusOK, products)
}

// filterByMarkets defines the logic of the handling the filter-by-markets requests.
func (h *HttpController) handleMarketsRequest(ctx echo.Context) error {
	const filterType = "markets-filter"

	requestInfo, err := h.valid.validProductRequest(ctx,
		h.valid.validQuery,
		h.valid.validMarkets,
		h.valid.validAmount,
		h.valid.validSample,
		h.valid.validSort,
		h.valid.validNoImage,
	)

	if err != nil {
		h.logger.Warn(fmt.Sprintf("error of the %v: %v", filterType, err))
		return ctx.JSON(http.StatusBadRequest, ResponseErr{err.Error()})
	}

	products, err := h.filter.FilterByMarkets(ctx, requestInfo)

	if err != nil {
		h.logger.Warn(fmt.Sprintf("error of the %v: %v", filterType, ErrServerHandling))
		return ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrServerHandling.Error()})
	}

	ctx.Response().Header().Add("Cache-Control", "public,max-age=43200")

	return ctx.JSON(http.StatusOK, products)
}
