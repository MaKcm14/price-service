package services

import (
	"github.com/labstack/echo/v4"

	"github.com/MaKcm14/best-price-service/price-service/internal/entities"
)

type ApiInteractor interface {
	GetProducts(ctx echo.Context, product entities.ProductRequest) ([]entities.ProductResponse, error)
	GetProductsByPriceRange(ctx echo.Context, product entities.ProductRequest, priceDown, priceUp int) ([]entities.ProductResponse, error)
	GetProductsByExactPrice(ctx echo.Context, product entities.ProductRequest, exactPrice int) ([]entities.ProductResponse, error)
	GetProductsByBestPrice(ctx echo.Context, product entities.ProductRequest) ([]entities.ProductResponse, error)
}
