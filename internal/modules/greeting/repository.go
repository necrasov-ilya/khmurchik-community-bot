package greeting

import (
	"go.uber.org/zap"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func sendDailyGreeting(bot *tgbotapi.BotAPI, targetChatID int64, logger *zap.Logger) {
	msg := tgbotapi.NewMessage(targetChatID, DefaultMessage)
	if _, err := bot.Send(msg); err != nil {
		logger.Error("failed to send greeting", zap.Error(err), zap.Int64("chat_id", targetChatID))
	}
}
