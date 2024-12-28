package services

import (
	"github.com/labstack/echo/v4"

	"github.com/MaKcm14/best-price-service/price-service/internal/entities"
)

type (
	MarketsFilterAdapter interface {
		FilterByMarkets(ctx echo.Context, product entities.ProductRequest) ([]entities.ProductSample, error)
	}

	PriceFilterAdapter interface {
		FilterByPriceRange(ctx echo.Context, product entities.ProductRequest, priceDown int, priceUp int) ([]entities.ProductSample, error)
		FilterByBestPrice(ctx echo.Context, product entities.ProductRequest) ([]entities.ProductSample, error)
		FilterByExactPrice(ctx echo.Context, product entities.ProductRequest, exactPrice int) ([]entities.ProductSample, error)
	}

	Filter interface {
		MarketsFilterAdapter
		PriceFilterAdapter
	}
)
