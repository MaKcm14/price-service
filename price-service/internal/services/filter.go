package services

import (
	"log/slog"

	"github.com/MaKcm14/best-price-service/price-service/internal/entities"
)

type (
	MarketFilter struct {
		logger *slog.Logger
	}

	PriceFilter struct {
		logger *slog.Logger
	}

	ProductsFilter struct {
		PriceFilter
		MarketFilter
	}
)

func NewFilter(log *slog.Logger) *ProductsFilter {
	return &ProductsFilter{
		PriceFilter:  NewPriceFilter(log),
		MarketFilter: NewMarketFilter(log),
	}
}

func NewMarketFilter(log *slog.Logger) MarketFilter {
	return MarketFilter{
		logger: log,
	}
}

// FilterByMarkets defines the logic of the getting and processing the products' sample
// from the markets' responses filtered only by markets.
func (filter *MarketFilter) FilterByMarkets(product ProductRequest) ([]entities.Product, error) {
	return nil, nil
}

func NewPriceFilter(log *slog.Logger) PriceFilter {
	return PriceFilter{
		logger: log,
	}
}

// FilterByPriceRange defines the logic of the getting and processing the products' sample
// from the markets' responses constrained by the markets' filters and two boundaries of
// the price range.
func (filter *PriceFilter) FilterByPriceRange(product ProductRequest, priceDown int, priceUp int) ([]entities.Product, error) {
	return nil, nil
}

// FilterBestPrice defines the logic of the getting and processing the products' sample
// from the markets' responses contrained by the markets' filters and the minimal price of the sample.
func (filter *PriceFilter) FilterByBestPrice(product ProductRequest) ([]entities.Product, error) {
	return nil, nil
}

// FilterByExactPrice defines the logic of the getting and processing the products' sample
// from the markets' responses constrained by the markets' filters and the products that
// have got the exactest prices to the client's price.
func (filter *PriceFilter) FilterByExactPrice(product ProductRequest, exactPrice int) ([]entities.Product, error) {
	return nil, nil
}
