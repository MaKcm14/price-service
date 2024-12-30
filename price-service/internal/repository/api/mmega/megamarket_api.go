package mmega

import (
	"context"
	"log/slog"
)

type MegaMarketAPI struct {
	logger    *slog.Logger
	loadCoeff float32
	ctx       context.Context
}

func NewMegaMarketAPI(ctx context.Context, log *slog.Logger, loadCoeff float32) MegaMarketAPI {
	return MegaMarketAPI{
		logger:    log,
		ctx:       ctx,
		loadCoeff: loadCoeff,
	}
}
