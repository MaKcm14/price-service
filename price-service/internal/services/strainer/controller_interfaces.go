package strainer

import (
	"github.com/labstack/echo/v4"

	"github.com/MaKcm14/best-price-service/price-service/internal/entities"
	"github.com/MaKcm14/best-price-service/price-service/internal/entities/dto"
)

type (
	MarketsFilterAdapter interface {
		FilterByMarkets(ctx echo.Context, request dto.ProductRequest) ([]entities.ProductSample, error)
	}

	PriceFilterAdapter interface {
		FilterByPriceRange(ctx echo.Context, request dto.ProductRequest, priceDown int, priceUp int) ([]entities.ProductSample, error)
		FilterByBestPrice(ctx echo.Context, request dto.ProductRequest) ([]entities.ProductSample, error)
		FilterByExactPrice(ctx echo.Context, request dto.ProductRequest, exactPrice int) ([]entities.ProductSample, error)
	}

	Filter interface {
		MarketsFilterAdapter
		PriceFilterAdapter
	}
)
