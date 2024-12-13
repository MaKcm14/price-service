package services

import "github.com/MaKcm14/best-price-service/price-service/internal/entities"

// TODO: delete the json's tags
type ProductRequest struct {
	ProductName string
	Sample      int
	Markets     []entities.Market
}

func NewProductRequest() ProductRequest {
	return ProductRequest{
		Markets: make([]entities.Market, 0, 10),
	}
}
