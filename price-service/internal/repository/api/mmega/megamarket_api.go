package mmega

import (
	"context"
	"fmt"
	"log/slog"
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

func (m MegaMarketAPI) getProducts(ctx echo.Context, request dto.ProductRequest, filters ...string) (entities.ProductSample, error) {
	//DEBUG:
	fmt.Println(m.view.getOpenApiURL(request, filters))

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
