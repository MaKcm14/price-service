package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/MaKcm14/best-price-service/price-service/internal/entities"
	"github.com/anaskhan96/soup"
	"github.com/chromedp/chromedp"
)

type (
	WildberriesPrice struct {
		Basic int `json:"basic"`
		Total int `json:"total"`
	}

	WildberriesProduct struct {
		ID       int    `json:"id"`
		Brand    string `json:"brand"`
		Name     string `json:"name"`
		Supplier string `json:"supplier"`
		Sizes    []struct {
			Price WildberriesPrice `json:"price"`
		} `json:"sizes"`
	}

	WildberriesAPI struct {
		scroll    int
		converter urlConverter
		logger    *slog.Logger
	}
)

func NewWildberriesAPI(log *slog.Logger, scroll int) WildberriesAPI {
	return WildberriesAPI{
		scroll: scroll,
		logger: log,
	}
}

// getHtmlPage gets the raw html using the filters and the main url's template.
func (api WildberriesAPI) getHtmlPage(product entities.ProductRequest, filters ...string) (string, error) {
	const serviceType = "wildberries.service"
	var html string

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	_, err := chromedp.RunResponse(ctx,
		chromedp.Navigate(fmt.Sprintf("https://www.wildberries.ru/catalog/0/search.aspx?%s",
			api.getOpenApiPath(product, filters))),
		chromedp.InnerHTML("[class='product-card-list']", &html),
	)

	///DEBUG:
	//file, err := os.Create("../../internal/repository/api/test.html")
	//file.Write([]byte(html))
	//defer file.Close()
	///TODO: delete

	if err != nil {
		api.logger.Error(fmt.Sprintf("error of the %s: %v: %v", serviceType, ErrChromeDriver, err))
		return "", fmt.Errorf("%v: %v", ErrChromeDriver, err)
	}

	return html, nil
}

// getImageLinks gets image links for the products.
func (api WildberriesAPI) getImageLinks(html string) []string {
	var imageLinks = make([]string, 0, 750)
	var parse = soup.HTMLParse(html)

	for _, tag := range parse.FindAll("article") {
		link := tag.Find("img", "class", "j-thumbnail")
		imageLinks = append(imageLinks, link.Attrs()["src"])
	}

	return imageLinks
}

// getOpenApiPath returns the correct URL's path for wildberries open API.
// It uses with domain "www.wildberries.ru"
func (api WildberriesAPI) getOpenApiPath(product entities.ProductRequest, filters []string) string {
	var path string
	filtersURL := api.converter.getFilters(filters)

	path += fmt.Sprintf("page=%d", product.Sample)
	path += "&sort=" + filtersURL["sort"]

	if priceRange, flagExist := filtersURL["priceU"]; flagExist {
		path += "&priceU=" + priceRange
	}

	path += "&search=" + strings.Join(strings.Split(product.ProductName, " "), "+")

	return path
}

// getHiddenApiPath returns the correct URL's path for wildberries hidden API.
// It uses with domain "search.wb.ru"
func (api WildberriesAPI) getHiddenApiPath(product entities.ProductRequest, filters []string) string {
	var path string
	filtersURL := api.converter.getFilters(filters)

	path += fmt.Sprintf("page=%d", product.Sample)

	if priceRange, flagExist := filtersURL["priceU"]; flagExist {
		path += "&priceU=" + priceRange
	}

	path += "&query=" + strings.Join(strings.Split(product.ProductName, " "), "%20")
	path += "&resultset=catalog&sort=" + filtersURL["sort"]
	path += "spp=30suppressSpellcheck=false"

	return path
}

// getProducts is the main function of getting the products with set filters.
// The current geo-string defines the Moscow info.
func (api WildberriesAPI) getProducts(product entities.ProductRequest, filters ...string) ([]entities.Product, error) {
	const serviceType = "wildberries.service"
	var products = make([]entities.Product, 0, 750)

	resp, err := http.Get(fmt.Sprintf("https://search.wb.ru/exactmatch/ru/common/v9/search?"+
		"ab_testing=false&appType=1&curr=rub&dest=-1257786&hide_dtype=10&lang=ru&%s",
		api.getHiddenApiPath(product, filters)),
	)

	if err != nil || resp.StatusCode > 299 {
		api.logger.Warn(fmt.Sprintf("error of the %v: %v: %v", serviceType, ErrServiceResponse, err))
		return nil, fmt.Errorf("%v: %v", ErrServiceResponse, err)
	}
	defer resp.Body.Close()

	respProd := struct {
		Data struct {
			Products []WildberriesProduct `json:"products"`
		} `json:"data"`
	}{
		Data: struct {
			Products []WildberriesProduct `json:"products"`
		}{
			make([]WildberriesProduct, 0, 750),
		},
	}

	respBody := make([]byte, 0, 200000)

	for {
		buffer := make([]byte, 100000)
		n, err := resp.Body.Read(buffer)

		if n != 0 && (err == nil || err == io.EOF) {
			respBody = append(respBody, buffer[:n]...)
			continue
		} else if err != nil && err != io.EOF {
			api.logger.Warn(fmt.Sprintf("error of the %v: %v: %v", serviceType, ErrBufferReading, err))
			return nil, fmt.Errorf("%v: %v", ErrBufferReading, err)
		} else if err == io.EOF {
			break
		}
	}

	err = json.Unmarshal(respBody, &respProd)

	if err != nil {
		api.logger.Error(fmt.Sprintf("error of parsing the json: %v", err))
		return nil, nil
	}

	html, err := api.getHtmlPage(product, filters...)
	imageLinks := api.getImageLinks(html)

	for i := 0; i != len(respProd.Data.Products); i++ {
		var imageLink = ""
		var product = respProd.Data.Products[i]

		if i < len(imageLinks) {
			imageLink = imageLinks[i]
		}

		discount := (product.Sizes[0].Price.Basic/100 - product.Sizes[0].Price.Total/100) * 100
		discount /= product.Sizes[0].Price.Basic / 100

		products = append(products, entities.Product{
			Name:          product.Name,
			Brand:         product.Brand,
			BasePrice:     product.Sizes[0].Price.Basic / 100,
			DiscountPrice: product.Sizes[0].Price.Total / 100,
			Discount:      discount,
			URL:           fmt.Sprintf("https://www.wildberries.ru/catalog/%v/detail.aspx", product.ID),
			Market:        Wildberries,
			Supplier:      product.Supplier,
			ImageLink:     imageLink,
		})
	}

	return products, nil
}

// GetProducts gets the products without any filters.
func (api WildberriesAPI) GetProducts(product entities.ProductRequest) ([]entities.Product, error) {
	return api.getProducts(product, "sort", "popular")
}
