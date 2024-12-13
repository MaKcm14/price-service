package services

type (
	MarketFilterAdapter interface {
		FilterByMarkets()
	}

	PriceFilterAdapter interface {
		FilterByPriceRange()
		FilterByBestPrice()
		FilterByExactPrice()
	}

	Filter interface {
		MarketFilterAdapter
		PriceFilterAdapter
	}
)
