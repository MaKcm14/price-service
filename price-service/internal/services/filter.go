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

func NewFilter(log *slog.Logger) ProductsFilter {
	return ProductsFilter{
		logger: log,
	}
}

// switchMarketApi swithches the context of market api according to the client's request.
func (filter *ProductsFilter) switchMarketApi(market entities.Market) error {
	if market == entities.Wildberries {
		filter.api = api.NewWildberriesAPI(filter.logger, 1)
		return nil
	}
	return ErrChooseMarket
}

// FilterByMarkets defines the logic of the getting and processing the products' sample
// from the markets' responses filtered only by markets.
func (filter ProductsFilter) FilterByMarkets(ctx echo.Context, product entities.ProductRequest) ([]entities.ProductSample, error) {
	const serviceType = "filter.service.filter-by-markets"
	var products = make([]entities.ProductSample, 0, 1000)

	for _, market := range product.Markets {
		if err := filter.switchMarketApi(market); err != nil {
			filter.logger.Warn(fmt.Sprintf("error of the %v: %v", serviceType, err))
			continue
		}

		sample, err := filter.api.GetProducts(ctx, product)

		if err != nil {
			filter.logger.Warn(fmt.Sprintf("error of the %v: %v", serviceType, err))
			continue
		}

		products = append(products, sample)
	}

	if len(products) == 0 {
		return nil, ErrGettingProducts
	}

	return products, nil
}

// FilterByPriceRange defines the logic of the getting and processing the products' sample
// from the markets' responses constrained by the markets' filters and two boundaries of
// the price range.
func (filter ProductsFilter) FilterByPriceRange(ctx echo.Context, product entities.ProductRequest, priceDown int, priceUp int) ([]entities.ProductSample, error) {
	const serviceType = "filter.service.filter-by-price-range"
	var products = make([]entities.ProductSample, 0, 1000)

	for _, market := range product.Markets {
		if err := filter.switchMarketApi(market); err != nil {
			filter.logger.Warn(fmt.Sprintf("error of the %v: %v", serviceType, err))
			continue
		}

		sample, err := filter.api.GetProductsWithPriceRange(ctx, product, priceDown, priceUp)

		if err != nil {
			filter.logger.Warn(fmt.Sprintf("error of the %v: %v", serviceType, err))
			continue
		}

		products = append(products, sample)
	}

	if len(products) == 0 {
		return nil, ErrGettingProducts
	}

	return products, nil
}

// FilterBestPrice defines the logic of the getting and processing the products' sample
// from the markets' responses contrained by the markets' filters and the minimal price of the sample.
func (filter ProductsFilter) FilterByBestPrice(ctx echo.Context, product entities.ProductRequest) ([]entities.ProductSample, error) {
	const serviceType = "filter.service.filter-by-best-price"
	var products = make([]entities.ProductSample, 0, 1000)

	for _, market := range product.Markets {
		if err := filter.switchMarketApi(market); err != nil {
			filter.logger.Warn(fmt.Sprintf("error of the %v: %v", serviceType, err))
			continue
		}

		sample, err := filter.api.GetProductsWithBestPrice(ctx, product)

		if err != nil {
			filter.logger.Warn(fmt.Sprintf("error of the %v: %v", serviceType, err))
			continue
		}

		products = append(products, sample)
	}

	if len(products) == 0 {
		return nil, ErrGettingProducts
	}

	return products, nil
}

// FilterByExactPrice defines the logic of the getting and processing the products' sample
// from the markets' responses constrained by the markets' filters and the products that
// have got the exactest prices to the client's price.
func (filter ProductsFilter) FilterByExactPrice(ctx echo.Context, product entities.ProductRequest, exactPrice int) ([]entities.ProductSample, error) {
	const serviceType = "filter.service.filter-by-exact-price"
	var products = make([]entities.ProductSample, 0, 1000)

	for _, market := range product.Markets {
		if err := filter.switchMarketApi(market); err != nil {
			filter.logger.Warn(fmt.Sprintf("error of the %v: %v", serviceType, err))
			continue
		}

		sample, err := filter.api.GetProductsWithExactPrice(ctx, product, exactPrice)

		if err != nil {
			filter.logger.Warn(fmt.Sprintf("error of the %v: %v", serviceType, err))
			continue
		}

		products = append(products, sample)
	}

	if len(products) == 0 {
		return nil, ErrGettingProducts
	}

	return products, nil
}
