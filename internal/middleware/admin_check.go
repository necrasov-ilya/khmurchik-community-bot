package middleware

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type AdminChecker struct {
	bot    *tgbotapi.BotAPI
	chatID int64
	logger *zap.Logger
}

func NewAdminChecker(bot *tgbotapi.BotAPI, chatID int64, logger *zap.Logger) *AdminChecker {
	return &AdminChecker{bot: bot, chatID: chatID, logger: logger}
}

func (a *AdminChecker) Check(userID int64) (bool, error) {
	member, err := a.bot.GetChatMember(tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: a.chatID,
			UserID: userID,
		},
	})
	if err != nil {
		return false, err
	}
	return member.Status == "creator" || member.Status == "administrator", nil
}
