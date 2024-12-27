package entities

// ProductRequest defines the request data from the client to this service.
type ProductRequest struct {
	ProductName string
	Sample      int
	Amount      string
	Markets     []Market
}

func NewProductRequest() ProductRequest {
	return ProductRequest{
		Markets: make([]Market, 0, 15),
	}
}
