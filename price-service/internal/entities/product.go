package entities

type ProductMetaData struct {
	URL       string `json:"url"`
	ImageLink string `json:"image_link"`
}

type Product struct {
	Name     string          `json:"name"`
	Brand    string          `json:"brand"`
	Price    Price           `json:"price"`
	MetaData ProductMetaData `json:"meta_data"`
	Market   string          `json:"market"`
	Supplier string          `json:"supplier"`
}
