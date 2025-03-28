package filter

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/MaKcm14/price-service/internal/entities/dto"
	"github.com/MaKcm14/price-service/internal/services"
	"github.com/MaKcm14/price-service/pkg/entities"
)

type productsFilterTestSuite struct {
	suite.Suite

	mockTestMarket1 *marketApiMock
	mockTestMarket2 *marketApiMock
	mockTestMarket3 *marketApiMock
}

func (s *productsFilterTestSuite) SetupTest() {
	s.mockTestMarket1 = newMarketApiMock(testMarket1, false)
	s.mockTestMarket2 = newMarketApiMock(testMarket2, false)
	s.mockTestMarket3 = newMarketApiMock(testMarket3, false)
}

func (s *productsFilterTestSuite) filterPositiveCaseOfFullMarketsHandlingSettings() {
	s.mockTestMarket1.On("GetProducts", nil, dto.ProductRequest{
		Markets: []entities.Market{
			testMarket1, testMarket2, testMarket3,
		},
	})

	s.mockTestMarket2.On("GetProducts", nil, dto.ProductRequest{
		Markets: []entities.Market{
			testMarket1, testMarket2, testMarket3,
		},
	})

	s.mockTestMarket3.On("GetProducts", nil, dto.ProductRequest{
		Markets: []entities.Market{
			testMarket1, testMarket2, testMarket3,
		},
	})
}

func (s *productsFilterTestSuite) filterExtremeCaseOfNotFullMarketsHandlingSettings() {
	s.mockTestMarket1.On("GetProducts", nil, dto.ProductRequest{
		Markets: []entities.Market{
			testMarket1, testMarket2,
		}})

	s.mockTestMarket2.On("GetProducts", nil, dto.ProductRequest{
		Markets: []entities.Market{
			testMarket1, testMarket2,
		}})
}

func (s *productsFilterTestSuite) filterFullNegativeMarketsApiInteractionSettings() {
	s.mockTestMarket1.negativeInteraction = true
	s.mockTestMarket2.negativeInteraction = true
	s.mockTestMarket3.negativeInteraction = true

	s.mockTestMarket1.On("GetProducts", nil, dto.ProductRequest{
		Markets: []entities.Market{
			testMarket1, testMarket2, testMarket3,
		},
	})

	s.mockTestMarket2.On("GetProducts", nil, dto.ProductRequest{
		Markets: []entities.Market{
			testMarket1, testMarket2, testMarket3,
		},
	})

	s.mockTestMarket3.On("GetProducts", nil, dto.ProductRequest{
		Markets: []entities.Market{
			testMarket1, testMarket2, testMarket3,
		},
	})
}

func (s *productsFilterTestSuite) filterPartialNegativeMarketsApiInteractionSettings() {
	s.mockTestMarket1.negativeInteraction = true

	s.mockTestMarket1.On("GetProducts", nil, dto.ProductRequest{
		Markets: []entities.Market{
			testMarket1, testMarket2, testMarket3,
		},
	})

	s.mockTestMarket2.On("GetProducts", nil, dto.ProductRequest{
		Markets: []entities.Market{
			testMarket1, testMarket2, testMarket3,
		},
	})

	s.mockTestMarket3.On("GetProducts", nil, dto.ProductRequest{
		Markets: []entities.Market{
			testMarket1, testMarket2, testMarket3,
		},
	})
}

func (s *productsFilterTestSuite) BeforeTest(suiteName, testName string) {
	if testName == "TestFilterPositiveCaseOfFullMarketsHandling" {
		s.filterPositiveCaseOfFullMarketsHandlingSettings()

	} else if testName == "TestFilterExtremeCaseOfNotFullMarketsHandling" {
		s.filterExtremeCaseOfNotFullMarketsHandlingSettings()

	} else if testName == "TestFilterFullNegativeMarketsApiInteraction" {
		s.filterFullNegativeMarketsApiInteractionSettings()

	} else if testName == "TestFilterPartialNegativeMarketsApiInteraction" {
		s.filterPartialNegativeMarketsApiInteractionSettings()
	}
}

func (s *productsFilterTestSuite) TestFilterPositiveCaseOfFullMarketsHandling() {
	var testFilterObj = New(
		slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})),
		map[entities.Market]services.ApiInteractor{
			testMarket1: s.mockTestMarket1,
			testMarket2: s.mockTestMarket2,
			testMarket3: s.mockTestMarket3,
		}, nil,
	)

	testProdSample, err := testFilterObj.filter(nil, dto.ProductRequest{
		Markets: []entities.Market{
			testMarket1, testMarket2, testMarket3,
		}}, "test_filter_service", commonFilter)

	if s.NoError(err) {
		count := 0
		for _, testSample := range testProdSample {
			if testSample.Market == nameTestMarket1 {
				count++
			} else if testSample.Market == nameTestMarket2 {
				count++
			} else if testSample.Market == nameTestMarket3 {
				count++
			}
		}
		s.Equal(3, count)
	}
	s.mockTestMarket1.AssertExpectations(s.T())
	s.mockTestMarket1.AssertNumberOfCalls(s.T(), "GetProducts", 1)

	s.mockTestMarket2.AssertExpectations(s.T())
	s.mockTestMarket2.AssertNumberOfCalls(s.T(), "GetProducts", 1)

	s.mockTestMarket3.AssertExpectations(s.T())
	s.mockTestMarket3.AssertNumberOfCalls(s.T(), "GetProducts", 1)
}

func (s *productsFilterTestSuite) TestFilterExtremeCaseOfNotFullMarketsHandling() {
	var testFilterObj = New(
		slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})),
		map[entities.Market]services.ApiInteractor{
			testMarket1: s.mockTestMarket1,
			testMarket2: s.mockTestMarket2,
			testMarket3: s.mockTestMarket3,
		}, nil,
	)

	testProdSample, err := testFilterObj.filter(nil, dto.ProductRequest{
		Markets: []entities.Market{
			testMarket1, testMarket2,
		}}, "test_filter_service", commonFilter)

	if s.NoError(err) {
		for _, testSample := range testProdSample {
			s.NotEqual(nameTestMarket3, testSample.Market)
		}
	}

	s.mockTestMarket1.AssertExpectations(s.T())
	s.mockTestMarket1.AssertNumberOfCalls(s.T(), "GetProducts", 1)

	s.mockTestMarket2.AssertExpectations(s.T())
	s.mockTestMarket2.AssertNumberOfCalls(s.T(), "GetProducts", 1)

	s.mockTestMarket3.AssertNumberOfCalls(s.T(), "GetProducts", 0)
}

func (s *productsFilterTestSuite) TestFilterFullNegativeMarketsApiInteraction() {
	var testFilterObj = New(
		slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})),
		map[entities.Market]services.ApiInteractor{
			testMarket1: s.mockTestMarket1,
			testMarket2: s.mockTestMarket2,
			testMarket3: s.mockTestMarket3,
		}, nil,
	)

	testProdSample, err := testFilterObj.filter(nil, dto.ProductRequest{
		Markets: []entities.Market{
			testMarket1, testMarket2, testMarket3,
		}}, "test_filter_service", commonFilter)

	if s.Error(err) {
		s.Equal(services.ErrGettingProducts, err)
		s.Equal(0, len(testProdSample))
	}

	s.mockTestMarket1.AssertExpectations(s.T())
	s.mockTestMarket2.AssertExpectations(s.T())
	s.mockTestMarket3.AssertExpectations(s.T())
}

func (s *productsFilterTestSuite) TestFilterPartialNegativeMarketsApiInteraction() {
	var testFilterObj = New(
		slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})),
		map[entities.Market]services.ApiInteractor{
			testMarket1: s.mockTestMarket1,
			testMarket2: s.mockTestMarket2,
			testMarket3: s.mockTestMarket3,
		}, nil,
	)

	testProdSample, err := testFilterObj.filter(nil, dto.ProductRequest{
		Markets: []entities.Market{
			testMarket1, testMarket2, testMarket3,
		}}, "test_filter_service", commonFilter)

	if s.NoError(err) {
		count := 0
		for _, testSample := range testProdSample {
			if testSample.Market == nameTestMarket2 {
				count++
			} else if testSample.Market == nameTestMarket3 {
				count++
			}
		}
		s.Equal(2, count)
	}

	s.mockTestMarket1.AssertExpectations(s.T())
	s.mockTestMarket2.AssertExpectations(s.T())
	s.mockTestMarket3.AssertExpectations(s.T())
}

func TestProductsFilter(t *testing.T) {
	suite.Run(t, new(productsFilterTestSuite))
}
