package chttp

import (
	"strings"

	"github.com/MaKcm14/price-service/pkg/entities"
)

// ProductResponse defines the response data.
type ProductResponse struct {
	Samples map[string]entities.ProductSample `json:"samples"`
}

// Header defines the header data.
type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// ExtraHeaders defines the extra-headers data for the async call.
type ExtraHeaders struct {
	Headers []Header `json:"headers"`
}

func NewExtraHeaders() *ExtraHeaders {
	return &ExtraHeaders{
		Headers: make([]Header, 0, 100),
	}
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
