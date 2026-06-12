package reports

import (
	"context"

	"github.com/evart2006/khmurchik-community-bot/internal/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Module struct {
	bot    *tgbotapi.BotAPI
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewModule(bot *tgbotapi.BotAPI, db *pgxpool.Pool, logger *zap.Logger) *Module {
	return &Module{bot: bot, db: db, logger: logger}
}

func (m *Module) Name() string    { return "reports" }
func (m *Module) Version() string { return "1.0.0" }

func (m *Module) Register(bot *tgbotapi.BotAPI, router *bot.Router) error {
	svc := NewService(bot, m.db, m.logger)
	router.Register("report", svc.ReportHandler(bot))
	return nil
}

func (m *Module) OnStart(ctx context.Context) error { return nil }
func (m *Module) OnStop(ctx context.Context) error  { return nil }
