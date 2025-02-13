package filter

import (
	"fmt"

	"github.com/MaKcm14/price-service/internal/entities/dto"
	"github.com/MaKcm14/price-service/pkg/entities"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
)

const (
	testMarket1 entities.Market = iota
	testMarket2
	testMarket3
)

const (
	nameTestMarket1 = "TestMarket1"
	nameTestMarket2 = "TestMarket2"
	nameTestMarket3 = "TestMarket3"
)

type marketApiMock struct {
	mock.Mock
	market              entities.Market
	negativeInteraction bool
}

func newMarketApiMock(market entities.Market, negativeInteraction bool) *marketApiMock {
	return &marketApiMock{
		market:              market,
		negativeInteraction: negativeInteraction,
	}
}

func (m *marketApiMock) getProductsForPositiveCaseInteraction() (entities.ProductSample, error) {
	var marketStr string

	if m.market == testMarket1 {
		marketStr = nameTestMarket1
	} else if m.market == testMarket2 {
		marketStr = nameTestMarket2
	} else if m.market == testMarket3 {
		marketStr = nameTestMarket3
	}

	return entities.ProductSample{
		Products: []entities.Product{{}},
		Market:   marketStr,
	}, nil
}

func (m *marketApiMock) getProductsForNegativeCaseInteraction() (entities.ProductSample, error) {
	return entities.ProductSample{}, fmt.Errorf("test error of the market's api interaction")
}

func (m *marketApiMock) GetProducts(ctx echo.Context, request dto.ProductRequest) (entities.ProductSample, error) {
	m.Called(ctx, request)

	if m.negativeInteraction {
		return m.getProductsForNegativeCaseInteraction()
	}

	return m.getProductsForPositiveCaseInteraction()
}

func (m *marketApiMock) GetProductsWithPriceRange(ctx echo.Context, request dto.ProductRequest) (entities.ProductSample, error) {
	if m.negativeInteraction {
		return m.getProductsForNegativeCaseInteraction()
	}
	return m.getProductsForPositiveCaseInteraction()
}

func (m *marketApiMock) GetProductsWithExactPrice(ctx echo.Context, request dto.ProductRequest) (entities.ProductSample, error) {
	if m.negativeInteraction {
		return m.getProductsForNegativeCaseInteraction()
	}
	return m.getProductsForPositiveCaseInteraction()
}

func (m *marketApiMock) GetProductsWithBestPrice(ctx echo.Context, request dto.ProductRequest) (entities.ProductSample, error) {
	if m.negativeInteraction {
		return m.getProductsForNegativeCaseInteraction()
	}
	return m.getProductsForPositiveCaseInteraction()
}
