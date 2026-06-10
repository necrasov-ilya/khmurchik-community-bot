package di

import (
	"github.com/evart2006/khmurchik-community-bot/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Container struct {
	Config *config.Config
	DB     *pgxpool.Pool
	Logger *zap.Logger
	Bot    *tgbotapi.BotAPI
}
