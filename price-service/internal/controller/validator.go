package controller

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/MaKcm14/best-price-service/price-service/internal/entities"
)

type queryOpts func(ctx echo.Context, request *entities.ProductRequest) error

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

	if strings.Contains(data, "drop ") || strings.Contains(data, "union ") || strings.Contains(data, "--") {
		return false
	}

	return true
}

// validQuery validates the product query.
func (v validator) validQuery(ctx echo.Context, request *entities.ProductRequest) error {
	if query := ctx.QueryParam("query"); v.check.isDataSafe(query) && len(query) != 0 {
		request.Query, _ = url.QueryUnescape(query)
		return nil
	}
	return ErrRequestInfo
}

// validSample validates the param "sample" that defines the num of the products' sample.
func (v validator) validSample(ctx echo.Context, request *entities.ProductRequest) error {
	sample, err := strconv.Atoi(ctx.QueryParam("sample"))

	if sample < 0 {
		return ErrRequestInfo
	} else if err != nil {
		return err
	}
	request.Sample = sample

	return nil
}

// validMarkets validates the param "markets" that defines the markets where products will be searched.
func (v validator) validMarkets(ctx echo.Context, request *entities.ProductRequest) error {
	for _, market := range strings.Split(ctx.QueryParam("markets"), " ") {
		if market == "wildberries" {
			request.Markets = append(request.Markets, entities.Wildberries)
		} else if market == "ozon" {
			request.Markets = append(request.Markets, entities.Ozon)
		} else if market == "megamarket" {
			request.Markets = append(request.Markets, entities.MegaMarket)
		}
	}

	if len(request.Markets) == 0 {
		return ErrRequestInfo
	}

	return nil
}

// validAmount validates the param "amount" that defines the amount of products.
func (v validator) validAmount(ctx echo.Context, request *entities.ProductRequest) error {
	request.Amount = ctx.QueryParam("amount")

	if amt := request.Amount; amt != "max" && amt != "min" {
		request.Amount = "min"
	}

	return nil
}

// validProductRequest validates the info from the URL-query's params.
func (v validator) validProductRequest(ctx echo.Context, opts ...queryOpts) (entities.ProductRequest, error) {
	var request entities.ProductRequest

	for _, opt := range opts {
		err := opt(ctx, &request)

		if err != nil {
			return entities.ProductRequest{}, err
		}
	}

	return request, nil
}
