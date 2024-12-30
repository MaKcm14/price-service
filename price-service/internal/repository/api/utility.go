package api

import (
	"context"

	"github.com/chromedp/chromedp"
	"github.com/labstack/echo/v4"
)

const (
	Wildberries string = "Wildberries"
)

type ChromePull struct {
	cancels []context.CancelFunc
}

func NewChromePull() ChromePull {
	return ChromePull{
		cancels: make([]context.CancelFunc, 0, 100),
	}
}

func (c ChromePull) NewContext() context.Context {
	ctx, cancel := chromedp.NewExecAllocator(context.Background(),
		append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))...)
	c.cancels = append(c.cancels, cancel)

	ctx, cancel = chromedp.NewContext(ctx)
	c.cancels = append(c.cancels, cancel)

	return ctx
}

func (c ChromePull) Close() {
	c.cancels[0]()
	c.cancels[1]()
}

// isConnectionClosed checks is the connection with the client is still alive.
func isConnectionClosed(ctx echo.Context) bool {
	select {
	case <-ctx.Request().Context().Done():
		return true

	default:
		return false
	}
}
