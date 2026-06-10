package moderation

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Executor struct {
	bot *tgbotapi.BotAPI
}

func NewExecutor(bot *tgbotapi.BotAPI) *Executor {
	return &Executor{bot: bot}
}

func (e *Executor) MuteUntil(chatID, userID int64, duration time.Duration) error {
	permissions := &tgbotapi.ChatPermissions{
		CanSendMessages:      false,
		CanSendPolls:         false,
		CanSendOtherMessages: false,
	}
	req := tgbotapi.RestrictChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatID,
			UserID: userID,
		},
		UntilDate: time.Now().Add(duration).Unix(),
		Permissions: permissions,
	}
	_, err := e.bot.Request(req)
	return err
}

func (e *Executor) Ban(chatID, userID int64) error {
	req := tgbotapi.BanChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatID,
			UserID: userID,
		},
		UntilDate: 0,
	}
	_, err := e.bot.Request(req)
	return err
}

func (e *Executor) Kick(chatID, userID int64) error {
	banReq := tgbotapi.BanChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatID,
			UserID: userID,
		},
		UntilDate: 0,
	}
	if _, err := e.bot.Request(banReq); err != nil {
		return err
	}
	unbanReq := tgbotapi.UnbanChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatID,
			UserID: userID,
		},
		OnlyIfBanned: true,
	}
	_, err := e.bot.Request(unbanReq)
	return err
}

func (e *Executor) Unmute(chatID, userID int64) error {
	permissions := &tgbotapi.ChatPermissions{
		CanSendMessages:      true,
		CanSendPolls:         true,
		CanSendOtherMessages: true,
		CanAddWebPagePreviews: true,
		CanInviteUsers:       true,
		CanPinMessages:       true,
	}
	req := tgbotapi.RestrictChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatID,
			UserID: userID,
		},
		UntilDate:   0,
		Permissions: permissions,
	}
	_, err := e.bot.Request(req)
	return err
}
