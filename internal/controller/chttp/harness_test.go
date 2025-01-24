package chttp

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"

	"github.com/MaKcm14/best-price-service/price-service/internal/entities"
	"github.com/MaKcm14/best-price-service/price-service/internal/entities/dto"
)

const (
	nameTestMarket1 = "TestMarket1"
)

type productsFilterMock struct {
	mock.Mock

	negativeInteraction bool
}

func newProductsFilterMock(negativeInteraction bool) *productsFilterMock {
	return &productsFilterMock{
		negativeInteraction: negativeInteraction,
	}
}

func (m *productsFilterMock) getPositiveCaseSample() []entities.ProductSample {
	return []entities.ProductSample{
		{
			Products: []entities.Product{{}},
			Market:   nameTestMarket1,
		},
	}
}

func (m *productsFilterMock) FilterByMarkets(ctx echo.Context, request dto.ProductRequest) ([]entities.ProductSample, error) {
	if m.negativeInteraction {
		return nil, fmt.Errorf("error of the filter's interaction")
	}
	return m.getPositiveCaseSample(), nil
}

func (m *productsFilterMock) FilterByPriceRange(ctx echo.Context, request dto.ProductRequest) ([]entities.ProductSample, error) {
	if m.negativeInteraction {
		return nil, fmt.Errorf("error of the filter's interaction")
	}
	return m.getPositiveCaseSample(), nil
}

func (m *productsFilterMock) FilterByBestPrice(ctx echo.Context, request dto.ProductRequest) ([]entities.ProductSample, error) {
	if m.negativeInteraction {
		return nil, fmt.Errorf("error of the filter's interaction")
	}
	return m.getPositiveCaseSample(), nil
}

func (m *productsFilterMock) FilterByExactPrice(ctx echo.Context, request dto.ProductRequest) ([]entities.ProductSample, error) {
	if m.negativeInteraction {
		return nil, fmt.Errorf("error of the filter's interaction")
	}
	return m.getPositiveCaseSample(), nil
}
