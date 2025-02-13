package filter

import (
	"github.com/labstack/echo/v4"

	"github.com/MaKcm14/price-service/internal/entities/dto"
	"github.com/MaKcm14/price-service/pkg/entities"
)

type (
	CommonFilterAdapter interface {
		FilterByMarkets(ctx echo.Context, request dto.ProductRequest) ([]entities.ProductSample, error)
	}

	PriceFilterAdapter interface {
		FilterByPriceRange(ctx echo.Context, request dto.ProductRequest) ([]entities.ProductSample, error)
		FilterByBestPrice(ctx echo.Context, request dto.ProductRequest) ([]entities.ProductSample, error)
		FilterByExactPrice(ctx echo.Context, request dto.ProductRequest) ([]entities.ProductSample, error)
	}

	Filter interface {
		CommonFilterAdapter
		PriceFilterAdapter
	}
)
