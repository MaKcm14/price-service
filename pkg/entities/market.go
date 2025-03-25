package entities

type Market int

// MarketView defines the data structure of the concrete market.
type MarketView struct {
	MarketName  string `json:"name"`
	Designation string `json:"emoji"`
}

// Markets defines the data structure of the supported markets.
type SupportedMarkets struct {
	Markets []MarketView `json:"markets"`
}

const (
	Wildberries Market = iota
	MegaMarket
	NotExists
)

// GetMarkets returns the current supported markets.
func GetSupportedMarkets() SupportedMarkets {
	return SupportedMarkets{
		Markets: []MarketView{
			{"Wildberries", "🌸"},
			{"Megamarket", "🛍️"},
		},
	}
}
