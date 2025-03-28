package services

import (
	"context"

	"github.com/labstack/echo/v4"

	"github.com/MaKcm14/price-service/internal/entities/dto"
	"github.com/MaKcm14/price-service/pkg/entities"
)

type (
	Driver interface {
		Closer
		NewContext() context.Context
	}

	Closer interface {
		Close()
	}

	AsyncWriter interface {
		Closer
		SendProductsMessage(products []entities.ProductSample, request dto.ProductRequest)
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
