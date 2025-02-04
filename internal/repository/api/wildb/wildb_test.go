package wildb

import (
	"log/slog"
	"os"
	"testing"

	"github.com/MaKcm14/best-price-service/price-service/internal/entities/dto"
	"github.com/stretchr/testify/assert"
)

func TestParseImageLinksPositiveCase(t *testing.T) {
	var testParserObj = wildberriesParser{
		logger: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})),
	}

	correctPage, err := os.Open("../../../../test/repository/api/wildb/positive_case.html")

	if err != nil {
		t.Fatal("error of test configuration: couldn't open the test html page")
	}
	defer correctPage.Close()

	html := make([]byte, 13000)
	correctPage.Read(html)

	resImageLinks := testParserObj.parseImageLinks(string(html))

	assert.Equal(t, 2, len(resImageLinks))
	assert.Equal(t, []string{
		"https://basket-16.wbbasket.ru/vol2425/part242589/242589892/images/c516x688/1.webp",
		"https://basket-16.wbbasket.ru/vol2611/part261162/261162615/images/c516x688/1.webp",
	}, resImageLinks)
}

func TestParseImageLinksExtremeCases(t *testing.T) {
	t.Run("Extreme Case: there isn't any article in the page", func(t *testing.T) {
		var testParserObj = wildberriesParser{
			logger: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})),
		}

		correctPage, err := os.Open("../../../../test/repository/api/wildb/extreme_case_no_products.html")

		if err != nil {
			t.Fatal("error of test configuration: couldn't open the test html page")
		}
		defer correctPage.Close()

		html := make([]byte, 13000)
		correctPage.Read(html)

		resImageLinks := testParserObj.parseImageLinks(string(html))

		assert.Equal(t, 0, len(resImageLinks))
		assert.Equal(t, []string(nil), resImageLinks)
	})

	t.Run("Extreme Case: the wrong structure of the page was given", func(t *testing.T) {
		var testParserObj = wildberriesParser{
			logger: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})),
		}

		correctPage, err := os.Open("../../../../test/repository/api/wildb/extreme_case_wrong_struct.html")

		if err != nil {
			t.Fatal("error of test configuration: couldn't open the test html page")
		}
		defer correctPage.Close()

		html := make([]byte, 13000)
		correctPage.Read(html)

		resImageLinks := testParserObj.parseImageLinks(string(html))

		assert.Equal(t, 0, len(resImageLinks))
		assert.Equal(t, []string{}, resImageLinks)
	})
}

func TestGetOpenApiPathPositiveCases(t *testing.T) {
	t.Run("Positive Case: path without price-range filter", func(t *testing.T) {
		var testViewObj = wildberriesViewer{}

		url := testViewObj.getOpenApiPath(dto.ProductRequest{
			Sample: 1,
			Query:  "Test Query",
		}, []string{"sort", "popular"})

		assert.Equal(t, "page=1&sort=popular&search=Test+Query", url)
	})

	t.Run("Positive Case: path with price-range filter", func(t *testing.T) {
		var testViewObj = wildberriesViewer{}

		url := testViewObj.getOpenApiPath(dto.ProductRequest{
			Sample: 1,
			Query:  "Test Query",
		}, []string{"sort", "popular", "priceU", "100000;500000"})

		assert.Equal(t, "page=1&sort=popular&priceU=100000;500000&search=Test+Query", url)
	})
}

func TestGetHiddenApiPathPositiveCases(t *testing.T) {
	t.Run("Positive Case: path without price-range filter", func(t *testing.T) {
		var testViewObj = wildberriesViewer{}

		url := testViewObj.getHiddenApiPath(dto.ProductRequest{
			Sample: 1,
			Query:  "Test Query",
		}, []string{"sort", "popular"})

		assert.Equal(t, "page=1&query=Test+Query&resultset=catalog&sort=popular&spp=30&suppressSpellcheck=false", url)
	})

	t.Run("Positive Case: path with price-range filter", func(t *testing.T) {
		var testViewObj = wildberriesViewer{}

		url := testViewObj.getHiddenApiPath(dto.ProductRequest{
			Sample: 1,
			Query:  "Test Query",
		}, []string{"sort", "popular", "priceU", "100000;500000"})

		assert.Equal(t, "page=1&priceU=100000;500000&query=Test+Query&resultset=catalog&sort=popular&spp=30&suppressSpellcheck=false", url)
	})
}
