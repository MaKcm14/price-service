package services

import (
	"github.com/labstack/echo/v4"

	"github.com/MaKcm14/best-price-service/price-service/internal/entities"
)

type ApiInteractor interface {
	GetProducts(ctx echo.Context, product entities.ProductRequest) ([]entities.Product, error)
	GetProductsByPriceRange(ctx echo.Context, product entities.ProductRequest, priceDown, priceUp int) ([]entities.Product, error)
	GetProductsByExactPrice(ctx echo.Context, product entities.ProductRequest, exactPrice int) ([]entities.Product, error)
	GetProductsByBestPrice(ctx echo.Context, product entities.ProductRequest) ([]entities.Product, error)
}
