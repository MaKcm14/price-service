package chttp

import (
	"strings"

	"github.com/MaKcm14/best-price-service/price-service/internal/entities"
)

// ProductResponse defines the response data.
type ProductResponse struct {
	Samples map[string]entities.ProductSample `json:"samples"`
}

func NewProductResponse(samples []entities.ProductSample) ProductResponse {
	marketSamples := make(map[string]entities.ProductSample)

	for _, sample := range samples {
		marketSamples[strings.ToLower(sample.Market)] = sample
	}

	return ProductResponse{
		Samples: marketSamples,
	}
}
