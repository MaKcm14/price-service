package entities

type Market int

type MarketView struct {
	MarketName  string `json:"name"`
	Designation string `json:"emoji"`
}

const (
	Wildberries Market = iota
	MegaMarket
	NotExists
)

// GetMarkets returns the current supported markets.
func GetMarkets() map[string][]MarketView {
	return map[string][]MarketView{
		"markets": {
			{"Wildberries", "🌸"},
			{"Megamarket", "🛍️"},
		},
	}
}
