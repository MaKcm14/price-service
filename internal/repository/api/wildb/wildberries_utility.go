package wildb

import (
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	"github.com/anaskhan96/soup"

	"github.com/MaKcm14/price-service/internal/entities/dto"
	"github.com/MaKcm14/price-service/internal/repository/api"
)

// URL paths' consts.
const (
	wildberriesGeoString   = "appType=1&curr=rub&dest=-1257786&hide_dtype=10&lang=ru"
	searchAPIPath          = "https://search.wb.ru/exactmatch/ru/common/v9/search"
	wildberriesOpenAPIPath = "https://www.wildberries.ru/catalog/0/search.aspx"
)

// URL query params' consts.
const (
	priceRangeID = "priceU"
	sortID       = "sort"
	searchID     = "search"
	queryID      = "query"
	pageID       = "page"
)

// parsing's consts.
const (
	mainProductsTagName       = "article"
	imageClassName            = "j-thumbnail"
	productContainerClassName = "product-card-list"
)

type (
	wildberriesProduct struct {
		ID       int    `json:"id"`
		Brand    string `json:"brand"`
		Name     string `json:"name"`
		Supplier string `json:"supplier"`
		Sizes    []struct {
			Price struct {
				Basic int
				Total int
			} `json:"price"`
		} `json:"sizes"`
	}

	// wildberriesViewer defines the logic of the queries' parameters format.
	wildberriesViewer struct {
		converter api.URLConverter
	}

	// wildberriesParser defines the logic of parsing the wildberries' service sources.
	wildberriesParser struct {
		logger *slog.Logger
	}
)

// getOpenApiPath returns the correct URL's path for wildberries open API.
// It uses with the origin "https://www.wildberries.ru".
func (v wildberriesViewer) getOpenApiPath(request dto.ProductRequest, filters []string) string {
	var path string
	filtersURL := v.converter.GetFilters(filters)

	if len(filtersURL) == 0 {
		return ""
	}

	path += fmt.Sprintf("%s=%d", pageID, request.Sample)
	path += fmt.Sprintf("&%s=%s", sortID, filtersURL[sortID])

	if priceRange, flagExist := filtersURL[priceRangeID]; flagExist {
		path += fmt.Sprintf("&%s=%s", priceRangeID, priceRange)
	}

	path += fmt.Sprintf("&%s=%s", searchID, strings.Join(strings.Split(request.Query, " "), "+"))

	return path
}

// getHiddenApiPath returns the correct URL's path for wildberries hidden API.
// It uses with the origin "https://search.wb.ru".
func (v wildberriesViewer) getHiddenApiPath(request dto.ProductRequest, filters []string) string {
	var path string
	filtersURL := v.converter.GetFilters(filters)

	if len(filtersURL) == 0 {
		return ""
	}

	path += fmt.Sprintf("%s=%d", pageID, request.Sample)

	if priceRange, flagExist := filtersURL[priceRangeID]; flagExist {
		path += fmt.Sprintf("&%s=%s", priceRangeID, priceRange)
	}

	path += fmt.Sprintf("&%s=%s", queryID, url.QueryEscape(request.Query))

	path += fmt.Sprintf("&resultset=catalog&%s=%s", sortID, filtersURL["sort"])
	path += "&spp=30&suppressSpellcheck=false"

	return path
}

func (v wildberriesViewer) getPriceRangeView(priceDown int, priceUp int) string {
	return fmt.Sprintf("%v00;%v00", priceDown, priceUp)
}

func (v wildberriesViewer) getProductCatalogLink(productID int) string {
	return fmt.Sprintf("https://www.wildberries.ru/catalog/%d/detail.aspx", productID)
}

// getHiddenApiURL returns the full hidden API-url for the "search.wb.ru" with the set filters.
func (p wildberriesViewer) getHiddenApiURL(request dto.ProductRequest, filters []string) string {
	return fmt.Sprintf("%s?"+
		"ab_testing=false&%s&%s", searchAPIPath,
		wildberriesGeoString, p.getHiddenApiPath(request, filters))
}

// getOpenApiURL returns the full open API-url for the "www.wildberries.ru" with the set filters.
func (p wildberriesViewer) getOpenApiURL(request dto.ProductRequest, filters []string) string {
	return fmt.Sprintf("%s?%s", wildberriesOpenAPIPath,
		p.getOpenApiPath(request, filters))
}

// parseImageLinks parses the image links for the products from the current html-page.
func (p wildberriesParser) parseImageLinks(html string) []string {
	const serviceType = "wildberries.service.image-links-getter"

	if !strings.Contains(html, mainProductsTagName) {
		p.logger.Warn(fmt.Sprintf("error of the %v: %v: images couldn't be load", serviceType, api.ErrServiceResponse))
		return nil
	}

	var imageLinks = make([]string, 0, 100)

	for _, tag := range soup.HTMLParse(html).FindAll(mainProductsTagName) {
		link := tag.Find("img", "class", imageClassName)

		if link.Pointer != nil {
			imageLinks = append(imageLinks, link.Attrs()["src"])
		}
	}

	return imageLinks
}
