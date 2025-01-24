package services

import (
	"context"

	"github.com/labstack/echo/v4"

	"github.com/MaKcm14/best-price-service/price-service/internal/entities"
	"github.com/MaKcm14/best-price-service/price-service/internal/entities/dto"
)

type (
	Driver interface {
		NewContext() context.Context
		Close()
	}

	CommonParser interface {
		GetProducts(ctx echo.Context, request dto.ProductRequest) (entities.ProductSample, error)
	}

	PriceParser interface {
		GetProductsWithPriceRange(ctx echo.Context, request dto.ProductRequest) (entities.ProductSample, error)
		GetProductsWithExactPrice(ctx echo.Context, request dto.ProductRequest) (entities.ProductSample, error)
		GetProductsWithBestPrice(ctx echo.Context, request dto.ProductRequest) (entities.ProductSample, error)
	}

	ApiInteractor interface {
		PriceParser
		CommonParser
	}
)
