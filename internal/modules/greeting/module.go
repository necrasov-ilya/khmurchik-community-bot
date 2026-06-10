package greeting

import (
	"context"

	"github.com/evart2006/khmurchik-community-bot/internal/bot"
	"go.uber.org/zap"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Module struct {
	bot          *tgbotapi.BotAPI
	targetChatID int64
	logger       *zap.Logger
}

func NewModule(bot *tgbotapi.BotAPI, targetChatID int64, logger *zap.Logger) *Module {
	return &Module{bot: bot, targetChatID: targetChatID, logger: logger}
}

func (m *Module) Name() string    { return "greeting" }
func (m *Module) Version() string { return "1.0.0" }

func (m *Module) Register(bot *tgbotapi.BotAPI, router *bot.Router) error {
	return nil
}

func (m *Module) OnStart(ctx context.Context) error {
	sched := NewScheduler(m.bot, m.targetChatID, m.logger)
	return sched.Start(ctx)
}

func (m *Module) OnStop(ctx context.Context) error {
	return nil
}
