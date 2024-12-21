package services

import "github.com/MaKcm14/best-price-service/price-service/internal/entities"

type ApiInteractor interface {
	GetProducts(product entities.ProductRequest)
}
