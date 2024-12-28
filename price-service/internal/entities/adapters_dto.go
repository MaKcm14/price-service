package entities

// ProductRequest defines the request data from the client to this service.
type ProductRequest struct {
	Query       string
	Sample      int
	Amount      string
	Sort        string
	FlagNoImage bool
	Markets     []Market
}

func NewProductRequest() ProductRequest {
	return ProductRequest{
		Markets: make([]Market, 0, 15),
	}
}
