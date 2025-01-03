package mmega

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/MaKcm14/best-price-service/price-service/internal/entities"
	"github.com/MaKcm14/best-price-service/price-service/internal/entities/dto"
	"github.com/MaKcm14/best-price-service/price-service/internal/repository/api"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"github.com/labstack/echo/v4"
)

type MegaMarketAPI struct {
	logger    *slog.Logger
	loadCoeff time.Duration
	ctx       context.Context
	parser    megaMarketParser
	view      megaMarketViewer
}

func NewMegaMarketAPI(ctx context.Context, log *slog.Logger, loadCoeff int) MegaMarketAPI {
	return MegaMarketAPI{
		logger:    log,
		ctx:       ctx,
		loadCoeff: time.Duration(loadCoeff) * time.Millisecond,
		parser: megaMarketParser{
			logger: log,
		},
	}
}

// getHtmlPage gets the html page from the open API url.
func (m MegaMarketAPI) getHtmlPage(url string, request dto.ProductRequest) (string, error) {
	const serviceType = "megamarket.service.html-page-getter"

	var err error
	var html string

	if request.Amount == "min" {
		_, err = chromedp.RunResponse(m.ctx,
			chromedp.Navigate(url),
			chromedp.Sleep(3000*time.Millisecond+m.loadCoeff),
			chromedp.Sleep(3000*time.Millisecond+m.loadCoeff),
			chromedp.Sleep(3000*time.Millisecond+m.loadCoeff),
			chromedp.InnerHTML(fmt.Sprintf("[class='%s']", productContainerClassName), &html),
		)
	} else if request.Amount == "max" {
		_, err = chromedp.RunResponse(m.ctx,
			chromedp.Navigate(url),
			chromedp.Sleep(3000*time.Millisecond+m.loadCoeff),
			chromedp.KeyEvent(kb.End),
			chromedp.Sleep(1000*time.Millisecond+m.loadCoeff),
			chromedp.KeyEvent(kb.End),
			chromedp.Sleep(1000*time.Millisecond+m.loadCoeff),
			chromedp.KeyEvent(kb.End),
			chromedp.Sleep(1000*time.Millisecond+m.loadCoeff),
			chromedp.KeyEvent(kb.End),
			chromedp.Sleep(4000*time.Millisecond+m.loadCoeff),
			chromedp.InnerHTML(fmt.Sprintf("[class='%s']", productContainerClassName), &html),
		)
	}

	if err != nil {
		m.logger.Error(fmt.Sprintf("error of the %s: %v: %v", serviceType, api.ErrChromeDriver, err))
		return "", fmt.Errorf("%w: %v", api.ErrChromeDriver, err)
	}

	return html, nil
}

func (m MegaMarketAPI) getProducts(ctx echo.Context, request dto.ProductRequest, filters ...string) (entities.ProductSample, error) {
	//DEBUG:
	html, _ := m.getHtmlPage(m.view.getOpenApiURL(request, filters), request)

	fmt.Println("parseProducts begun")
	m.parser.parseProducts(html)
	fmt.Println("parseProducts finished")

	return entities.ProductSample{}, api.ErrBufferReading
	///TODO: check and delete
}

// GetProducts gets the products without any filters.
func (m MegaMarketAPI) GetProducts(ctx echo.Context, request dto.ProductRequest) (entities.ProductSample, error) {
	return m.getProducts(ctx, request, sortID, m.view.getSortParamView(string(request.Sort)))
}

// GetProductsByPriceRange gets the products with filter by price range.
func (m MegaMarketAPI) GetProductsWithPriceRange(ctx echo.Context, request dto.ProductRequest, priceDown, priceUp int) (entities.ProductSample, error) {
	return m.getProducts(ctx, request, sortID, m.view.getSortParamView(string(request.Sort)),
		priceRangeID, m.view.getPriceRangeView(priceDown, priceUp))
}

// GetProductsByExactPrice gets the products with filter by price
// in range [exactPrice, exactPrice + 10% off exactPrice].
func (m MegaMarketAPI) GetProductsWithExactPrice(ctx echo.Context, request dto.ProductRequest, exactPrice int) (entities.ProductSample, error) {
	return m.getProducts(ctx, request, sortID, m.view.getSortParamView(string(request.Sort)),
		priceRangeID, m.view.getPriceRangeView(exactPrice, int(float32(exactPrice)*1.1)))
}

// GetProductsByBestPrice gets the products with filter by min price.
func (m MegaMarketAPI) GetProductsWithBestPrice(ctx echo.Context, request dto.ProductRequest) (entities.ProductSample, error) {
	return m.getProducts(ctx, request, sortID, m.view.getSortParamView(string(request.Sort)))
}
