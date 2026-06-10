package moderation

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/evart2006/khmurchik-community-bot/internal/bot"
	"github.com/evart2006/khmurchik-community-bot/internal/middleware"
	"go.uber.org/zap"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Module struct {
	bot          *tgbotapi.BotAPI
	targetChatID int64
	db           *pgxpool.Pool
	logger       *zap.Logger
}

func NewModule(bot *tgbotapi.BotAPI, targetChatID int64, db *pgxpool.Pool, logger *zap.Logger) *Module {
	return &Module{bot: bot, targetChatID: targetChatID, db: db, logger: logger}
}

func (m *Module) Name() string    { return "moderation" }
func (m *Module) Version() string { return "1.0.0" }

func (m *Module) Register(bot *tgbotapi.BotAPI, router *bot.Router) error {
	checker := middleware.NewAdminChecker(bot, m.targetChatID, m.logger)
	svc := NewService(bot, m.db, m.logger, checker, m.targetChatID)

	router.Register("mute", svc.MuteHandler(bot))
	router.Register("ban", svc.BanHandler(bot))
	router.Register("kick", svc.KickHandler(bot))
	router.Register("unmute", svc.UnmuteHandler(bot))

	return nil
}

func (m *Module) OnStart(ctx context.Context) error { return nil }
func (m *Module) OnStop(ctx context.Context) error  { return nil }
