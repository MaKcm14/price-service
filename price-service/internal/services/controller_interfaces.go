package services

import "github.com/MaKcm14/best-price-service/price-service/internal/entities"

type (
	MarketFilterAdapter interface {
		FilterByMarkets(product entities.ProductRequest) ([]entities.Product, error)
	}

	PriceFilterAdapter interface {
		FilterByPriceRange(product entities.ProductRequest, priceDown int, priceUp int) ([]entities.Product, error)
		FilterByBestPrice(product entities.ProductRequest) ([]entities.Product, error)
		FilterByExactPrice(product entities.ProductRequest, exactPrice int) ([]entities.Product, error)
	}

	Filter interface {
		MarketFilterAdapter
		PriceFilterAdapter
	}
)
