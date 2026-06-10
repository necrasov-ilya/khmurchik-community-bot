package moderation

import (
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *Service) MuteHandler(bot *tgbotapi.BotAPI) func(tgbotapi.Message) {
	return func(msg tgbotapi.Message) {
		ok, err := s.checker.Check(msg.From.ID)
		if err != nil || !ok {
			s.reply(msg, bot, "У вас нет прав для этой команды.")
			return
		}

		targetID, duration, reason, err := parseMuteArgs(msg.CommandArguments(), msg.ReplyToMessage, s.userRepo)
		if err != nil {
			s.reply(msg, bot, fmt.Sprintf("Ошибка: %s", err.Error()))
			return
		}

		if err := s.executor.MuteUntil(s.targetChatID, targetID, duration); err != nil {
			s.reply(msg, bot, fmt.Sprintf("Ошибка мута: %s", err.Error()))
			return
		}

		s.logAction(msg.From.ID, targetID, s.targetChatID, "mute", reason, duration)
		s.reply(msg, bot, fmt.Sprintf("Пользователь замучен на %s", formatDuration(duration)))
	}
}

func (s *Service) BanHandler(bot *tgbotapi.BotAPI) func(tgbotapi.Message) {
	return func(msg tgbotapi.Message) {
		ok, err := s.checker.Check(msg.From.ID)
		if err != nil || !ok {
			s.reply(msg, bot, "У вас нет прав для этой команды.")
			return
		}

		targetID, reason, err := parseBanArgs(msg.CommandArguments(), msg.ReplyToMessage, s.userRepo)
		if err != nil {
			s.reply(msg, bot, fmt.Sprintf("Ошибка: %s", err.Error()))
			return
		}

		if err := s.executor.Ban(s.targetChatID, targetID); err != nil {
			s.reply(msg, bot, fmt.Sprintf("Ошибка бана: %s", err.Error()))
			return
		}

		s.logAction(msg.From.ID, targetID, s.targetChatID, "ban", reason, 0)
		s.reply(msg, bot, "Пользователь забанен.")
	}
}

func (s *Service) KickHandler(bot *tgbotapi.BotAPI) func(tgbotapi.Message) {
	return func(msg tgbotapi.Message) {
		ok, err := s.checker.Check(msg.From.ID)
		if err != nil || !ok {
			s.reply(msg, bot, "У вас нет прав для этой команды.")
			return
		}

		targetID, err := parseKickArgs(msg.CommandArguments(), msg.ReplyToMessage, s.userRepo)
		if err != nil {
			s.reply(msg, bot, fmt.Sprintf("Ошибка: %s", err.Error()))
			return
		}

		if err := s.executor.Kick(s.targetChatID, targetID); err != nil {
			s.reply(msg, bot, fmt.Sprintf("Ошибка кика: %s", err.Error()))
			return
		}

		s.logAction(msg.From.ID, targetID, s.targetChatID, "kick", "", 0)
		s.reply(msg, bot, "Пользователь кикнут.")
	}
}

func (s *Service) UnmuteHandler(bot *tgbotapi.BotAPI) func(tgbotapi.Message) {
	return func(msg tgbotapi.Message) {
		ok, err := s.checker.Check(msg.From.ID)
		if err != nil || !ok {
			s.reply(msg, bot, "У вас нет прав для этой команды.")
			return
		}

		targetID, err := parseUnmuteArgs(msg.CommandArguments(), msg.ReplyToMessage, s.userRepo)
		if err != nil {
			s.reply(msg, bot, fmt.Sprintf("Ошибка: %s", err.Error()))
			return
		}

		if err := s.executor.Unmute(s.targetChatID, targetID); err != nil {
			s.reply(msg, bot, fmt.Sprintf("Ошибка размута: %s", err.Error()))
			return
		}

		s.logAction(msg.From.ID, targetID, s.targetChatID, "unmute", "", 0)
		s.reply(msg, bot, "Пользователь размучен.")
	}
}

func (s *Service) reply(msg tgbotapi.Message, bot *tgbotapi.BotAPI, text string) {
	msgText := tgbotapi.NewMessage(msg.Chat.ID, text)
	msgText.ReplyToMessageID = msg.MessageID
	bot.Send(msgText)
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
