package dto

import "github.com/MaKcm14/price-service/pkg/entities"

const (
	PopularSort   SortType = "popular"
	PriceUpSort   SortType = "priceup"
	PriceDownSort SortType = "pricedown"
	NewlySort     SortType = "newly"
	RateSort      SortType = "rate"
)

type SortType string

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

	Async   bool
	Headers map[string]string

	PriceRange PriceRangeRequest
	ExactPrice int
}

func NewProductRequest() ProductRequest {
	return ProductRequest{
		Markets: make([]entities.Market, 0, 15),
		Headers: make(map[string]string),
		Async:   false,
	}
}
