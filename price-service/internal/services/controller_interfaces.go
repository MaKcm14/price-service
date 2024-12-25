package services

import (
	"github.com/labstack/echo/v4"

	"github.com/MaKcm14/best-price-service/price-service/internal/entities"
)

type (
	MarketFilterAdapter interface {
		FilterByMarkets(ctx echo.Context, product entities.ProductRequest) ([]entities.ProductResponse, error)
	}

	PriceFilterAdapter interface {
		FilterByPriceRange(ctx echo.Context, product entities.ProductRequest, priceDown int, priceUp int) ([]entities.ProductResponse, error)
		FilterByBestPrice(ctx echo.Context, product entities.ProductRequest) ([]entities.ProductResponse, error)
		FilterByExactPrice(ctx echo.Context, product entities.ProductRequest, exactPrice int) ([]entities.ProductResponse, error)
	}

	Filter interface {
		MarketFilterAdapter
		PriceFilterAdapter
	}
)
