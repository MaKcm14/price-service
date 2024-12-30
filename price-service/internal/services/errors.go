package services

import "errors"

var (
	ErrGettingProducts = errors.New("error of getting the products from all the market/markets")
	ErrMarketApi       = errors.New("error of getting the market api")
)
