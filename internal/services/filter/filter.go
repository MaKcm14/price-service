package filter

import (
	"fmt"
	"log/slog"

	"github.com/labstack/echo/v4"

	"github.com/MaKcm14/price-service/internal/entities/dto"
	"github.com/MaKcm14/price-service/internal/services"
	"github.com/MaKcm14/price-service/pkg/entities"
)

type filterType int

const (
	priceRangeFilter filterType = iota
	exactPriceFilter
	bestPriceFilter
	commonFilter
)

type ProductsFilter struct {
	logger     *slog.Logger
	marketsApi map[entities.Market]services.ApiInteractor
}

func New(log *slog.Logger, markets map[entities.Market]services.ApiInteractor) ProductsFilter {
	return ProductsFilter{
		logger:     log,
		marketsApi: markets,
	}
}

// getMarketApi returns the requested marketApi wrapped in the ApiInteractor.
func (p *ProductsFilter) getMarketApi(market entities.Market) (services.ApiInteractor, error) {
	marketApi, flagExist := p.marketsApi[market]

	if !flagExist {
		return nil, services.ErrMarketApi
	}

	return marketApi, nil
}

// filter defines the main filter logic which defines the flow of control according to the set filter type.
func (p *ProductsFilter) filter(ctx echo.Context, request dto.ProductRequest, serviceType string, filter filterType) ([]entities.ProductSample, error) {
	var products = make([]entities.ProductSample, 0, 1000)

	for _, market := range request.Markets {
		var sample entities.ProductSample
		var err error

		marketApi, err := p.getMarketApi(market)

		if err != nil {
			p.logger.Warn(fmt.Sprintf("error of the %v: %v", serviceType, err))
			continue
		}

		if filter == commonFilter {
			sample, err = marketApi.GetProducts(ctx, request)
		} else if filter == priceRangeFilter {
			sample, err = marketApi.GetProductsWithPriceRange(ctx, request)
		} else if filter == exactPriceFilter {
			sample, err = marketApi.GetProductsWithExactPrice(ctx, request)
		} else if filter == bestPriceFilter {
			sample, err = marketApi.GetProductsWithBestPrice(ctx, request)
		} else {
			continue
		}

		if err != nil {
			p.logger.Warn(fmt.Sprintf("error of the %v: %v", serviceType, err))
			continue
		}

		products = append(products, sample)
	}

	if len(products) == 0 {
		return nil, services.ErrGettingProducts
	}

	return products, nil
}

// FilterByMarkets defines the logic of the getting and processing the products' sample
// from the markets' responses filtered only by markets and non-specified parameters.
func (p ProductsFilter) FilterByMarkets(ctx echo.Context, request dto.ProductRequest) ([]entities.ProductSample, error) {
	const serviceType = "filter.service.filter-by-markets"
	return p.filter(ctx, request, serviceType, commonFilter)
}

// FilterByPriceRange defines the logic of the getting and processing the products' sample
// from the markets' responses constrained by the markets' filters and two boundaries of
// the price range.
func (p ProductsFilter) FilterByPriceRange(ctx echo.Context, request dto.ProductRequest) ([]entities.ProductSample, error) {
	const serviceType = "filter.service.filter-by-price-range"
	return p.filter(ctx, request, serviceType, priceRangeFilter)
}

// FilterBestPrice defines the logic of the getting and processing the products' sample
// from the markets' responses contrained by the markets' filters and the minimal price of the sample.
func (p ProductsFilter) FilterByBestPrice(ctx echo.Context, request dto.ProductRequest) ([]entities.ProductSample, error) {
	const serviceType = "filter.service.filter-by-best-price"
	return p.filter(ctx, request, serviceType, bestPriceFilter)
}

// FilterByExactPrice defines the logic of the getting and processing the products' sample
// from the markets' responses constrained by the markets' filters and the products that
// have got the exactest prices to the client's price.
func (p ProductsFilter) FilterByExactPrice(ctx echo.Context, request dto.ProductRequest) ([]entities.ProductSample, error) {
	const serviceType = "filter.service.filter-by-exact-price"
	return p.filter(ctx, request, serviceType, bestPriceFilter)
}
