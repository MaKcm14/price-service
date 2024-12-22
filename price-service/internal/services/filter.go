package services

import (
	"fmt"
	"log/slog"

	"github.com/MaKcm14/best-price-service/price-service/internal/entities"
	"github.com/MaKcm14/best-price-service/price-service/internal/repository/api"
)

type ProductsFilter struct {
	logger *slog.Logger
	api    ApiInteractor
}

func NewFilter(log *slog.Logger) *ProductsFilter {
	return &ProductsFilter{
		logger: log,
		api:    api.NewMarketsApi(log),
	}
}

// FilterByMarkets defines the logic of the getting and processing the products' sample
// from the markets' responses filtered only by markets.
func (filter *ProductsFilter) FilterByMarkets(product entities.ProductRequest) ([]entities.Product, error) {
	const serviceType = "filter.service.filter_by_markets"
	products, err := filter.api.GetProducts(product)

	if err != nil {
		return nil, fmt.Errorf("error of the %v: %v", serviceType, err)
	}

	return products, nil
}

// FilterByPriceRange defines the logic of the getting and processing the products' sample
// from the markets' responses constrained by the markets' filters and two boundaries of
// the price range.
func (filter *ProductsFilter) FilterByPriceRange(product entities.ProductRequest, priceDown int, priceUp int) ([]entities.Product, error) {
	return nil, nil
}

// FilterBestPrice defines the logic of the getting and processing the products' sample
// from the markets' responses contrained by the markets' filters and the minimal price of the sample.
func (filter *ProductsFilter) FilterByBestPrice(product entities.ProductRequest) ([]entities.Product, error) {
	return nil, nil
}

// FilterByExactPrice defines the logic of the getting and processing the products' sample
// from the markets' responses constrained by the markets' filters and the products that
// have got the exactest prices to the client's price.
func (filter *ProductsFilter) FilterByExactPrice(product entities.ProductRequest, exactPrice int) ([]entities.Product, error) {
	return nil, nil
}
