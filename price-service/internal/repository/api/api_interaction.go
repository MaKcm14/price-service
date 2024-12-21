package api

import (
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
		wildb:  NewWildberriesAPI(log, 10000),
	}
}

// GetProducts defines getting the products from the needed markets, that set in ProductRequest DTO.
func (api MarketsApi) GetProducts(product entities.ProductRequest) {
	// DEBUG:
	api.wildb.GetProducts(product)
	// TODO: check this
}
