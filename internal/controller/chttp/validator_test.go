package chttp

import (
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/MaKcm14/best-price-service/price-service/internal/entities/dto"
	"github.com/MaKcm14/best-price-service/price-service/pkg/entities"
)

func TestValidMarketsPositiveCase(t *testing.T) {
	var testValidatorObj = validator{}
	var testRequestObj = dto.ProductRequest{}

	request, err := http.NewRequest("GET", "http://localhost/products/filter/markets?query=test+query&markets=wildberries+megamarket", nil)

	if err != nil {
		t.Fatal("error of the test configuration")
	}

	err = testValidatorObj.validMarkets(echo.New().NewContext(request, nil), &testRequestObj)

	if assert.NoError(t, err) {
		assert.Equal(t, []entities.Market{entities.Wildberries, entities.MegaMarket}, testRequestObj.Markets)
	}
}

func TestValidMarketsExtremeCase(t *testing.T) {
	var testValidatorObj = validator{}
	var testRequestObj = dto.ProductRequest{}

	request, err := http.NewRequest("GET", "http://localhost/products/filter/markets?query=test+query&markets=wildberries+megamarket+wildberries", nil)

	if err != nil {
		t.Fatal("error of the test configuration")
	}

	err = testValidatorObj.validMarkets(echo.New().NewContext(request, nil), &testRequestObj)

	if assert.NoError(t, err) {
		assert.Equal(t, []entities.Market{entities.Wildberries, entities.MegaMarket}, testRequestObj.Markets)
	}
}

func TestValidMarketsNegativeCases(t *testing.T) {
	t.Run("Negative Case: error of the markets' names", func(t *testing.T) {
		var testValidatorObj = validator{}
		var testRequestObj = dto.ProductRequest{}

		request, err := http.NewRequest("GET", "http://localhost/products/filter/markets?query=test+query&markets=wildslbareies", nil)

		if err != nil {
			t.Fatal("error of the test configuration")
		}

		err = testValidatorObj.validMarkets(echo.New().NewContext(request, nil), &testRequestObj)

		if assert.Error(t, err) {
			assert.Equal(t, 0, len(testRequestObj.Markets))
		}
	})

	t.Run("Negative Case: error of the markets: the absence of the markets' param", func(t *testing.T) {
		var testValidatorObj = validator{}
		var testRequestObj = dto.ProductRequest{}

		request, err := http.NewRequest("GET", "http://localhost/products/filter/markets?query=test+query", nil)

		if err != nil {
			t.Fatal("error of the test configuration")
		}

		err = testValidatorObj.validMarkets(echo.New().NewContext(request, nil), &testRequestObj)

		if assert.Error(t, err) {
			assert.Equal(t, 0, len(testRequestObj.Markets))
		}
	})
}

func TestIsDataSafePositiveCases(t *testing.T) {
	t.Run("Clear data: without any SQL-injection string", func(t *testing.T) {
		var testCheckerObj = checker{}

		flagDataSafe := testCheckerObj.isDataSafe("Hello, World!")

		assert.True(t, flagDataSafe)
	})
}

func TestIsDataSafeExtremeCases(t *testing.T) {
	t.Run("Extreme data: requests that have SQL instructions but aren't really this: 'Union Case'", func(t *testing.T) {
		var testCheckerObj = checker{}

		flagDataSafe := testCheckerObj.isDataSafe("Union Jack")

		assert.True(t, flagDataSafe)
	})

	t.Run("Extreme data: requests that have SQL instructions but aren't really this: 'Drop Case'", func(t *testing.T) {
		var testCheckerObj = checker{}

		flagDataSafe := testCheckerObj.isDataSafe("Dropping watering")

		assert.True(t, flagDataSafe)
	})
}

func TestIsDataSafeNegativeCases(t *testing.T) {
	t.Run("Negative data: requests that are the SQL-injections: 'Drop Case'", func(t *testing.T) {
		var testCheckerObj = checker{}

		flagDataSafe := testCheckerObj.isDataSafe("DROP database")

		assert.False(t, flagDataSafe)

		flagDataSafe = testCheckerObj.isDataSafe("' and 1=1--")

		assert.False(t, flagDataSafe)
	})

	t.Run("Negative data: requests that are the SQL-injections: 'Union Case'", func(t *testing.T) {
		var testCheckerObj = checker{}

		flagDataSafe := testCheckerObj.isDataSafe("SELECT data from table UNION SELECT data from table")

		assert.False(t, flagDataSafe)
	})

	t.Run("Negative data: requests that are the SQL-injections: 'Comment Case'", func(t *testing.T) {
		var testCheckerObj = checker{}

		flagDataSafe := testCheckerObj.isDataSafe("' and 1=1--")

		assert.False(t, flagDataSafe)
	})
}
