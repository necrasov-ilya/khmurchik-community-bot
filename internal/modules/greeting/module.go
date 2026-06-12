package greeting

import (
	"context"

	"github.com/evart2006/khmurchik-community-bot/internal/bot"
	"github.com/evart2006/khmurchik-community-bot/internal/config"
	"github.com/evart2006/khmurchik-community-bot/internal/middleware"
	"github.com/evart2006/khmurchik-community-bot/internal/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Module struct {
	bot       *tgbotapi.BotAPI
	db        *pgxpool.Pool
	defaults  config.DefaultsConfig
	logger    *zap.Logger
	scheduler *Scheduler
}

func NewModule(bot *tgbotapi.BotAPI, db *pgxpool.Pool, defaults config.DefaultsConfig, logger *zap.Logger) *Module {
	return &Module{bot: bot, db: db, defaults: defaults, logger: logger}
}

func (m *Module) Name() string    { return "greeting" }
func (m *Module) Version() string { return "1.0.0" }

func (m *Module) Register(api *tgbotapi.BotAPI, router *bot.Router) error {
	chatRepo := repository.NewChatRepository(m.db, m.defaults)
	checker := middleware.NewAdminChecker(api, m.logger)
	svc := NewService(api, chatRepo, checker, m.logger)

	router.Register("greeting_on", svc.GreetingOnHandler())
	router.Register("greeting_off", svc.GreetingOffHandler())
	router.Register("greeting_time", svc.GreetingTimeHandler())
	router.Register("greeting_text", svc.GreetingTextHandler())
	router.Register("greeting_status", svc.GreetingStatusHandler())
	return nil
}

func (m *Module) OnStart(ctx context.Context) error {
	m.scheduler = NewScheduler(m.bot, repository.NewChatRepository(m.db, m.defaults), m.logger)
	return m.scheduler.Start(ctx)
}

func (m *Module) OnStop(ctx context.Context) error {
	if m.scheduler != nil {
		m.scheduler.Stop()
	}
	return nil
}
