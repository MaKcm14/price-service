package entities

// ProductRequest defines the request data from the client to this service.
type ProductRequest struct {
	ProductName string
	Sample      int
	Markets     []Market
	Client      ClientType
}

type ProductResponse struct {
	Products         []Product `json:"products"`
	ParentSampleLink string    `json:"main_products_sample"`
}

func NewProductRequest() ProductRequest {
	return ProductRequest{
		Markets: make([]Market, 0, 10),
	}
}

func NewProductResponse(products []Product, sampleLink string) ProductResponse {
	return ProductResponse{
		Products:         products,
		ParentSampleLink: sampleLink,
	}
}
