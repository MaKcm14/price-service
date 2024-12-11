package main

import "github.com/MaKcm14/best-price-service/price-service/internal/app"

func main() {
	app := app.NewApp()
	app.Run()
}
