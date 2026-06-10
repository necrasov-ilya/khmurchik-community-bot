package greeting

import (
	"context"

	"github.com/robfig/cron/v3"
	"github.com/evart2006/khmurchik-community-bot/internal/timeutil"
	"go.uber.org/zap"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Scheduler struct {
	bot          *tgbotapi.BotAPI
	targetChatID int64
	logger       *zap.Logger
}

func NewScheduler(bot *tgbotapi.BotAPI, targetChatID int64, logger *zap.Logger) *Scheduler {
	return &Scheduler{bot: bot, targetChatID: targetChatID, logger: logger}
}

func (s *Scheduler) Start(ctx context.Context) error {
	loc, err := timeutil.LoadLocation()
	if err != nil {
		return err
	}
	c := cron.New(cron.WithLocation(loc))
	_, err = c.AddFunc("0 10 * * *", func() {
		sendDailyGreeting(s.bot, s.targetChatID, s.logger)
	})
	if err != nil {
		return err
	}
	c.Start()
	return nil
}
