package api_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/MaKcm14/best-price-service/price-service/internal/repository/api"
)

func TestGetFiltersPositiveCases(t *testing.T) {
	var testConvertObj = api.URLConverter{}

	res := testConvertObj.GetFilters([]string{"filter_1", "value_1", "filter_2", "value_2"})

	assert.Equal(t, res, map[string]string{
		"filter_1": "value_1",
		"filter_2": "value_2",
	})
}

func TestExtremeCases(t *testing.T) {
	t.Run("Extreme: the filter's slice is nil", func(t *testing.T) {
		var testConvertObj = api.URLConverter{}

		res := testConvertObj.GetFilters(nil)

		assert.Equal(t, map[string]string{}, res)
	})

	t.Run("Extreme: the filter's slice has the len = 0", func(t *testing.T) {
		var testConvertObj = api.URLConverter{}

		res := testConvertObj.GetFilters([]string{})

		assert.Equal(t, map[string]string{}, res)
	})
}

func TestGetFilterNegativeCases(t *testing.T) {
	t.Run("Negative: the filter's slice has the non-even len", func(t *testing.T) {
		var testConvertObj = api.URLConverter{}

		res := testConvertObj.GetFilters([]string{"filter_1", "value_1", "filter_2"})

		assert.Equal(t, map[string]string{
			"filter_1": "value_1",
		}, res)
	})
}
