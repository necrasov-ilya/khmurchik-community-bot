package modulesystem

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/evart2006/khmurchik-community-bot/internal/bot"
)

type Module interface {
	Name() string
	Version() string
	Register(bot *tgbotapi.BotAPI, router *bot.Router) error
	OnStart(ctx context.Context) error
	OnStop(ctx context.Context) error
}
