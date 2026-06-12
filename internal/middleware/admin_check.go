package middleware

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type AdminChecker struct {
	bot    *tgbotapi.BotAPI
	logger *zap.Logger
}

func NewAdminChecker(bot *tgbotapi.BotAPI, logger *zap.Logger) *AdminChecker {
	return &AdminChecker{bot: bot, logger: logger}
}

func (a *AdminChecker) Check(chatID, userID int64) (bool, error) {
	member, err := a.GetMember(chatID, userID)
	if err != nil {
		return false, err
	}
	return IsAdmin(member), nil
}

func (a *AdminChecker) CheckCanRestrict(chatID, userID int64) (bool, error) {
	member, err := a.GetMember(chatID, userID)
	if err != nil {
		return false, err
	}
	return CanRestrict(member), nil
}

func (a *AdminChecker) GetMember(chatID, userID int64) (tgbotapi.ChatMember, error) {
	member, err := a.bot.GetChatMember(tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: chatID,
			UserID: userID,
		},
	})
	if err != nil {
		a.logger.Warn("get chat member failed", zap.Error(err), zap.Int64("chat_id", chatID), zap.Int64("user_id", userID))
		return tgbotapi.ChatMember{}, err
	}
	return member, nil
}

func IsAdmin(member tgbotapi.ChatMember) bool {
	return member.Status == "creator" || member.Status == "administrator"
}

func CanRestrict(member tgbotapi.ChatMember) bool {
	if member.Status == "creator" {
		return true
	}
	return member.Status == "administrator" && member.CanRestrictMembers
}
