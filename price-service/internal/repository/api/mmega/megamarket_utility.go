package mmega

import (
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"strings"

	"github.com/MaKcm14/best-price-service/price-service/internal/entities"
	"github.com/MaKcm14/best-price-service/price-service/internal/entities/dto"
	"github.com/MaKcm14/best-price-service/price-service/internal/repository/api"
	"github.com/anaskhan96/soup"
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

// parsing's consts.
const (
	productContainerClassName = "catalog-items-list__container"
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

// parseProductPrices parses the current html page for getting the product prices.
func (p megaMarketParser) parseProductPrices(html string) []entities.Price {
	return nil
}

// parseProductSuppliers parses the current html page for getting the product suppliers.
func (p megaMarketParser) parseProductSuppliers(html string) []string {
	return nil
}

// getPrice returns the int price view of the parsed product's price.
func (p megaMarketParser) getPrice(respPrice string) int {
	price := ""

	for _, elem := range respPrice {
		if elem >= 48 && elem <= 57 {
			price += string(elem)
		}
	}

	res, _ := strconv.Atoi(price)

	return res
}

// TODO: divide it. This is a test function that can be deleted in the future versions.
func (p megaMarketParser) parseProducts(html string) []entities.Product {
	var products = make([]entities.Product, 0, 50)

	for _, tag := range soup.HTMLParse(html).FindAll("a", "data-test", "product-image-link") {
		link := tag.Find("img")

		products = append(products, entities.Product{
			Name: link.Attrs()["alt"],
			Links: entities.ProductLink{
				URL:       megaMarketOrigin + strings.Join(strings.Split(tag.Attrs()["href"], " "), "%20"),
				ImageLink: link.Attrs()["src"],
			},
		})

		fmt.Println(products[len(products)-1])
	}

	discPrices := make([]int, 0, len(products))
	priceTags := soup.HTMLParse(html).FindAll("div", "class", "item-price")
	for i := 0; i != len(priceTags); i++ {
		var tag = priceTags[i]

		tag = tag.Find("div", "data-test", "product-price")

		price := p.getPrice(tag.Text())

		discPrices = append(discPrices, price)
		fmt.Println(price)
	}

	basePrices := make([]int, 0, len(products))
	for _, tag := range soup.HTMLParse(html).FindAll("div", "class", "item-price-discount__number") {
		if tag.Pointer != nil {
			price := p.getPrice(tag.Text())
			basePrices = append(basePrices, price)
			fmt.Println(price)
		}
	}

	return products
}
