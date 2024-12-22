package api

import (
	"fmt"
	"log/slog"

	"github.com/MaKcm14/best-price-service/price-service/internal/entities"
)

const (
	Wildberries string = "Wildberries"
)

// MarketsApi is the common type that combines all logic of marketplaces' interaction.
type MarketsApi struct {
	logger *slog.Logger
	wildb  WildberriesAPI
}

func NewMarketsApi(log *slog.Logger) MarketsApi {
	return MarketsApi{
		logger: log,
		wildb:  NewWildberriesAPI(log),
	}
}

// GetProducts defines getting the products from the needed markets, that set in ProductRequest DTO.
func (api MarketsApi) GetProducts(product entities.ProductRequest) ([]entities.Product, error) {
	const serviceType = "api.service"
	var products = make([]entities.Product, 0, 10000)

	wildProd, err := api.wildb.GetProducts(product)

	if err != nil {
		return nil, fmt.Errorf("error of the %v: %v", serviceType, err)
	}
	products = append(products, wildProd...)

	return products, nil
}

// GetProductsByPriceRange defines getting the products from the needed markets.
func (api MarketsApi) GetProductsByPriceRange(product entities.ProductRequest, priceDown, priceUp int) ([]entities.Product, error) {
	const serviceType = "api.service"
	var products = make([]entities.Product, 0, 10000)

	wildProd, err := api.wildb.GetProductsByPriceRange(product, priceDown, priceUp)

	if err != nil {
		return nil, fmt.Errorf("error of the %v: %v", serviceType, err)
	}
	products = append(products, wildProd...)

	return products, nil
}

// GetProductsByExactPrice defines getting the products from the needed markets
// that have the price in range [exactPrice, exactPrice + 10% off exactPrice].
func (api MarketsApi) GetProductsByExactPrice(product entities.ProductRequest, exactPrice int) ([]entities.Product, error) {
	const serviceType = "api.service"
	var products = make([]entities.Product, 0, 10000)

	wildProd, err := api.wildb.GetProductsByExactPrice(product, exactPrice)

	if err != nil {
		return nil, fmt.Errorf("error of the %v: %v", serviceType, err)
	}
	products = append(products, wildProd...)

	return products, nil
}

// GetProductsByBestPrice defines getting the products from the needed markets
// that have the min price.
func (api MarketsApi) GetProductsByBestPrice(product entities.ProductRequest) ([]entities.Product, error) {
	const serviceType = "api.service"
	var products = make([]entities.Product, 0, 10000)

	wildProd, err := api.wildb.GetProductsByBestPrice(product)

	if err != nil {
		return nil, fmt.Errorf("error of the %v: %v", serviceType, err)
	}
	products = append(products, wildProd...)

	return products, nil
}
