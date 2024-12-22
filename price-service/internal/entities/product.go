package entities

type Product struct {
	Name      string `json:"name"`
	Brand     string `json:"brand"`
	Price     Price  `json:"price"`
	URL       string `json:"url"`
	Market    string `json:"market"`
	Supplier  string `json:"supplier"`
	ImageLink string `json:"image_link"`
}
