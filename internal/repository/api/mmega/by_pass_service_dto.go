package mmega

import "strings"

type (
	byPassPriceRangeFilter struct {
		PriceDown string `json:"price_down"`
		PriceUp   string `json:"price_up"`
	}

	byPassServiceRequest struct {
		Query              string                 `json:"query"`
		Sample             string                 `json:"sample"`
		Sort               string                 `json:"sort"`
		FlagShowNotAvail   bool                   `json:"show_not_available"`
		FlagPriceFilterSet bool                   `json:"is_price_filter_set"`
		PriceRange         byPassPriceRangeFilter `json:"price_filter"`
	}
)

func newByPassServiceRequest(query, sample, sort string, filters map[string]string) byPassServiceRequest {
	byPassRequest := byPassServiceRequest{
		Query:            query,
		Sample:           sample,
		Sort:             sort,
		FlagShowNotAvail: false,
	}

	if priceParam, flagExist := filters[priceRangeID]; flagExist {
		prices := strings.Split(priceParam, " ")

		byPassRequest.FlagPriceFilterSet = true

		byPassRequest.PriceRange = byPassPriceRangeFilter{
			PriceDown: prices[0],
			PriceUp:   prices[1],
		}
	}

	return byPassRequest
}
