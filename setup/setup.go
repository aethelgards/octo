package setup

import (
	"context"
	"log/slog"

	"github.com/aethelgards/octo/config"
	"github.com/aethelgards/octo/llm"
	"github.com/aethelgards/octo/tui"
)

func Init(ctx context.Context) {
	err := initInner(ctx)
	if err != nil {
		slog.Error("init octo failed", slog.String("error", err.Error()))
		panic(err)
	}
}

func initInner(ctx context.Context) error {
	if err := config.LoadConfig(ctx); err != nil {
		return err
	}
	if err := InitLog(&config.OctoConfig.LogConfig); err != nil {
		return err
	}
	if err := llm.InitModel(ctx, config.OctoConfig.LLMConfig); err != nil {
		return err
	}
	if err := llm.InitAgent(ctx); err != nil {
		return err
	}
	if err := tui.Init(ctx, config.OctoConfig); err != nil {
		return err
	}
	return nil
}
