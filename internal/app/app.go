package app

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/evart2006/khmurchik-community-bot/internal/bot"
	"github.com/evart2006/khmurchik-community-bot/internal/config"
	"github.com/evart2006/khmurchik-community-bot/internal/handlers"
	"github.com/evart2006/khmurchik-community-bot/internal/modules/greeting"
	"github.com/evart2006/khmurchik-community-bot/internal/modules/moderation"
	"github.com/evart2006/khmurchik-community-bot/internal/modulesystem"
	"github.com/evart2006/khmurchik-community-bot/internal/repository"
	"go.uber.org/zap"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type App struct {
	cfg    *config.Config
	db     *pgxpool.Pool
	logger *zap.Logger
	bot    *tgbotapi.BotAPI
	router *bot.Router
	reg    *modulesystem.Registry
}

func New(cfg *config.Config, db *pgxpool.Pool, logger *zap.Logger) *App {
	return &App{cfg: cfg, db: db, logger: logger}
}

func (a *App) Start(ctx context.Context) error {
	b, err := tgbotapi.NewBotAPI(a.cfg.Bot.Token)
	if err != nil {
		return err
	}
	a.bot = b

	a.router = bot.NewRouter()
	a.router.Register("start", handlers.StartHandler(a.bot))
	a.router.SetUnknownHandler(handlers.UnknownHandler(a.bot))

	a.reg = modulesystem.NewRegistry()

	greetingMod := greeting.NewModule(a.bot, a.cfg.Bot.TargetChatID, a.logger)
	a.reg.Register(greetingMod)

	moderationMod := moderation.NewModule(a.bot, a.cfg.Bot.TargetChatID, a.db, a.logger)
	a.reg.Register(moderationMod)

	if err := a.reg.RegisterAll(a.bot, a.router); err != nil {
		return err
	}

	if err := a.reg.StartAll(ctx); err != nil {
		return err
	}

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = int(a.cfg.Server.PollTimeout)
	updates := a.bot.GetUpdatesChan(updateConfig)

	userRepo := repository.NewUserRepository(a.db)
	chatRepo := repository.NewChatRepository(a.db)

	go func() {
		for update := range updates {
			if update.Message != nil && update.Message.From != nil {
				if err := userRepo.UpsertUser(update.Message.From); err != nil {
					a.logger.Error("upsert user failed", zap.Error(err), zap.Int64("user_id", update.Message.From.ID))
				}
				if err := chatRepo.UpsertChat(*update.Message.Chat); err != nil {
					a.logger.Error("upsert chat failed", zap.Error(err), zap.Int64("chat_id", update.Message.Chat.ID))
				}
			}
			a.router.Handle(update)
		}
	}()

	<-ctx.Done()

	a.reg.StopAll(context.Background())
	a.bot.StopReceivingUpdates()

	return nil
}
