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

func NewProductsFilter(log *slog.Logger) ProductsFilter {
	return ProductsFilter{
		logger: log,
	}
}

// switchMarketApi swithches the context of market api according to the client's request.
func (p *ProductsFilter) switchMarketApi(market entities.Market) error {
	if market == entities.Wildberries {
		p.api = api.NewWildberriesAPI(p.logger, 1.4)
		return nil
	}
	return ErrChooseMarket
}

// FilterByMarkets defines the logic of the getting and processing the products' sample
// from the markets' responses filtered only by markets and non-specified parameters.
func (p ProductsFilter) FilterByMarkets(ctx echo.Context, request entities.ProductRequest) ([]entities.ProductSample, error) {
	const serviceType = "filter.service.filter-by-markets"
	var products = make([]entities.ProductSample, 0, 1000)

	for _, market := range request.Markets {
		if err := p.switchMarketApi(market); err != nil {
			p.logger.Warn(fmt.Sprintf("error of the %v: %v", serviceType, err))
			continue
		}

		sample, err := p.api.GetProducts(ctx, request)

		if err != nil {
			p.logger.Warn(fmt.Sprintf("error of the %v: %v", serviceType, err))
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
func (p ProductsFilter) FilterByPriceRange(ctx echo.Context, request entities.ProductRequest, priceDown int, priceUp int) ([]entities.ProductSample, error) {
	const serviceType = "filter.service.filter-by-price-range"
	var products = make([]entities.ProductSample, 0, 1000)

	for _, market := range request.Markets {
		if err := p.switchMarketApi(market); err != nil {
			p.logger.Warn(fmt.Sprintf("error of the %v: %v", serviceType, err))
			continue
		}

		sample, err := p.api.GetProductsWithPriceRange(ctx, request, priceDown, priceUp)

		if err != nil {
			p.logger.Warn(fmt.Sprintf("error of the %v: %v", serviceType, err))
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
func (p ProductsFilter) FilterByBestPrice(ctx echo.Context, request entities.ProductRequest) ([]entities.ProductSample, error) {
	const serviceType = "filter.service.filter-by-best-price"
	var products = make([]entities.ProductSample, 0, 1000)

	for _, market := range request.Markets {
		if err := p.switchMarketApi(market); err != nil {
			p.logger.Warn(fmt.Sprintf("error of the %v: %v", serviceType, err))
			continue
		}

		sample, err := p.api.GetProductsWithBestPrice(ctx, request)

		if err != nil {
			p.logger.Warn(fmt.Sprintf("error of the %v: %v", serviceType, err))
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
func (p ProductsFilter) FilterByExactPrice(ctx echo.Context, request entities.ProductRequest, exactPrice int) ([]entities.ProductSample, error) {
	const serviceType = "filter.service.filter-by-exact-price"
	var products = make([]entities.ProductSample, 0, 1000)

	for _, market := range request.Markets {
		if err := p.switchMarketApi(market); err != nil {
			p.logger.Warn(fmt.Sprintf("error of the %v: %v", serviceType, err))
			continue
		}

		sample, err := p.api.GetProductsWithExactPrice(ctx, request, exactPrice)

		if err != nil {
			p.logger.Warn(fmt.Sprintf("error of the %v: %v", serviceType, err))
			continue
		}

		products = append(products, sample)
	}

	if len(products) == 0 {
		return nil, ErrGettingProducts
	}

	return products, nil
}
