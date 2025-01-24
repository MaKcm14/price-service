package wildb

import (
	"log/slog"
	"os"
	"testing"

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
