package services

import (
	"fmt"
	"log/slog"

	"github.com/labstack/echo/v4"

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
func (filter *ProductsFilter) FilterByMarkets(ctx echo.Context, product entities.ProductRequest) ([]entities.Product, error) {
	const serviceType = "filter.service.filter_by_markets"
	products, err := filter.api.GetProducts(ctx, product)

	if err != nil {
		return nil, fmt.Errorf("error of the %v: %v", serviceType, err)
	}

	return products, nil
}

// FilterByPriceRange defines the logic of the getting and processing the products' sample
// from the markets' responses constrained by the markets' filters and two boundaries of
// the price range.
func (filter *ProductsFilter) FilterByPriceRange(ctx echo.Context, product entities.ProductRequest, priceDown int, priceUp int) ([]entities.Product, error) {
	const serviceType = "filter.service.filter_by_price_range"
	products, err := filter.api.GetProductsByPriceRange(ctx, product, priceDown, priceUp)

	if err != nil {
		return nil, fmt.Errorf("error of the %v: %v", serviceType, err)
	}

	return products, nil
}

// FilterBestPrice defines the logic of the getting and processing the products' sample
// from the markets' responses contrained by the markets' filters and the minimal price of the sample.
func (filter *ProductsFilter) FilterByBestPrice(ctx echo.Context, product entities.ProductRequest) ([]entities.Product, error) {
	const serviceType = "filter.service.filter_by_best_price"
	products, err := filter.api.GetProductsByBestPrice(ctx, product)

	if err != nil {
		return nil, fmt.Errorf("error of the %v: %v", serviceType, err)
	}

	return products, nil
}

// FilterByExactPrice defines the logic of the getting and processing the products' sample
// from the markets' responses constrained by the markets' filters and the products that
// have got the exactest prices to the client's price.
func (filter *ProductsFilter) FilterByExactPrice(ctx echo.Context, product entities.ProductRequest, exactPrice int) ([]entities.Product, error) {
	const serviceType = "filter.service.filter_by_exact_price"
	products, err := filter.api.GetProductsByExactPrice(ctx, product, exactPrice)

	if err != nil {
		return nil, fmt.Errorf("error of the %v: %v", serviceType, err)
	}

	return products, nil
}
