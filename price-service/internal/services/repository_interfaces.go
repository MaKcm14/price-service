package services

import "github.com/MaKcm14/best-price-service/price-service/internal/entities"

type ApiInteractor interface {
	GetProducts(product entities.ProductRequest) ([]entities.Product, error)
	GetProductsByPriceRange(product entities.ProductRequest, priceDown, priceUp int) ([]entities.Product, error)
	GetProductsByExactPrice(product entities.ProductRequest, exactPrice int) ([]entities.Product, error)
	GetProductsByBestPrice(product entities.ProductRequest) ([]entities.Product, error)
}
