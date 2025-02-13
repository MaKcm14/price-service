package dto

import "github.com/MaKcm14/price-service/pkg/entities"

type SortType string

const (
	PopularSort   SortType = "popular"
	PriceUpSort   SortType = "priceup"
	PriceDownSort SortType = "pricedown"
	NewlySort     SortType = "newly"
	RateSort      SortType = "rate"
)

// PriceRangeRequest defines the request data specially for price-range filter.
type PriceRangeRequest struct {
	PriceDown int
	PriceUp   int
}

// ProductRequest defines the request data from the client to this service.
type ProductRequest struct {
	Query       string
	Sample      int
	Amount      string
	Sort        SortType
	FlagNoImage bool
	Markets     []entities.Market

	PriceRange PriceRangeRequest
	ExactPrice int
}

func NewProductRequest() ProductRequest {
	return ProductRequest{
		Markets: make([]entities.Market, 0, 15),
	}
}
