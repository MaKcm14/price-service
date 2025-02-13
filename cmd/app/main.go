package main

import (
	"github.com/MaKcm14/price-service/internal/app"
)

func main() {
	service := app.NewService()
	service.Run()
}
