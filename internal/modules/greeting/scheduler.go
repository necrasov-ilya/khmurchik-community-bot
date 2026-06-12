package greeting

import (
	"context"
	"strings"
	"time"

	"github.com/evart2006/khmurchik-community-bot/internal/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type Scheduler struct {
	bot    *tgbotapi.BotAPI
	chats  *repository.ChatRepository
	logger *zap.Logger
	cancel context.CancelFunc
}

func NewScheduler(bot *tgbotapi.BotAPI, chats *repository.ChatRepository, logger *zap.Logger) *Scheduler {
	return &Scheduler{bot: bot, chats: chats, logger: logger}
}

func (s *Scheduler) Start(ctx context.Context) error {
	runCtx, cancel := context.WithCancel(ctx)
	s.cancel = cancel

	go s.loop(runCtx)
	return nil
}

func (s *Scheduler) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
}

func (s *Scheduler) loop(ctx context.Context) {
	s.tick(ctx, time.Now())

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			s.tick(ctx, now)
		}
	}
}

func (s *Scheduler) tick(ctx context.Context, now time.Time) {
	chats, err := s.chats.ListGreetingEnabled(ctx)
	if err != nil {
		s.logger.Error("list greeting chats failed", zap.Error(err))
		return
	}

	for _, chat := range chats {
		shouldSend, localNow, err := ShouldSendGreeting(chat, now)
		if err != nil {
			s.logger.Warn("invalid greeting settings", zap.Error(err), zap.Int64("chat_id", chat.TelegramID))
			continue
		}
		if !shouldSend {
			continue
		}

		text := strings.TrimSpace(chat.GreetingMessage)
		if text == "" {
			continue
		}

		msg := tgbotapi.NewMessage(chat.TelegramID, text)
		if _, err := s.bot.Send(msg); err != nil {
			s.logger.Error("failed to send greeting", zap.Error(err), zap.Int64("chat_id", chat.TelegramID))
			continue
		}

		if err := s.chats.MarkGreetingSent(ctx, chat.TelegramID, localNow); err != nil {
			s.logger.Error("mark greeting sent failed", zap.Error(err), zap.Int64("chat_id", chat.TelegramID))
		}
	}
}

func ShouldSendGreeting(chat repository.ChatSettings, now time.Time) (bool, time.Time, error) {
	loc, err := time.LoadLocation(chat.GreetingTimezone)
	if err != nil {
		return false, time.Time{}, err
	}

	localNow := now.In(loc)
	target, err := time.Parse("15:04", chat.GreetingTime)
	if err != nil {
		return false, time.Time{}, err
	}

	if localNow.Hour() != target.Hour() || localNow.Minute() != target.Minute() {
		return false, localNow, nil
	}
	if chat.LastGreetingSentDate != nil && sameLocalDate(*chat.LastGreetingSentDate, localNow, loc) {
		return false, localNow, nil
	}
	return true, localNow, nil
}

func sameLocalDate(a time.Time, b time.Time, loc *time.Location) bool {
	aa := a.In(loc)
	bb := b.In(loc)
	return aa.Year() == bb.Year() && aa.YearDay() == bb.YearDay()
}
