package moderation

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/evart2006/khmurchik-community-bot/internal/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func (s *Service) MuteHandler(bot *tgbotapi.BotAPI) func(tgbotapi.Message) {
	return func(msg tgbotapi.Message) {
		if msg.From == nil {
			return
		}
		chatID := msg.Chat.ID

		if err := s.ensureAdmin(chatID, msg.From.ID); err != nil {
			s.reply(msg, bot, fmt.Sprintf("Ошибка: %s.", err.Error()))
			return
		}

		targetID, duration, reason, err := parseMuteArgs(msg.CommandArguments(), msg.ReplyToMessage, s.userRepo)
		if err != nil {
			s.reply(msg, bot, fmt.Sprintf("Ошибка: %s", err.Error()))
			return
		}

		if err := s.ensureRestrictAllowed(chatID, msg.From.ID, targetID); err != nil {
			s.reply(msg, bot, fmt.Sprintf("Ошибка: %s.", err.Error()))
			return
		}

		if err := s.executor.MuteUntil(chatID, targetID, duration); err != nil {
			s.reply(msg, bot, fmt.Sprintf("Ошибка мута: %s", err.Error()))
			return
		}

		s.logAction(msg.From.ID, targetID, chatID, "mute", reason, duration)
		s.reply(msg, bot, fmt.Sprintf("Пользователь замучен на %s", formatDuration(duration)))
	}
}

func (s *Service) BanHandler(bot *tgbotapi.BotAPI) func(tgbotapi.Message) {
	return func(msg tgbotapi.Message) {
		if msg.From == nil {
			return
		}
		chatID := msg.Chat.ID

		if err := s.ensureAdmin(chatID, msg.From.ID); err != nil {
			s.reply(msg, bot, fmt.Sprintf("Ошибка: %s.", err.Error()))
			return
		}

		targetID, reason, err := parseBanArgs(msg.CommandArguments(), msg.ReplyToMessage, s.userRepo)
		if err != nil {
			s.reply(msg, bot, fmt.Sprintf("Ошибка: %s", err.Error()))
			return
		}

		if err := s.ensureRestrictAllowed(chatID, msg.From.ID, targetID); err != nil {
			s.reply(msg, bot, fmt.Sprintf("Ошибка: %s.", err.Error()))
			return
		}

		if err := s.executor.Ban(chatID, targetID); err != nil {
			s.reply(msg, bot, fmt.Sprintf("Ошибка бана: %s", err.Error()))
			return
		}

		s.logAction(msg.From.ID, targetID, chatID, "ban", reason, 0)
		s.reply(msg, bot, "Пользователь забанен.")
	}
}

func (s *Service) KickHandler(bot *tgbotapi.BotAPI) func(tgbotapi.Message) {
	return func(msg tgbotapi.Message) {
		if msg.From == nil {
			return
		}
		chatID := msg.Chat.ID

		if err := s.ensureAdmin(chatID, msg.From.ID); err != nil {
			s.reply(msg, bot, fmt.Sprintf("Ошибка: %s.", err.Error()))
			return
		}

		targetID, err := parseKickArgs(msg.CommandArguments(), msg.ReplyToMessage, s.userRepo)
		if err != nil {
			s.reply(msg, bot, fmt.Sprintf("Ошибка: %s", err.Error()))
			return
		}

		if err := s.ensureRestrictAllowed(chatID, msg.From.ID, targetID); err != nil {
			s.reply(msg, bot, fmt.Sprintf("Ошибка: %s.", err.Error()))
			return
		}

		if err := s.executor.Kick(chatID, targetID); err != nil {
			s.reply(msg, bot, fmt.Sprintf("Ошибка кика: %s", err.Error()))
			return
		}

		s.logAction(msg.From.ID, targetID, chatID, "kick", "", 0)
		s.reply(msg, bot, "Пользователь кикнут.")
	}
}

func (s *Service) UnmuteHandler(bot *tgbotapi.BotAPI) func(tgbotapi.Message) {
	return func(msg tgbotapi.Message) {
		if msg.From == nil {
			return
		}
		chatID := msg.Chat.ID

		if err := s.ensureAdmin(chatID, msg.From.ID); err != nil {
			s.reply(msg, bot, fmt.Sprintf("Ошибка: %s.", err.Error()))
			return
		}

		targetID, err := parseUnmuteArgs(msg.CommandArguments(), msg.ReplyToMessage, s.userRepo)
		if err != nil {
			s.reply(msg, bot, fmt.Sprintf("Ошибка: %s", err.Error()))
			return
		}

		if err := s.ensureRestrictAllowed(chatID, msg.From.ID, targetID); err != nil {
			s.reply(msg, bot, fmt.Sprintf("Ошибка: %s.", err.Error()))
			return
		}

		if err := s.executor.Unmute(chatID, targetID); err != nil {
			s.reply(msg, bot, fmt.Sprintf("Ошибка размута: %s", err.Error()))
			return
		}

		s.logAction(msg.From.ID, targetID, chatID, "unmute", "", 0)
		s.reply(msg, bot, "Пользователь размучен.")
	}
}

func (s *Service) BalabolHandler(bot *tgbotapi.BotAPI) func(tgbotapi.Message) {
	return func(msg tgbotapi.Message) {
		if msg.From == nil {
			return
		}
		chatID := msg.Chat.ID

		if err := s.ensureAdmin(chatID, msg.From.ID); err != nil {
			s.reply(msg, bot, fmt.Sprintf("Ошибка: %s.", err.Error()))
			return
		}

		targetID, reason, err := parseBalabolArgs(msg.CommandArguments(), msg.ReplyToMessage, s.userRepo)
		if err != nil {
			s.reply(msg, bot, fmt.Sprintf("Ошибка: %s", err.Error()))
			return
		}
		if targetID == s.api.Self.ID {
			s.reply(msg, bot, "Боту метка балабола не нужна, он и так всё помнит.")
			return
		}

		if err := s.markRepo.SetMark(context.Background(), msg.Chat.ID, targetID, msg.From.ID, repository.MarkBalabol, reason); err != nil {
			s.reply(msg, bot, fmt.Sprintf("Ошибка метки: %s", err.Error()))
			return
		}
		s.reply(msg, bot, "Метка балабола выдана. Репорты от пользователя теперь не принимаются.")
	}
}

func (s *Service) UnbalabolHandler(bot *tgbotapi.BotAPI) func(tgbotapi.Message) {
	return func(msg tgbotapi.Message) {
		if msg.From == nil {
			return
		}
		chatID := msg.Chat.ID

		if err := s.ensureAdmin(chatID, msg.From.ID); err != nil {
			s.reply(msg, bot, fmt.Sprintf("Ошибка: %s.", err.Error()))
			return
		}

		targetID, err := parseUnmuteArgs(msg.CommandArguments(), msg.ReplyToMessage, s.userRepo)
		if err != nil {
			s.reply(msg, bot, fmt.Sprintf("Ошибка: %s", err.Error()))
			return
		}
		if err := s.markRepo.RemoveMark(context.Background(), msg.Chat.ID, targetID, repository.MarkBalabol); err != nil {
			s.reply(msg, bot, fmt.Sprintf("Ошибка снятия метки: %s", err.Error()))
			return
		}
		s.reply(msg, bot, "Метка балабола снята. Репорты от пользователя снова принимаются.")
	}
}

func (s *Service) reply(msg tgbotapi.Message, bot *tgbotapi.BotAPI, text string) {
	msgText := tgbotapi.NewMessage(msg.Chat.ID, text)
	msgText.ReplyToMessageID = msg.MessageID
	if _, err := bot.Send(msgText); err != nil {
		s.logger.Warn("send moderation reply failed", zap.Error(err), zap.Int64("chat_id", msg.Chat.ID))
	}
}

func formatDuration(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	mins := int(d.Minutes()) % 60
	parts := []string{}
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%dd", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
	}
	if mins > 0 {
		parts = append(parts, fmt.Sprintf("%dm", mins))
	}
	if len(parts) == 0 {
		return "0m"
	}
	return strings.Join(parts, " ")
}
