package entities

type Product struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	URL         string `json:"url"`
	Market      string `json:"market"`
	Image       []byte `json:"image"`
}
