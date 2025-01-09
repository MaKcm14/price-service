package dto

import "github.com/MaKcm14/best-price-service/price-service/internal/entities"

type SortType string

const (
	PopularSort   SortType = "popular"
	PriceUpSort   SortType = "priceup"
	PriceDownSort SortType = "pricedown"
	NewlySort     SortType = "newly"
	RateSort      SortType = "rate"
)

// ProductRequest defines the request data from the client to this service.
type ProductRequest struct {
	Query       string
	Sample      int
	Amount      string
	Sort        SortType
	FlagNoImage bool
	Markets     []entities.Market
}

func NewProductRequest() ProductRequest {
	return ProductRequest{
		Markets: make([]entities.Market, 0, 15),
	}
}
