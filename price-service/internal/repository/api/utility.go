package api

import (
	"github.com/labstack/echo/v4"
)

const (
	Wildberries string = "Wildberries"
)

func IsConnectionClosed(ctx echo.Context) bool {
	select {
	case <-ctx.Request().Context().Done():
		return true

	default:
		return false
	}
}
