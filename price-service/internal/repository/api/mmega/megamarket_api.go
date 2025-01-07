package mmega

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/MaKcm14/best-price-service/price-service/internal/entities"
	"github.com/MaKcm14/best-price-service/price-service/internal/entities/dto"
	"github.com/MaKcm14/best-price-service/price-service/internal/repository/api"
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

// getProducts is the main function of getting the products from the MegaMarket-API calls.
func (m MegaMarketAPI) getProducts(ctx echo.Context, request dto.ProductRequest, filters ...string) (entities.ProductSample, error) {
	const serviceType = "megamarket.service.main-products-getter"

	products := make([]entities.Product, 0, 50)

	byPassRequest := newByPassServiceRequest(
		request.Query,
		fmt.Sprint(request.Sample),
		m.view.getSortParamURLView(string(request.Sort)),
		m.view.converter.GetFilters(filters),
	)
	requestBody, _ := json.Marshal(byPassRequest)

	if api.IsConnectionClosed(ctx) {
		m.logger.Warn(fmt.Sprintf("error of processing the %v: %v", serviceType, api.ErrConnectionClosed))
		return entities.ProductSample{}, fmt.Errorf("error of processing the %v: %w", serviceType, api.ErrConnectionClosed)
	}

	resp, err := http.Post("http://localhost:8081/mmarket", "application/json", bytes.NewBuffer(requestBody))

	if err != nil || resp.StatusCode > 299 {
		m.logger.Warn(fmt.Sprintf("error of the %s: %v", serviceType, err))
		return entities.ProductSample{}, fmt.Errorf("error of the %s: %w", serviceType, api.ErrByPassServiceResponse)
	}
	defer resp.Body.Close()

	jsonProducts, err := api.ReadResponseBody(resp.Body, m.logger, serviceType)

	if err != nil {
		return entities.ProductSample{}, fmt.Errorf("error of the %v: error of the reading response body: %v",
			serviceType, err)
	}

	_ = jsonProducts

	// parse the products to the []entities.Product

	return entities.NewProductSample(products, m.view.getOpenApiURL(request, filters), entities.MegaMarket), nil
}

// GetProducts gets the products without any filters.
func (m MegaMarketAPI) GetProducts(ctx echo.Context, request dto.ProductRequest) (entities.ProductSample, error) {
	return m.getProducts(ctx, request, sortID, m.view.getSortParamURLView(string(request.Sort)))
}

// GetProductsByPriceRange gets the products with filter by price range.
func (m MegaMarketAPI) GetProductsWithPriceRange(ctx echo.Context, request dto.ProductRequest, priceDown, priceUp int) (entities.ProductSample, error) {
	return m.getProducts(ctx, request, sortID, m.view.getSortParamURLView(string(request.Sort)),
		priceRangeID, fmt.Sprintf("%d %d", priceDown, priceUp))
}

// GetProductsByExactPrice gets the products with filter by price
// in range [exactPrice, exactPrice + 10% off exactPrice].
func (m MegaMarketAPI) GetProductsWithExactPrice(ctx echo.Context, request dto.ProductRequest, exactPrice int) (entities.ProductSample, error) {
	return m.getProducts(ctx, request, sortID, m.view.getSortParamURLView(string(request.Sort)),
		priceRangeID, fmt.Sprintf("%d %d", exactPrice, int(float32(exactPrice)*1.1)))
}

// GetProductsByBestPrice gets the products with filter by min price.
func (m MegaMarketAPI) GetProductsWithBestPrice(ctx echo.Context, request dto.ProductRequest) (entities.ProductSample, error) {
	return m.getProducts(ctx, request, sortID, m.view.getSortParamURLView(string(request.Sort)))
}
