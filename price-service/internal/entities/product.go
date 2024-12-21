package entities

type Product struct {
	Name          string `json:"name"`
	Brand         string `json:"brand"`
	BasePrice     int    `json:"base_price"`
	DiscountPrice int    `json:"discount_price"`
	Discount      int    `json:"discount"`
	URL           string `json:"url"`
	Market        string `json:"market"`
	Supplier      string `json:"supplier"`
	Image         []byte `json:"image"`
}
