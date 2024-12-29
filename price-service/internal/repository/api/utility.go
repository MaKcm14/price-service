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
	ctx     context.Context
	cancels []context.CancelFunc
}

func NewChromePull() ChromePull {
	cancels := make([]context.CancelFunc, 0, 2)

	ctx, cancel := chromedp.NewExecAllocator(context.Background(),
		append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))...)
	cancels = append(cancels, cancel)

	ctx, cancel = chromedp.NewContext(ctx)
	cancels = append(cancels, cancel)

	return ChromePull{
		ctx:     ctx,
		cancels: cancels,
	}
}

func (c ChromePull) GetContext() context.Context {
	return c.ctx
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
