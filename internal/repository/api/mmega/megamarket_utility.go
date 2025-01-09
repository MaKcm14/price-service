package mmega

import (
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	"github.com/MaKcm14/best-price-service/price-service/internal/entities/dto"
	"github.com/MaKcm14/best-price-service/price-service/internal/repository/api"
)

// URL paths' consts.
const (
	megaMarketOrigin      = "https://megamarket.ru"
	megaMarketOpenApiPath = "/catalog/page-"
	priceRangeKey         = "88C83F68482F447C9F4E401955196697"
)

// URL query params' consts.
const (
	priceRangeID = "filters"
	sortID       = "sort"
)

const (
	// minValue is the min amount of products that can be got from the service in one sample.
	minValue = 15
)

type (
	megaMarketProductOffer struct {
		MerchantName string `json:"merchantName"`
	}

	megaMarketProductDecription struct {
		Title          string `json:"title"`
		TitleImageLink string `json:"titleImage"`
		URL            string `json:"webUrl"`
		Brand          string `json:"brand"`
	}

	megaMarketProduct struct {
		Goods      megaMarketProductDecription `json:"goods"`
		Price      int                         `json:"price"`
		FinalPrice int                         `json:"finalPrice"`
		Offer      megaMarketProductOffer      `json:"favoriteOffer"`
	}

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
		path += fmt.Sprintf("&%s=%s", priceRangeID, v.getPriceRangeURLView(priceRange))
	}

	return path
}

// getPriceRangeURLView returns the correct URL-view of the price-range filter.
func (v megaMarketViewer) getPriceRangeURLView(priceRange string) string {
	prices := strings.Split(priceRange, " ")
	priceRangeObj := fmt.Sprintf("{\"%s\":{\"min\":%s,\"max\":%s}}", priceRangeKey, prices[0], prices[1])
	urlPriceRange := ""

	for _, elem := range priceRangeObj {
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

// getSortParamURLView returns the mapped value of the sort parameter specified for megamarket-service.
func (v megaMarketViewer) getSortParamURLView(sort string) string {
	if sort == "priceup" {
		return "1"
	} else if sort == "pricedown" {
		return "2"
	} else if sort == "newly" {
		return "5"
	}
	return "0"
}
