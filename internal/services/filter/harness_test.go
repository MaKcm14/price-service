package filter

import (
	"fmt"

	"github.com/MaKcm14/best-price-service/price-service/internal/entities"
	"github.com/MaKcm14/best-price-service/price-service/internal/entities/dto"
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

func (s *marketApiMock) getProductsForPositiveCaseInteraction() (entities.ProductSample, error) {
	var marketStr string

	if s.market == testMarket1 {
		marketStr = nameTestMarket1
	} else if s.market == testMarket2 {
		marketStr = nameTestMarket2
	} else if s.market == testMarket3 {
		marketStr = nameTestMarket3
	}

	return entities.ProductSample{
		Products: []entities.Product{{}},
		Market:   marketStr,
	}, nil
}

func (s *marketApiMock) getProductsForNegativeCaseInteraction() (entities.ProductSample, error) {
	return entities.ProductSample{}, fmt.Errorf("test error of the market's api interaction")
}

func (s *marketApiMock) GetProducts(ctx echo.Context, request dto.ProductRequest) (entities.ProductSample, error) {
	s.Called(ctx, request)

	if s.negativeInteraction {
		return s.getProductsForNegativeCaseInteraction()
	}

	return s.getProductsForPositiveCaseInteraction()
}

func (s *marketApiMock) GetProductsWithPriceRange(ctx echo.Context, request dto.ProductRequest) (entities.ProductSample, error) {
	if s.negativeInteraction {
		return s.getProductsForNegativeCaseInteraction()
	}
	return s.getProductsForPositiveCaseInteraction()
}

func (s *marketApiMock) GetProductsWithExactPrice(ctx echo.Context, request dto.ProductRequest) (entities.ProductSample, error) {
	if s.negativeInteraction {
		return s.getProductsForNegativeCaseInteraction()
	}
	return s.getProductsForPositiveCaseInteraction()
}

func (s *marketApiMock) GetProductsWithBestPrice(ctx echo.Context, request dto.ProductRequest) (entities.ProductSample, error) {
	if s.negativeInteraction {
		return s.getProductsForNegativeCaseInteraction()
	}
	return s.getProductsForPositiveCaseInteraction()
}
