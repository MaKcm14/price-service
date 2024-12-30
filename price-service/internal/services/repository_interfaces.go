package services

import (
	"context"

	"github.com/labstack/echo/v4"

	"github.com/MaKcm14/best-price-service/price-service/internal/entities"
)

type (
	Driver interface {
		NewContext() context.Context
		Close()
	}

	MarketsSifter interface {
		GetProducts(ctx echo.Context, request entities.ProductRequest) (entities.ProductSample, error)
	}

	PriceSifter interface {
		GetProductsWithPriceRange(ctx echo.Context, request entities.ProductRequest, priceDown, priceUp int) (entities.ProductSample, error)
		GetProductsWithExactPrice(ctx echo.Context, request entities.ProductRequest, exactPrice int) (entities.ProductSample, error)
		GetProductsWithBestPrice(ctx echo.Context, request entities.ProductRequest) (entities.ProductSample, error)
	}

	ApiInteractor interface {
		PriceSifter
		MarketsSifter
	}
)
