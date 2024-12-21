package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/MaKcm14/best-price-service/price-service/internal/entities"
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
	var html string

	/*
		err := chromedp.Run(context.Background(),
			chromedp.Navigate(fmt.Sprintf("https://www.wildberries.ru/catalog/0/search.aspx?&search=%s%s",
				api.converter.convertToURLView(product.ProductName), api.converter.convertFilters(filters))),
			chromedp.Evaluate(fmt.Sprintf("window.scrollBy(0, %v)", api.scroll), nil),
			chromedp.OuterHTML("html", &html),
		)

		if err != nil {
			api.logger.Error(fmt.Sprintf("%v: %v", ErrChromeDriver, err))
			return "", fmt.Errorf("%v: %v", ErrChromeDriver, err)
		}
	*/

	return html, nil
}

// parseHtmlForImage gets the image's links for products.
func (api WildberriesAPI) parseHtmlForImage(html string) ([]string, error) {
	///DEBUG:
	fmt.Println(html)
	///TODO: delete
	return nil, nil
}

// composeURLPath returns the correct URL's path for wildberries API.
func (api WildberriesAPI) composeURLPath(product entities.ProductRequest, filters []string) string {
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
		api.composeURLPath(product, filters)),
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

	// TODO: get url for image data
	for _, product := range respProd.Data.Products {
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
		})
	}

	return products, nil
}

// GetProducts gets the products without any filters.
func (api WildberriesAPI) GetProducts(product entities.ProductRequest) {
	api.getProducts(product, "sort", "popular")
}
