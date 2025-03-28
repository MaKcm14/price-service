package mmega

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/MaKcm14/price-service/internal/entities/dto"
	"github.com/MaKcm14/price-service/internal/repository/api"
	"github.com/MaKcm14/price-service/pkg/entities"
)

type MegaMarketAPI struct {
	logger       *slog.Logger
	ctx          context.Context
	parser       megaMarketParser
	view         megaMarketViewer
	byPassSocket string
}

func NewMegaMarketAPI(ctx context.Context, log *slog.Logger, socket string) MegaMarketAPI {
	return MegaMarketAPI{
		logger: log,
		ctx:    ctx,
		parser: megaMarketParser{
			logger: log,
		},
		byPassSocket: socket,
	}
}

// getByPassProducts gets the products from by-pass-service.
func (m MegaMarketAPI) getByPassProducts(ctx echo.Context, request dto.ProductRequest, filters []string) (*http.Response, error) {
	const serviceType = "megamarket.service.by-pass-interaction"

	byPassRequest := newByPassServiceRequest(
		request.Query,
		fmt.Sprint(request.Sample),
		m.view.getSortParamURLView(string(request.Sort)),
		m.view.converter.GetFilters(filters),
	)
	requestBody, _ := json.Marshal(byPassRequest)

	if !request.Async && api.IsConnectionClosed(ctx) {
		m.logger.Warn(fmt.Sprintf("error of processing the %v: %v", serviceType, api.ErrConnectionClosed))
		return nil, fmt.Errorf("error of processing the %v: %w", serviceType, api.ErrConnectionClosed)
	}

	resp, err := http.Post(fmt.Sprintf("http://%s/mmarket", m.byPassSocket), "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		m.logger.Warn(fmt.Sprintf("error of the %s: %v", serviceType, err))
		return nil, fmt.Errorf("error of the %s: %w", serviceType, api.ErrByPassServiceResponse)
	} else if resp.StatusCode > 299 {
		errDescr := make(map[string]string)

		json.Marshal(errDescr)

		m.logger.Warn(fmt.Sprintf("error of the %s: %s", serviceType, errDescr["error"]))
		return nil, fmt.Errorf("error of the %s: %w", serviceType, api.ErrByPassServiceResponse)
	}

	return resp, nil
}

// getProducts is the main function of getting the products from the MegaMarket-API calls.
func (m MegaMarketAPI) getProducts(ctx echo.Context, request dto.ProductRequest, filters ...string) (entities.ProductSample, error) {
	const serviceType = "megamarket.service.main-products-getter"

	respByPassProds := struct {
		Items []megaMarketProduct `json:"items"`
	}{
		Items: make([]megaMarketProduct, 0, 100),
	}
	products := make([]entities.Product, 0, 50)

	resp, err := m.getByPassProducts(ctx, request, filters)

	if err != nil {
		m.logger.Warn(fmt.Sprintf("error of the %v: %v", serviceType, err))
		return entities.ProductSample{}, err
	}
	defer resp.Body.Close()

	jsonProducts, err := api.ReadResponseBody(resp.Body, m.logger, serviceType)

	if err != nil {
		return entities.ProductSample{}, fmt.Errorf("error of the %v: error of the reading response body: %v",
			serviceType, err)
	}

	err = json.Unmarshal(jsonProducts, &respByPassProds)

	if err != nil {
		m.logger.Warn(fmt.Sprintf("error of %v: %v: %v", serviceType, api.ErrJSONResponseParsing, err))
		return entities.ProductSample{}, fmt.Errorf("error of the %v: %w: %v", serviceType, api.ErrJSONResponseParsing, err)
	}

	amount := len(respByPassProds.Items)

	if request.Amount == "min" {
		amount = int(math.Min(float64(minValue), float64(amount)))
	}

	for i := 0; i != amount; i++ {
		if respByPassProds.Items[i].FinalPrice == 0 {
			continue
		}
		products = append(products, entities.Product{
			Name:  respByPassProds.Items[i].Goods.Title,
			Brand: respByPassProds.Items[i].Goods.Brand,
			Price: entities.NewPrice(respByPassProds.Items[i].Price, respByPassProds.Items[i].FinalPrice),
			Links: entities.ProductLink{
				URL:       respByPassProds.Items[i].Goods.URL,
				ImageLink: respByPassProds.Items[i].Goods.TitleImageLink,
			},
			Supplier: respByPassProds.Items[i].Offer.MerchantName,
		})
	}

	return entities.NewProductSample(products, m.view.getOpenApiURL(request, filters), entities.MegaMarket), nil
}

// GetProducts gets the products without any filters.
func (m MegaMarketAPI) GetProducts(ctx echo.Context, request dto.ProductRequest) (entities.ProductSample, error) {
	return m.getProducts(ctx, request, sortID, m.view.getSortParamURLView(string(request.Sort)))
}

// GetProductsByPriceRange gets the products with filter by price range.
func (m MegaMarketAPI) GetProductsWithPriceRange(ctx echo.Context, request dto.ProductRequest) (entities.ProductSample, error) {
	return m.getProducts(ctx, request, sortID, m.view.getSortParamURLView(string(request.Sort)),
		priceRangeID, fmt.Sprintf("%d %d", request.PriceRange.PriceDown, request.PriceRange.PriceUp))
}

// GetProductsByExactPrice gets the products with filter by price
// in range [exactPrice, exactPrice + 10% off exactPrice].
func (m MegaMarketAPI) GetProductsWithExactPrice(ctx echo.Context, request dto.ProductRequest) (entities.ProductSample, error) {
	return m.getProducts(ctx, request, sortID, m.view.getSortParamURLView(string(request.Sort)),
		priceRangeID, fmt.Sprintf("%d %d", request.ExactPrice, int(float32(request.ExactPrice)*1.1)))
}

// GetProductsByBestPrice gets the products with filter by min price.
func (m MegaMarketAPI) GetProductsWithBestPrice(ctx echo.Context, request dto.ProductRequest) (entities.ProductSample, error) {
	return m.getProducts(ctx, request, sortID, m.view.getSortParamURLView(string(request.Sort)))
}
