package greeting

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/evart2006/khmurchik-community-bot/internal/middleware"
	"github.com/evart2006/khmurchik-community-bot/internal/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type Service struct {
	api     *tgbotapi.BotAPI
	chats   *repository.ChatRepository
	checker *middleware.AdminChecker
	logger  *zap.Logger
}

func NewService(api *tgbotapi.BotAPI, chats *repository.ChatRepository, checker *middleware.AdminChecker, logger *zap.Logger) *Service {
	return &Service{api: api, chats: chats, checker: checker, logger: logger}
}

func (s *Service) GreetingOnHandler() func(tgbotapi.Message) {
	return func(msg tgbotapi.Message) {
		if !s.ensureAdmin(msg) {
			return
		}
		if err := s.chats.SetGreetingEnabled(context.Background(), msg.Chat.ID, true); err != nil {
			s.reply(msg, fmt.Sprintf("Ошибка настройки приветствия: %s", err.Error()))
			return
		}
		s.reply(msg, "Утреннее приветствие включено.")
	}
}

func (s *Service) GreetingOffHandler() func(tgbotapi.Message) {
	return func(msg tgbotapi.Message) {
		if !s.ensureAdmin(msg) {
			return
		}
		if err := s.chats.SetGreetingEnabled(context.Background(), msg.Chat.ID, false); err != nil {
			s.reply(msg, fmt.Sprintf("Ошибка настройки приветствия: %s", err.Error()))
			return
		}
		s.reply(msg, "Утреннее приветствие выключено.")
	}
}

func (s *Service) GreetingTimeHandler() func(tgbotapi.Message) {
	return func(msg tgbotapi.Message) {
		if !s.ensureAdmin(msg) {
			return
		}
		args := strings.Fields(msg.CommandArguments())
		if len(args) != 2 {
			s.reply(msg, "Формат: /greeting_time 10:00 Europe/Minsk")
			return
		}
		if _, err := time.Parse("15:04", args[0]); err != nil {
			s.reply(msg, "Время должно быть в формате HH:MM, например 10:00.")
			return
		}
		if _, err := time.LoadLocation(args[1]); err != nil {
			s.reply(msg, "Не знаю такую таймзону. Пример: Europe/Minsk.")
			return
		}
		if err := s.chats.SetGreetingTime(context.Background(), msg.Chat.ID, args[0], args[1]); err != nil {
			s.reply(msg, fmt.Sprintf("Ошибка настройки времени: %s", err.Error()))
			return
		}
		s.reply(msg, fmt.Sprintf("Время приветствия: %s %s.", args[0], args[1]))
	}
}

func (s *Service) GreetingTextHandler() func(tgbotapi.Message) {
	return func(msg tgbotapi.Message) {
		if !s.ensureAdmin(msg) {
			return
		}
		text := strings.TrimSpace(msg.CommandArguments())
		if text == "" {
			s.reply(msg, "Формат: /greeting_text текст утреннего сообщения")
			return
		}
		if err := s.chats.SetGreetingMessage(context.Background(), msg.Chat.ID, text); err != nil {
			s.reply(msg, fmt.Sprintf("Ошибка настройки текста: %s", err.Error()))
			return
		}
		s.reply(msg, "Текст утреннего приветствия обновлён.")
	}
}

func (s *Service) GreetingStatusHandler() func(tgbotapi.Message) {
	return func(msg tgbotapi.Message) {
		if !s.ensureAdmin(msg) {
			return
		}
		settings, err := s.chats.GetSettings(context.Background(), msg.Chat.ID)
		if err != nil {
			s.reply(msg, fmt.Sprintf("Не смог получить настройки: %s", err.Error()))
			return
		}
		state := "выключено"
		if settings.GreetingEnabled {
			state = "включено"
		}
		s.reply(msg, fmt.Sprintf("Приветствие: %s\nВремя: %s %s\nТекст: %s", state, settings.GreetingTime, settings.GreetingTimezone, settings.GreetingMessage))
	}
}

func (s *Service) ensureAdmin(msg tgbotapi.Message) bool {
	if msg.From == nil {
		return false
	}
	if msg.Chat.Type != "group" && msg.Chat.Type != "supergroup" {
		s.reply(msg, "Настройки приветствий работают только в группах.")
		return false
	}
	ok, err := s.checker.Check(msg.Chat.ID, msg.From.ID)
	if err != nil {
		s.logger.Warn("admin check failed", zap.Error(err), zap.Int64("chat_id", msg.Chat.ID), zap.Int64("user_id", msg.From.ID))
		s.reply(msg, "Не смог проверить права администратора.")
		return false
	}
	if !ok {
		s.reply(msg, "У вас нет прав для этой команды.")
		return false
	}
	return true
}

func (s *Service) reply(msg tgbotapi.Message, text string) {
	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	reply.ReplyToMessageID = msg.MessageID
	if _, err := s.api.Send(reply); err != nil {
		s.logger.Warn("send greeting reply failed", zap.Error(err), zap.Int64("chat_id", msg.Chat.ID))
	}
}
