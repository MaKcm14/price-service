package chttp

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/MaKcm14/best-price-service/price-service/internal/entities/dto"
	"github.com/MaKcm14/best-price-service/price-service/pkg/entities"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type handlersTestSuite struct {
	suite.Suite
	ctx        echo.Context
	filterMock *productsFilterMock
}

func (s *handlersTestSuite) testInputCase(testName string, path string, testF func(echo.Context) error) {
	s.T().Run(testName, func(t *testing.T) {
		s.ctx = echo.New().NewContext(
			httptest.NewRequest("GET", path, nil),
			httptest.NewRecorder(),
		)

		err := testF(s.ctx)

		assert.NoError(s.T(), err)

		assert.Equal(s.T(), http.StatusBadRequest, s.ctx.Response().Status)
	})
}

func (s *handlersTestSuite) SetupTest() {
	s.filterMock = newProductsFilterMock(false, false)
	s.ctx = echo.New().NewContext(
		httptest.NewRequest("GET", "/test/path?query=test+query&markets=wildberries&sort=popular&sample=1&no-image=1&amount=min&price_down=1000&price_up=5000&price=5000", nil),
		httptest.NewRecorder(),
	)
}

func (s *handlersTestSuite) BeforeTest(suiteName, testName string) {
	if testName == "TestHandlePriceRangeRequestPositiveCase" ||
		testName == "TestHandlePriceRangeRequestNegativeCasesFilterInteraction" {
		s.filterMock.On("FilterByPriceRange", mock.Anything, dto.ProductRequest{
			Query:       "test query",
			Sample:      1,
			Amount:      "min",
			Sort:        "popular",
			FlagNoImage: true,
			Markets:     []entities.Market{entities.Wildberries},
			PriceRange: dto.PriceRangeRequest{
				PriceDown: 1000,
				PriceUp:   5000,
			},
		})
	} else if testName == "TestHandleExactPriceRequestPositiveCase" ||
		testName == "TestHandleExactPriceRequestNegativeCasesFilterInteraction" {
		s.filterMock.On("FilterByExactPrice", mock.Anything, dto.ProductRequest{
			Query:       "test query",
			Sample:      1,
			Amount:      "min",
			Sort:        "popular",
			FlagNoImage: true,
			Markets:     []entities.Market{entities.Wildberries},
			ExactPrice:  5000,
		})
	}
}

func (s *handlersTestSuite) TestHandleMarketsRequestPositiveCase() {
	var testContrObj = Controller{
		logger: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})),
		filter: s.filterMock,
	}

	err := testContrObj.handleMarketsRequest(s.ctx)

	assert.NoError(s.T(), err)

	assert.Equal(s.T(), http.StatusOK, s.ctx.Response().Status)
}

func (s *handlersTestSuite) TestHandleMarketsRequestNegativeCasesFilterInteraction() {
	s.T().Run("Negative Case: error of the service", func(t *testing.T) {
		var testContrObj = Controller{
			logger: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})),
			filter: s.filterMock,
		}
		s.filterMock.negativeInteraction = false
		s.filterMock.serviceError = true

		err := testContrObj.handleMarketsRequest(s.ctx)

		assert.NoError(s.T(), err)

		assert.Equal(s.T(), http.StatusInternalServerError, s.ctx.Response().Status)
	})

	s.T().Run("Negative Case: error of the gateway services", func(t *testing.T) {
		var testContrObj = Controller{
			logger: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})),
			filter: s.filterMock,
		}
		s.filterMock.serviceError = false
		s.filterMock.negativeInteraction = true

		err := testContrObj.handleMarketsRequest(s.ctx)

		assert.NoError(s.T(), err)

		assert.Equal(s.T(), http.StatusBadGateway, s.ctx.Response().Status)
	})
}

func (s *handlersTestSuite) TestHandleMarketsRequestNegativeInputCases() {
	type args struct {
		path string
	}

	tests := []struct {
		name string
		args args
	}{
		{"Negative case: the wrong markets parameter", args{"/test/path?query=test+query&markets=asvnjas&sort=popular&sample=1&no-image=1&amount=min"}},
		{"Negative case: the empty markets parameter", args{"/test/path?query=test+query&sort=popular&sample=1&no-image=1&amount=min"}},
		{"Negative case: the wrong query parameter", args{"/test/path?&markets=wildberries&sort=popular&sample=1&no-image=1&amount=min"}},
		{"Negative case: the empty query parameter", args{"/test/path?&markets=wildberries&sort=popular&sample=1&no-image=1&amount=min"}},
	}

	testContrObj := Controller{
		logger: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})),
		filter: s.filterMock,
	}

	for _, test := range tests {
		s.testInputCase(test.name, test.args.path, testContrObj.handleMarketsRequest)
	}
}

func (s *handlersTestSuite) TestHandleBestPriceRequestPositiveCase() {
	var testContrObj = Controller{
		logger: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})),
		filter: s.filterMock,
	}

	err := testContrObj.handleBestPriceRequest(s.ctx)

	assert.NoError(s.T(), err)

	assert.Equal(s.T(), http.StatusOK, s.ctx.Response().Status)
}

func (s *handlersTestSuite) TestHandleBestPriceRequestNegativeCasesFilterInteraction() {
	s.T().Run("Negative Case: error of the service", func(t *testing.T) {
		var testContrObj = Controller{
			logger: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})),
			filter: s.filterMock,
		}
		s.filterMock.negativeInteraction = false
		s.filterMock.serviceError = true

		err := testContrObj.handleBestPriceRequest(s.ctx)

		assert.NoError(s.T(), err)

		assert.Equal(s.T(), http.StatusInternalServerError, s.ctx.Response().Status)
	})

	s.T().Run("Negative Case: error of the gateway services", func(t *testing.T) {
		var testContrObj = Controller{
			logger: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})),
			filter: s.filterMock,
		}
		s.filterMock.serviceError = false
		s.filterMock.negativeInteraction = true

		err := testContrObj.handleBestPriceRequest(s.ctx)

		assert.NoError(s.T(), err)

		assert.Equal(s.T(), http.StatusBadGateway, s.ctx.Response().Status)
	})
}

func (s *handlersTestSuite) TestHandleBestPriceRequestNegativeInputCases() {
	type args struct {
		path string
	}

	tests := []struct {
		name string
		args args
	}{
		{"Negative case: the wrong markets parameter", args{"/test/path?query=test+query&markets=asvnjas&sort=popular&sample=1&no-image=1&amount=min"}},
		{"Negative case: the empty markets parameter", args{"/test/path?query=test+query&sort=popular&sample=1&no-image=1&amount=min"}},
		{"Negative case: the wrong query parameter", args{"/test/path?&markets=wildberries&sort=popular&sample=1&no-image=1&amount=min"}},
		{"Negative case: the empty query parameter", args{"/test/path?&markets=wildberries&sort=popular&sample=1&no-image=1&amount=min"}},
	}

	testContrObj := Controller{
		logger: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})),
		filter: s.filterMock,
	}

	for _, test := range tests {
		s.testInputCase(test.name, test.args.path, testContrObj.handleBestPriceRequest)
	}
}

func (s *handlersTestSuite) TestHandlePriceRangeRequestPositiveCase() {
	var testContrObj = Controller{
		logger: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})),
		filter: s.filterMock,
	}

	err := testContrObj.handlePriceRangeRequest(s.ctx)

	assert.NoError(s.T(), err)

	assert.Equal(s.T(), http.StatusOK, s.ctx.Response().Status)

	s.filterMock.AssertExpectations(s.T())
}

func (s *handlersTestSuite) TestHandlePriceRangeRequestNegativeCasesFilterInteraction() {
	s.T().Run("Negative Case: error of the service", func(t *testing.T) {
		var testContrObj = Controller{
			logger: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})),
			filter: s.filterMock,
		}
		s.filterMock.negativeInteraction = false
		s.filterMock.serviceError = true

		err := testContrObj.handlePriceRangeRequest(s.ctx)

		assert.NoError(s.T(), err)

		assert.Equal(s.T(), http.StatusInternalServerError, s.ctx.Response().Status)

		s.filterMock.AssertExpectations(s.T())
	})

	s.T().Run("Negative Case: error of the gateway services", func(t *testing.T) {
		var testContrObj = Controller{
			logger: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})),
			filter: s.filterMock,
		}
		s.filterMock.serviceError = false
		s.filterMock.negativeInteraction = true

		err := testContrObj.handlePriceRangeRequest(s.ctx)

		assert.NoError(s.T(), err)

		assert.Equal(s.T(), http.StatusBadGateway, s.ctx.Response().Status)

		s.filterMock.AssertExpectations(s.T())
	})
}

func (s *handlersTestSuite) TestHandlePriceRangeRequestNegativeInputCases() {
	type args struct {
		path string
	}

	tests := []struct {
		name string
		args args
	}{
		{"Negative case: the wrong markets parameter", args{"/test/path?query=test+query&markets=asvnjas&sort=popular&sample=1&no-image=1&amount=min"}},
		{"Negative case: the empty markets parameter", args{"/test/path?query=test+query&sort=popular&sample=1&no-image=1&amount=min"}},
		{"Negative case: the wrong query parameter", args{"/test/path?&markets=wildberries&sort=popular&sample=1&no-image=1&amount=min"}},
		{"Negative case: the empty query parameter", args{"/test/path?&markets=wildberries&sort=popular&sample=1&no-image=1&amount=min"}},
		{"Negative case: the wrong price_down parameter: wrong type of data got", args{"/test/path?&markets=wildberries&sort=popular&sample=1&no-image=1&amount=min&price_down=hello&price_up=1000"}},
		{"Negative case: the wrong price_up parameter: wrong type of data got", args{"/test/path?&markets=wildberries&sort=popular&sample=1&no-image=1&amount=min&price_down=1000&price_up=hello"}},
		{"Negative case: the wrong price_down parameter: wrong value", args{"/test/path?&markets=wildberries&sort=popular&sample=1&no-image=1&amount=min&price_down=-1&price_up=1000"}},
		{"Negative case: the wrong price_up parameter: wrong value", args{"/test/path?&markets=wildberries&sort=popular&sample=1&no-image=1&amount=min&price_down=1000&price_up=-1"}},
		{"Negative case: the wrong price_down parameter: empty value", args{"/test/path?&markets=wildberries&sort=popular&sample=1&no-image=1&amount=min&price_up=1000"}},
		{"Negative case: the wrong price_up parameter: empty value", args{"/test/path?&markets=wildberries&sort=popular&sample=1&no-image=1&amount=min&price_down=1000"}},
	}

	testContrObj := Controller{
		logger: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})),
		filter: s.filterMock,
	}

	for _, test := range tests {
		s.testInputCase(test.name, test.args.path, testContrObj.handlePriceRangeRequest)
	}
}

func (s *handlersTestSuite) TestHandleExactPriceRequestPositiveCase() {
	var testContrObj = Controller{
		logger: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})),
		filter: s.filterMock,
	}

	err := testContrObj.handleExactPriceRequest(s.ctx)

	assert.NoError(s.T(), err)

	assert.Equal(s.T(), http.StatusOK, s.ctx.Response().Status)

	s.filterMock.AssertExpectations(s.T())
}

func (s *handlersTestSuite) TestHandleExactPriceRequestNegativeCasesFilterInteraction() {
	s.T().Run("Negative Case: error of the service", func(t *testing.T) {
		var testContrObj = Controller{
			logger: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})),
			filter: s.filterMock,
		}
		s.filterMock.negativeInteraction = false
		s.filterMock.serviceError = true

		err := testContrObj.handleExactPriceRequest(s.ctx)

		assert.NoError(s.T(), err)

		assert.Equal(s.T(), http.StatusInternalServerError, s.ctx.Response().Status)

		s.filterMock.AssertExpectations(s.T())
	})

	s.T().Run("Negative Case: error of the gateway services", func(t *testing.T) {
		var testContrObj = Controller{
			logger: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})),
			filter: s.filterMock,
		}
		s.filterMock.serviceError = false
		s.filterMock.negativeInteraction = true

		err := testContrObj.handleExactPriceRequest(s.ctx)

		assert.NoError(s.T(), err)

		assert.Equal(s.T(), http.StatusBadGateway, s.ctx.Response().Status)

		s.filterMock.AssertExpectations(s.T())
	})
}

func (s *handlersTestSuite) TestHandleExactPriceRequestNegativeInputCases() {
	type args struct {
		path string
	}

	tests := []struct {
		name string
		args args
	}{
		{"Negative case: the wrong markets parameter", args{"/test/path?query=test+query&markets=asvnjas&sort=popular&sample=1&no-image=1&amount=min"}},
		{"Negative case: the empty markets parameter", args{"/test/path?query=test+query&sort=popular&sample=1&no-image=1&amount=min"}},
		{"Negative case: the wrong query parameter", args{"/test/path?&markets=wildberries&sort=popular&sample=1&no-image=1&amount=min"}},
		{"Negative case: the empty query parameter", args{"/test/path?&markets=wildberries&sort=popular&sample=1&no-image=1&amount=min"}},
		{"Negative case: the wrong price parameter: wrong type of data", args{"/test/path?&markets=wildberries&sort=popular&sample=1&no-image=1&amount=min&price_down=1000&price_up=5000&price=as"}},
		{"Negative case: the wrong price parameter: wrong value", args{"/test/path?&markets=wildberries&sort=popular&sample=1&no-image=1&amount=min&price_down=1000&price_up=5000&price=-1"}},
		{"Negative case: the wron price parameter: empty value", args{"/test/path?&markets=wildberries&sort=popular&sample=1&no-image=1&amount=min&price_down=1000&price_up=5000"}},
	}

	testContrObj := Controller{
		logger: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})),
		filter: s.filterMock,
	}

	for _, test := range tests {
		s.testInputCase(test.name, test.args.path, testContrObj.handleExactPriceRequest)
	}
}

func TestContollerHandlers(t *testing.T) {
	suite.Run(t, new(handlersTestSuite))
}
