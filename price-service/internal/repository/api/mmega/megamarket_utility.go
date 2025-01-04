package mmega

import (
	"fmt"
	"log/slog"
	"net/url"

	"github.com/MaKcm14/best-price-service/price-service/internal/entities/dto"
	"github.com/MaKcm14/best-price-service/price-service/internal/repository/api"
)

// URL paths' consts.
const (
	megaMarketOrigin      = "https://megamarket.ru"
	megaMarketOpenApiPath = "/catalog/page-"

	//TODO: realocate it to the .env:
	priceRangeKey = "88C83F68482F447C9F4E401955196697"
)

// URL query params' consts.
const (
	priceRangeID = "filters"
	sortID       = "sort"
)

type (
	// megaMarketViewer defines the specific view of the filters and url for MegaMarket.
	megaMarketViewer struct {
		converter api.URLConverter
	}

	// megaMarketParser defines the logic of the parsing.
	megaMarketParser struct {
		logger *slog.Logger
	}
)

// getOpenApiURL returns the full open API-url with the set filters.
// It uses with the origin "https://megamarket.ru".
func (v megaMarketViewer) getOpenApiURL(request dto.ProductRequest, filter []string) string {
	path := megaMarketOrigin + megaMarketOpenApiPath +
		fmt.Sprintf("%d/?q=%s", request.Sample, url.QueryEscape(request.Query))

	filters := v.converter.GetFilters(filter)
	path += fmt.Sprintf("#?%s=%s", sortID, filters[sortID])

	if priceRange, flagExist := filters[priceRangeID]; flagExist {
		path += fmt.Sprintf("&%s=%s", priceRangeID, priceRange)
	}

	return path
}

func (v megaMarketViewer) getPriceRangeView(priceDown int, priceUp int) string {
	priceRange := fmt.Sprintf("{\"%s\":{\"min\":%d,\"max\":%d}}", priceRangeKey, priceDown, priceUp)
	urlPriceRange := ""

	for _, elem := range priceRange {
		if string(elem) == "{" {
			urlPriceRange += "%7B"
		} else if string(elem) == "}" {
			urlPriceRange += "%7D"
		} else if string(elem) == ":" {
			urlPriceRange += "%3A"
		} else if string(elem) == "," {
			urlPriceRange += "%2C"
		} else {
			urlPriceRange += string(elem)
		}
	}

	return urlPriceRange
}

func (v megaMarketViewer) getSortParamView(sort string) string {
	if sort == "priceup" {
		return "1"
	} else if sort == "pricedown" {
		return "2"
	} else if sort == "newly" {
		return "5"
	}
	return "0"
}
