package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"

	"github.com/labstack/echo/v4"

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
		converter urlConverter
		loadCoeff time.Duration
		logger    *slog.Logger
	}
)

func NewWildberriesAPI(log *slog.Logger, loadCoeff int) WildberriesAPI {
	return WildberriesAPI{
		logger:    log,
		loadCoeff: time.Duration(loadCoeff),
	}
}

// getHtmlPage gets the raw html using the filters and the main url's template.
func (api WildberriesAPI) getHtmlPage(url string, request entities.ProductRequest) (string, error) {
	const serviceType = "wildberries.service"
	var html string
	var err error

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
	)
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	if request.Amount == "min" {
		_, err = chromedp.RunResponse(ctx,
			chromedp.Navigate(url),
			chromedp.InnerHTML("[class='product-card-list']", &html),
		)
	} else if request.Amount == "max" {
		_, err = chromedp.RunResponse(ctx,
			chromedp.Navigate(url),
			chromedp.Sleep(3000*api.loadCoeff*time.Millisecond),
			chromedp.KeyEvent(kb.End),
			chromedp.Sleep(1000*api.loadCoeff*time.Millisecond),
			chromedp.KeyEvent(kb.End),
			chromedp.Sleep(1000*api.loadCoeff*time.Millisecond),
			chromedp.KeyEvent(kb.End),
			chromedp.Sleep(1000*api.loadCoeff*time.Millisecond),
			chromedp.KeyEvent(kb.End),
			chromedp.Sleep(4000*api.loadCoeff*time.Millisecond),
			chromedp.InnerHTML("[class='product-card-list']", &html),
		)
	}

	if err != nil {
		api.logger.Error(fmt.Sprintf("error of the %s: %v: %v", serviceType, ErrChromeDriver, err))
		return "", fmt.Errorf("%v: %v", ErrChromeDriver, err)
	}

	return html, nil
}

// getImageLinks gets image links for the products.
func (api WildberriesAPI) getImageLinks(html string) []string {
	const serviceType = "wildberries.service"

	if !strings.Contains(html, "article") {
		api.logger.Warn(fmt.Sprintf("error of the %v: %v: images couldn't be load", serviceType, ErrServiceResponse))
		return nil
	}

	var imageLinks = make([]string, 0, 750)
	var parse = soup.HTMLParse(html)

	for _, tag := range parse.FindAll("article") {
		link := tag.Find("img", "class", "j-thumbnail")
		imageLinks = append(imageLinks, link.Attrs()["src"])
	}

	return imageLinks
}

// getOpenApiPath returns the correct URL's path for wildberries open API.
// It uses with domain "www.wildberries.ru".
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
// It uses with domain "search.wb.ru".
func (api WildberriesAPI) getHiddenApiPath(product entities.ProductRequest, filters []string) string {
	var path string
	filtersURL := api.converter.getFilters(filters)

	path += fmt.Sprintf("page=%d", product.Sample)

	if priceRange, flagExist := filtersURL["priceU"]; flagExist {
		path += "&priceU=" + priceRange
	}

	path += "&query=" + url.QueryEscape(product.ProductName)

	path += "&resultset=catalog&sort=" + filtersURL["sort"]
	path += "&spp=30&suppressSpellcheck=false"

	return path
}

func (api WildberriesAPI) getResponse(url string) ([]WildberriesProduct, error) {
	const serviceType = "wildberries.service"
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

	for len(respProd.Data.Products) < 10 {
		respBody = respBody[:0]
		respProd.Data.Products = respProd.Data.Products[:0]

		resp, err := http.Get(url)

		if err != nil || resp.StatusCode > 299 {
			api.logger.Warn(fmt.Sprintf("error of the %v: %v: %v", serviceType, ErrServiceResponse, err))
			return nil, fmt.Errorf("%v: %v", ErrServiceResponse, err)
		}
		defer resp.Body.Close()

		for {
			buffer := make([]byte, 100000)
			n, err := resp.Body.Read(buffer)

			if n != 0 && (err == nil || err == io.EOF) {
				respBody = append(respBody, buffer[:n]...)
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
	}
	return respProd.Data.Products, nil
}

// getProducts is the main function of getting the products with set filters.
// The current geo-string defines the Moscow info.
func (api WildberriesAPI) getProducts(ctx echo.Context, product entities.ProductRequest, filters ...string) (entities.ProductSample, error) {
	const serviceType = "wildberries.service"
	var products = make([]entities.Product, 0, 750)

	if IsConnectionClosed(ctx) {
		api.logger.Warn(fmt.Sprintf("error of processing the %v: %v", serviceType, ErrConnectionClosed))
		return entities.ProductSample{}, fmt.Errorf("error of processing the %v: %v", serviceType, ErrConnectionClosed)
	}

	respProd, err := api.getResponse(fmt.Sprintf("https://search.wb.ru/exactmatch/ru/common/v9/search?"+
		"ab_testing=false&appType=1&curr=rub&dest=-1257786&hide_dtype=10&lang=ru&%s",
		api.getHiddenApiPath(product, filters)))

	if err != nil {
		return entities.ProductSample{}, err
	}

	if IsConnectionClosed(ctx) {
		api.logger.Warn(fmt.Sprintf("error of processing the %v: %v", serviceType, ErrConnectionClosed))
		return entities.ProductSample{}, fmt.Errorf("error of processing the %v: %v", serviceType, ErrConnectionClosed)
	}

	prodsLink := fmt.Sprintf("https://www.wildberries.ru/catalog/0/search.aspx?%s",
		api.getOpenApiPath(product, filters))

	html, err := api.getHtmlPage(prodsLink, product)

	if err != nil {
		return entities.ProductSample{}, err
	}

	imageLinks := api.getImageLinks(html)

	for i, j := 0, 0; i != int(math.Min(float64(len(respProd)), float64(len(imageLinks)))); i++ {
		var imageLink string

		if j < len(imageLinks) {
			imageLink = imageLinks[j]
			j++
		}

		products = append(products, entities.Product{
			Name:     respProd[i].Name,
			Brand:    respProd[i].Brand,
			Price:    entities.NewPrice(respProd[i].Sizes[0].Price.Basic/100, respProd[i].Sizes[0].Price.Total/100),
			Market:   Wildberries,
			Supplier: respProd[i].Supplier,
			MetaData: entities.ProductMetaData{
				URL:       fmt.Sprintf("https://www.wildberries.ru/catalog/%v/detail.aspx", respProd[i].ID),
				ImageLink: imageLink,
			},
		})
	}

	return entities.NewProductSample(products, prodsLink, entities.Wildberries), nil
}

// GetProducts gets the products without any filters.
func (api WildberriesAPI) GetProducts(ctx echo.Context, product entities.ProductRequest) (entities.ProductSample, error) {
	return api.getProducts(ctx, product, "sort", "popular")
}

// GetProductsByPriceRange gets the products with filter by price range.
func (api WildberriesAPI) GetProductsByPriceRange(ctx echo.Context, product entities.ProductRequest, priceDown, priceUp int) (entities.ProductSample, error) {
	return api.getProducts(ctx, product, "sort", "popular",
		"priceU", fmt.Sprintf("%v00;%v00", priceDown, priceUp))
}

// GetProductsByExactPrice gets the products with filter by price
// in range [exactPrice, exactPrice + 10% off exactPrice].
func (api WildberriesAPI) GetProductsByExactPrice(ctx echo.Context, product entities.ProductRequest, exactPrice int) (entities.ProductSample, error) {
	return api.getProducts(ctx, product, "sort", "priceup",
		"priceU", fmt.Sprintf("%v00;%v00", exactPrice, int(float32(exactPrice)*1.1)))
}

// GetProductsByBestPrice gets the products with filter by min price.
func (api WildberriesAPI) GetProductsByBestPrice(ctx echo.Context, product entities.ProductRequest) (entities.ProductSample, error) {
	return api.getProducts(ctx, product, "sort", "priceup")
}
