package chttp

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type handlersTestSuite struct {
	suite.Suite
	ctx        echo.Context
	filterMock *productsFilterMock
}

func (s *handlersTestSuite) SetupTest() {
	s.filterMock = newProductsFilterMock(false)
}

func (s *handlersTestSuite) BeforeTest(suiteName, testName string) {
	if testName == "TestHandleMarketsRequestNegativeCaseFilterInteraction" ||
		testName == "TestHandleBestPriceRequestNegativeCaseFilterInteraction" {
		s.filterMock.negativeInteraction = true
		s.ctx = echo.New().NewContext(
			httptest.NewRequest("GET", "/products/filter/markets?query=test+query&markets=wildberries&sort=popular&sample=1&no-image=1&amount=min", nil),
			httptest.NewRecorder(),
		)

	} else if testName == "TestHandleMarketsRequestNegativeInputCase" ||
		testName == "TestHandleBestPriceRequestNegativeInputCase" {
		s.ctx = echo.New().NewContext(
			httptest.NewRequest("GET", "/products/filter/markets?query=test+query&markets=asvnjas&sort=popular&sample=1&no-image=1&amount=min", nil),
			httptest.NewRecorder(),
		)
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

func (s *handlersTestSuite) TestHandleMarketsRequestNegativeCaseFilterInteraction() {
	var testContrObj = Controller{
		logger: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})),
		filter: s.filterMock,
	}

	err := testContrObj.handleMarketsRequest(s.ctx)

	assert.NoError(s.T(), err)

	assert.Equal(s.T(), http.StatusInternalServerError, s.ctx.Response().Status)
}

func (s *handlersTestSuite) TestHandleMarketsRequestNegativeInputCase() {
	var testContrObj = Controller{
		logger: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})),
		filter: s.filterMock,
	}

	err := testContrObj.handleMarketsRequest(s.ctx)

	assert.NoError(s.T(), err)

	assert.Equal(s.T(), http.StatusBadRequest, s.ctx.Response().Status)
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

func (s *handlersTestSuite) TestHandleBestPriceRequestNegativeCaseFilterInteraction() {
	var testContrObj = Controller{
		logger: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})),
		filter: s.filterMock,
	}

	err := testContrObj.handleBestPriceRequest(s.ctx)

	assert.NoError(s.T(), err)

	assert.Equal(s.T(), http.StatusInternalServerError, s.ctx.Response().Status)
}

func (s *handlersTestSuite) TestHandleBestPriceRequestNegativeInputCase() {
	var testContrObj = Controller{
		logger: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})),
		filter: s.filterMock,
	}

	err := testContrObj.handleBestPriceRequest(s.ctx)

	assert.NoError(s.T(), err)

	assert.Equal(s.T(), http.StatusBadRequest, s.ctx.Response().Status)
}

func TestContollerHandlers(t *testing.T) {
	suite.Run(t, new(handlersTestSuite))
}
