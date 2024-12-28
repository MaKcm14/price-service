package services

import "errors"

var (
	ErrGettingProducts = errors.New("error of getting the products from all the market/markets")
	ErrChooseMarket    = errors.New("error of choosing the market api")
)
