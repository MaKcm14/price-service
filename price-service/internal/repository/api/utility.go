package api

import (
	"github.com/labstack/echo/v4"
)

const (
	Wildberries string = "Wildberries"
)

// isConnectionClosed checks is the connection with the client is still alive.
func isConnectionClosed(ctx echo.Context) bool {
	select {
	case <-ctx.Request().Context().Done():
		return true

	default:
		return false
	}
}
