package services

import (
	"github.com/labstack/echo/v4"

	"github.com/MaKcm14/best-price-service/price-service/internal/entities"
)

type ApiInteractor interface {
	GetProducts(ctx echo.Context, product entities.ProductRequest) ([]entities.ProductSample, error)
	GetProductsByPriceRange(ctx echo.Context, product entities.ProductRequest, priceDown, priceUp int) ([]entities.ProductSample, error)
	GetProductsByExactPrice(ctx echo.Context, product entities.ProductRequest, exactPrice int) ([]entities.ProductSample, error)
	GetProductsByBestPrice(ctx echo.Context, product entities.ProductRequest) ([]entities.ProductSample, error)
}
