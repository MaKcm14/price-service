package wildb

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"github.com/labstack/echo/v4"

	"github.com/MaKcm14/best-price-service/price-service/internal/entities"
	"github.com/MaKcm14/best-price-service/price-service/internal/entities/dto"
	"github.com/MaKcm14/best-price-service/price-service/internal/repository/api"
)

// WildberriesAPI defines the rules of interaction with the wildberries service and
// provides the interface of getting the products with set clients' filters.
type WildberriesAPI struct {
	logger    *slog.Logger
	loadCoeff time.Duration
	parser    wildberriesParser
	view      wildberriesViewer
	ctx       context.Context
}

func NewWildberriesAPI(ctx context.Context, log *slog.Logger, loadCoeff int) WildberriesAPI {
	return WildberriesAPI{
		logger:    log,
		loadCoeff: time.Duration(loadCoeff) * time.Millisecond,
		parser: wildberriesParser{
			logger: log,
		},
		ctx: ctx,
	}
}

// getHtmlPage gets the raw html (through the open API path) using the filters and the main url's template.
func (w WildberriesAPI) getHtmlPage(url string, request dto.ProductRequest) (string, error) {
	const serviceType = "wildberries.service.html-page-getter"

	var err error
	var html string

	if request.Amount == "min" {
		_, err = chromedp.RunResponse(w.ctx,
			chromedp.Navigate(url),
			chromedp.InnerHTML(fmt.Sprintf("[class='%s']", productContainerClassName), &html),
		)
	} else if request.Amount == "max" {
		_, err = chromedp.RunResponse(w.ctx,
			chromedp.Navigate(url),
			chromedp.Sleep(3000*time.Millisecond+w.loadCoeff),
			chromedp.KeyEvent(kb.End),
			chromedp.Sleep(1000*time.Millisecond+w.loadCoeff),
			chromedp.KeyEvent(kb.End),
			chromedp.Sleep(1000*time.Millisecond+w.loadCoeff),
			chromedp.KeyEvent(kb.End),
			chromedp.Sleep(1000*time.Millisecond+w.loadCoeff),
			chromedp.KeyEvent(kb.End),
			chromedp.Sleep(4000*time.Millisecond+w.loadCoeff),
			chromedp.InnerHTML(fmt.Sprintf("[class='%s']", productContainerClassName), &html),
		)
	}

	if err != nil {
		w.logger.Error(fmt.Sprintf("error of the %s: %v: %v", serviceType, api.ErrChromeDriver, err))
		return "", fmt.Errorf("%w: %v", api.ErrChromeDriver, err)
	}

	return html, nil
}

// getProductsSample gets the json-view structs of the products connected with the current "sample".
func (w WildberriesAPI) getProductSample(url string) ([]wildberriesProduct, error) {
	const serviceType = "wildberries.service.search.wb.ru-products-getter"

	sample := struct {
		Data struct {
			Products []wildberriesProduct `json:"products"`
		} `json:"data"`
	}{
		Data: struct {
			Products []wildberriesProduct `json:"products"`
		}{
			make([]wildberriesProduct, 0, 750),
		},
	}

	for respBody := []byte{}; len(sample.Data.Products) < 10; {
		respBody = respBody[:0]
		sample.Data.Products = sample.Data.Products[:0]

		resp, err := http.Get(url)

		if err != nil || resp.StatusCode > 299 {
			resp.Body.Close()
			w.logger.Warn(fmt.Sprintf("error of the %v: %v: %v", serviceType, api.ErrServiceResponse, err))
			return nil, fmt.Errorf("%w: %v", api.ErrServiceResponse, err)
		}
		defer resp.Body.Close()

		respBody, err := api.ReadResponseBody(resp.Body, w.logger, serviceType)

		if err != nil {
			return nil, fmt.Errorf("error of the %v: error of the reading response body: %v",
				serviceType, err)
		}

		err = json.Unmarshal(respBody, &sample)

		if err != nil {
			w.logger.Error(fmt.Sprintf("error of parsing the json: %v", err))
			return nil, nil
		}
	}

	return sample.Data.Products, nil
}

// getProducts is the main function of getting the products with set filters.
// The current geo-string defines the Moscow info.
func (w WildberriesAPI) getProducts(ctx echo.Context, request dto.ProductRequest, filters ...string) (entities.ProductSample, error) {
	const serviceType = "wildberries.service.main-products-getter"
	var products = make([]entities.Product, 0, 100)

	if api.IsConnectionClosed(ctx) {
		w.logger.Warn(fmt.Sprintf("error of processing the %v: %v", serviceType, api.ErrConnectionClosed))
		return entities.ProductSample{}, fmt.Errorf("error of processing the %v: %w", serviceType, api.ErrConnectionClosed)
	}

	sample, err := w.getProductSample(w.view.getHiddenApiURL(request, filters))

	if err != nil {
		return entities.ProductSample{}, err
	}

	if api.IsConnectionClosed(ctx) {
		w.logger.Warn(fmt.Sprintf("error of processing the %v: %v", serviceType, api.ErrConnectionClosed))
		return entities.ProductSample{}, fmt.Errorf("error of processing the %v: %w", serviceType, api.ErrConnectionClosed)
	}

	htmlSourceLink := w.view.getOpenApiURL(request, filters)

	imageLinks := make([]string, 0, 100)

	if !request.FlagNoImage {
		html, err := w.getHtmlPage(htmlSourceLink, request)

		if err != nil {
			return entities.ProductSample{}, err
		}

		imageLinks = w.parser.parseImageLinks(html)
	}

	for i, j := 0, 0; i != len(sample); i++ {
		var image string

		if j < len(imageLinks) && !request.FlagNoImage {
			image = imageLinks[j]
			j++
		}

		products = append(products, entities.Product{
			Name:     sample[i].Name,
			Brand:    sample[i].Brand,
			Price:    entities.NewPrice(sample[i].Sizes[0].Price.Basic/100, sample[i].Sizes[0].Price.Total/100),
			Supplier: sample[i].Supplier,
			Links: entities.ProductLink{
				URL:       w.view.getProductCatalogLink(sample[i].ID),
				ImageLink: image,
			},
		})
	}

	return entities.NewProductSample(products, htmlSourceLink, entities.Wildberries), nil
}

// GetProducts gets the products without any filters.
func (w WildberriesAPI) GetProducts(ctx echo.Context, request dto.ProductRequest) (entities.ProductSample, error) {
	return w.getProducts(ctx, request, sortID, string(request.Sort))
}

// GetProductsByPriceRange gets the products with filter by price range.
func (w WildberriesAPI) GetProductsWithPriceRange(ctx echo.Context, request dto.ProductRequest) (entities.ProductSample, error) {
	return w.getProducts(ctx, request, sortID, string(request.Sort),
		priceRangeID, w.view.getPriceRangeView(request.PriceRange.PriceDown, request.PriceRange.PriceUp))
}

// GetProductsByExactPrice gets the products with filter by price
// in range [exactPrice, exactPrice + 10% off exactPrice].
func (w WildberriesAPI) GetProductsWithExactPrice(ctx echo.Context, request dto.ProductRequest) (entities.ProductSample, error) {
	return w.getProducts(ctx, request, sortID, string(request.Sort),
		priceRangeID, w.view.getPriceRangeView(request.ExactPrice, int(float32(request.ExactPrice)*1.1)))
}

// GetProductsByBestPrice gets the products with filter by min price.
func (w WildberriesAPI) GetProductsWithBestPrice(ctx echo.Context, request dto.ProductRequest) (entities.ProductSample, error) {
	return w.getProducts(ctx, request, sortID, string(dto.PriceUpSort))
}
