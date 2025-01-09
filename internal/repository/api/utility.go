package api

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/chromedp/chromedp"
	"github.com/labstack/echo/v4"
)

// ChromePull supports the safe opening and closing the connection with the instances of the browser.
type ChromePull struct {
	cancels []context.CancelFunc
}

func NewChromePull() *ChromePull {
	return &ChromePull{
		cancels: make([]context.CancelFunc, 0, 100),
	}
}

// NewContext creates a new context.Context that connected with the allocated browser.
func (c *ChromePull) NewContext() context.Context {
	ctx, cancel := chromedp.NewExecAllocator(context.Background(),
		append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))...)
	c.cancels = append(c.cancels, cancel)

	ctx, cancel = chromedp.NewContext(ctx)
	c.cancels = append(c.cancels, cancel)

	return ctx
}

// Close closes all the connections that were created with the NewContext.
func (c ChromePull) Close() {
	for _, cancel := range c.cancels {
		cancel()
	}
}

// IsConnectionClosed checks is the connection with the client is still alive.
func IsConnectionClosed(ctx echo.Context) bool {
	select {
	case <-ctx.Request().Context().Done():
		return true

	default:
		return false
	}
}

func ReadResponseBody(source io.Reader, logger *slog.Logger, serviceType string) ([]byte, error) {
	respBody := make([]byte, 0, 100000)
	for {
		buffer := make([]byte, 10000)
		n, err := source.Read(buffer)

		if n != 0 && (err == nil || err == io.EOF) {
			respBody = append(respBody, buffer[:n]...)
		} else if err != nil && err != io.EOF {
			logger.Warn(fmt.Sprintf("error of the %v: %v: %v", serviceType, ErrBufferReading, err))
			return nil, fmt.Errorf("%w: %v", ErrBufferReading, err)
		} else if err == io.EOF {
			break
		}
	}

	return respBody, nil
}
