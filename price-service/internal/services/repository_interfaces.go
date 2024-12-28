package services

import (
	"github.com/labstack/echo/v4"

	"github.com/MaKcm14/best-price-service/price-service/internal/entities"
)

type (
	MarketsSifter interface {
		GetProducts(ctx echo.Context, product entities.ProductRequest) (entities.ProductSample, error)
	}

	PriceSifter interface {
		GetProductsWithPriceRange(ctx echo.Context, product entities.ProductRequest, priceDown, priceUp int) (entities.ProductSample, error)
		GetProductsWithExactPrice(ctx echo.Context, product entities.ProductRequest, exactPrice int) (entities.ProductSample, error)
		GetProductsWithBestPrice(ctx echo.Context, product entities.ProductRequest) (entities.ProductSample, error)
	}

	ApiInteractor interface {
		PriceSifter
		MarketsSifter
	}
)
