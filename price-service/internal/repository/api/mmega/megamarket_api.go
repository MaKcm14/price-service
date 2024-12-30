package mmega

import (
	"context"
	"log/slog"
	"time"
)

type MegaMarketAPI struct {
	logger    *slog.Logger
	loadCoeff time.Duration
	ctx       context.Context
}

func NewMegaMarketAPI(ctx context.Context, log *slog.Logger, loadCoeff int) MegaMarketAPI {
	return MegaMarketAPI{
		logger:    log,
		ctx:       ctx,
		loadCoeff: time.Duration(loadCoeff) * time.Millisecond,
	}
}
