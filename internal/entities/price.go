package entities

type Price struct {
	BasePrice     int `json:"base_price"`
	DiscountPrice int `json:"discount_price"`
	Discount      int `json:"discount"`
}

func NewPrice(basePrice, discPrice int) Price {
	if basePrice <= 0 || discPrice <= 0 || basePrice < discPrice {
		return Price{}
	}

	discount := (basePrice - discPrice) * 100
	discount /= basePrice

	return Price{
		BasePrice:     basePrice,
		DiscountPrice: discPrice,
		Discount:      discount,
	}
}
