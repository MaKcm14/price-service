package services

import (
	"context"
	"log/slog"
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
		ctx context.Context
	}
)

func NewFilter(log *slog.Logger, ctx context.Context) *ProductsFilter {
	return &ProductsFilter{
		PriceFilter:  NewPriceFilter(log),
		MarketFilter: NewMarketFilter(log),
		ctx:          ctx,
	}
}

func NewMarketFilter(log *slog.Logger) MarketFilter {
	return MarketFilter{
		logger: log,
	}
}

// FilterByMarkets defines the logic of the getting and processing the products' sample
// from the markets' responses filtered only by markets.
func (filter *MarketFilter) FilterByMarkets() {

}

func NewPriceFilter(log *slog.Logger) PriceFilter {
	return PriceFilter{
		logger: log,
	}
}

// FilterByPriceRange defines the logic of the getting and processing the products' sample
// from the markets' responses constrained by the markets' filters and two boundaries of
// the price range.
func (filter *PriceFilter) FilterByPriceRange() {

}

// FilterBestPrice defines the logic of the getting and processing the products' sample
// from the markets' responses contrained by the markets' filters and the minimal price of the sample.
func (filter *PriceFilter) FilterByBestPrice() {

}

// FilterByExactPrice defines the logic of the getting and processing the products' sample
// from the markets' responses constrained by the markets' filters and the products that
// have got the exactest prices to the client's price.
func (filter *PriceFilter) FilterByExactPrice() {

}
