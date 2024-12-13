package services

import "github.com/MaKcm14/best-price-service/price-service/internal/entities"

type (
	MarketFilterAdapter interface {
		FilterByMarkets(product ProductRequest) ([]entities.Product, error)
	}

	PriceFilterAdapter interface {
		FilterByPriceRange(product ProductRequest, priceDown int, priceUp int) ([]entities.Product, error)
		FilterByBestPrice(product ProductRequest) ([]entities.Product, error)
		FilterByExactPrice(product ProductRequest, exactPrice int) ([]entities.Product, error)
	}

	Filter interface {
		MarketFilterAdapter
		PriceFilterAdapter
	}
)
