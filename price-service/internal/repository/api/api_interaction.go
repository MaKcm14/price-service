package api

import (
	"fmt"
	"log/slog"

	"github.com/labstack/echo/v4"

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
		wildb:  NewWildberriesAPI(log, 1),
	}
}

func IsConnectionClosed(ctx echo.Context) bool {
	select {
	case <-ctx.Request().Context().Done():
		return true

	default:
		return false
	}
}

// GetProducts defines getting the products from the needed markets, that set in ProductRequest DTO.
func (api MarketsApi) GetProducts(ctx echo.Context, product entities.ProductRequest) ([]entities.ProductResponse, error) {
	const serviceType = "api.service"
	var products = make([]entities.ProductResponse, 0, 1000)

	wildProd, err := api.wildb.GetProducts(ctx, product)

	if err != nil {
		return nil, fmt.Errorf("error of the %v: %v", serviceType, err)
	}
	products = append(products, wildProd)

	return products, nil
}

// GetProductsByPriceRange defines getting the products from the needed markets.
func (api MarketsApi) GetProductsByPriceRange(ctx echo.Context, product entities.ProductRequest, priceDown, priceUp int) ([]entities.ProductResponse, error) {
	const serviceType = "api.service"
	var products = make([]entities.ProductResponse, 0, 10000)

	wildProd, err := api.wildb.GetProductsByPriceRange(ctx, product, priceDown, priceUp)

	if err != nil {
		return nil, fmt.Errorf("error of the %v: %v", serviceType, err)
	}
	products = append(products, wildProd)

	return products, nil
}

// GetProductsByExactPrice defines getting the products from the needed markets
// that have the price in range [exactPrice, exactPrice + 10% off exactPrice].
func (api MarketsApi) GetProductsByExactPrice(ctx echo.Context, product entities.ProductRequest, exactPrice int) ([]entities.ProductResponse, error) {
	const serviceType = "api.service"
	var products = make([]entities.ProductResponse, 0, 10000)

	wildProd, err := api.wildb.GetProductsByExactPrice(ctx, product, exactPrice)

	if err != nil {
		return nil, fmt.Errorf("error of the %v: %v", serviceType, err)
	}
	products = append(products, wildProd)

	return products, nil
}

// GetProductsByBestPrice defines getting the products from the needed markets
// that have the min price.
func (api MarketsApi) GetProductsByBestPrice(ctx echo.Context, product entities.ProductRequest) ([]entities.ProductResponse, error) {
	const serviceType = "api.service"
	var products = make([]entities.ProductResponse, 0, 10000)

	wildProd, err := api.wildb.GetProductsByBestPrice(ctx, product)

	if err != nil {
		return nil, fmt.Errorf("error of the %v: %v", serviceType, err)
	}
	products = append(products, wildProd)

	return products, nil
}
