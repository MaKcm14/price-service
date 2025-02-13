package chttp

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/MaKcm14/price-service/internal/entities/dto"
	"github.com/MaKcm14/price-service/pkg/entities"
)

type queryOpt func(ctx echo.Context, request *dto.ProductRequest) error

type (
	checker struct{}

	validator struct {
		check checker
	}
)

// isDataSafe checks the data doesn't have the SQL-injection signatures.
func (c checker) isDataSafe(data string) bool {
	data = strings.ToLower(data)

	for _, elem := range data {
		if string(elem) == "=" {
			return false
		}
	}

	if strings.Contains(data, "drop ") || (strings.Contains(data, "union ") && strings.Contains(data, "select")) || strings.Contains(data, "--") {
		return false
	}

	return true
}

// validQuery validates the product query.
func (v validator) validQuery(ctx echo.Context, request *dto.ProductRequest) error {
	if query := ctx.QueryParam("query"); v.check.isDataSafe(query) && len(query) != 0 {
		request.Query, _ = url.QueryUnescape(query)
		return nil
	}
	return ErrRequestInfo
}

// validSample validates the param "sample" that defines the num of the products' sample.
func (v validator) validSample(ctx echo.Context, request *dto.ProductRequest) error {
	sample, err := strconv.Atoi(ctx.QueryParam("sample"))

	if sample < 0 || err != nil {
		sample = 1
	}
	request.Sample = sample

	return nil
}

// validMarkets validates the param "markets" that defines the markets where products will be searched.
func (v validator) validMarkets(ctx echo.Context, request *dto.ProductRequest) error {
	wildbFlag := false
	mmegaFlag := false

	for _, market := range strings.Split(ctx.QueryParam("markets"), " ") {
		if market == "wildberries" && !wildbFlag {
			request.Markets = append(request.Markets, entities.Wildberries)
			wildbFlag = true
		} else if market == "megamarket" && !mmegaFlag {
			request.Markets = append(request.Markets, entities.MegaMarket)
			mmegaFlag = true
		}
	}

	if len(request.Markets) == 0 {
		return ErrRequestInfo
	}

	return nil
}

// validAmount validates the param "amount" that defines the amount of products.
func (v validator) validAmount(ctx echo.Context, request *dto.ProductRequest) error {
	request.Amount = ctx.QueryParam("amount")

	if amt := request.Amount; amt != "max" && amt != "min" {
		request.Amount = "min"
	}

	return nil
}

// validSort validates the param "sort" that defines the sort of sample.
func (v validator) validSort(ctx echo.Context, request *dto.ProductRequest) error {
	request.Sort = dto.SortType(ctx.QueryParam("sort"))

	if sort := request.Sort; sort != dto.PopularSort && sort != dto.PriceDownSort && sort != dto.PriceUpSort &&
		sort != dto.RateSort && sort != dto.NewlySort {
		request.Sort = dto.PopularSort
	}

	return nil
}

// validNoImage validates the param "no-image" that defines the presense the image-links in
// the response.
func (v validator) validNoImage(ctx echo.Context, request *dto.ProductRequest) error {
	flagNoImage := ctx.QueryParam("no-image")

	request.FlagNoImage = true

	if flagNoImage == "0" {
		request.FlagNoImage = false
	}

	return nil
}

// validProductRequest validates the info from the URL-query's params.
func (v validator) validProductRequest(ctx echo.Context, opts ...queryOpt) (dto.ProductRequest, error) {
	var request dto.ProductRequest

	for _, opt := range opts {
		err := opt(ctx, &request)

		if err != nil {
			return dto.ProductRequest{}, err
		}
	}

	return request, nil
}
