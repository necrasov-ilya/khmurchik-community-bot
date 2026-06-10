package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/evart2006/khmurchik-community-bot/internal/app"
	"github.com/evart2006/khmurchik-community-bot/internal/config"
	"github.com/evart2006/khmurchik-community-bot/internal/database"
	"github.com/evart2006/khmurchik-community-bot/internal/logger"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		panic(err)
	}

	l, err := logger.New()
	if err != nil {
		panic(err)
	}
	defer l.Sync()

	pool, err := database.NewPool(cfg.Database)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	a := app.New(cfg, pool, l)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		cancel()
	}()

	if err := a.Start(ctx); err != nil {
		l.Fatal("app start failed", zap.Error(err))
	}
}
