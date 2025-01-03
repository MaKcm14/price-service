package entities

type ProductLink struct {
	URL       string `json:"url"`
	ImageLink string `json:"image_link"`
}

type Product struct {
	Name     string      `json:"name"`
	Brand    string      `json:"brand"`
	Price    Price       `json:"price"`
	Links    ProductLink `json:"related_links"`
	Supplier string      `json:"supplier"`
}

// ProductSample defines the sample of the products from the one market.
type ProductSample struct {
	Products   []Product `json:"products"`
	SampleLink string    `json:"main_products_sample"`
	Market     string    `json:"market"`
}

func NewProductSample(products []Product, sampleLink string, sampleMarket Market) ProductSample {
	var market string

	if sampleMarket == Wildberries {
		market = "Wildberries"
	} else if sampleMarket == Ozon {
		market = "Ozon"
	} else if sampleMarket == MegaMarket {
		market = "MegaMarket"
	}

	return ProductSample{
		Products:   products,
		SampleLink: sampleLink,
		Market:     market,
	}
}
