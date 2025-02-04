package mmega

import (
	"fmt"
	"testing"

	"github.com/MaKcm14/best-price-service/price-service/internal/entities/dto"
	"github.com/stretchr/testify/assert"
)

func TestGetOpenApiURLPositiveCases(t *testing.T) {
	t.Run("Positive Case: URL without price-range", func(t *testing.T) {
		var testViewerObj = megaMarketViewer{}

		url := testViewerObj.getOpenApiURL(dto.ProductRequest{
			Sample: 1,
			Query:  "Test query",
		}, []string{"sort", "0"})

		assert.Equal(t, "https://megamarket.ru/catalog/page-1/?q=Test+query#?sort=0", url)
	})

	t.Run("Positive Case: URL with price-range", func(t *testing.T) {
		var testViewerObj = megaMarketViewer{}

		url := testViewerObj.getOpenApiURL(dto.ProductRequest{
			Sample: 1,
			Query:  "Test query",
		}, []string{"sort", "0", "filters", "2000 10000"})

		assert.Equal(t,
			fmt.Sprintf("https://megamarket.ru/catalog/page-1/?q=Test+query#?sort=0&filters=%s", testViewerObj.getPriceRangeURLView("2000 10000")),
			url)
	})
}
