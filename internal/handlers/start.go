package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartHandler(bot *tgbotapi.BotAPI) func(tgbotapi.Message) {
	return func(msg tgbotapi.Message) {
		_, _ = bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Привет! Я бот модерации сообщества. Используй /help для списка команд."))
	}
}
