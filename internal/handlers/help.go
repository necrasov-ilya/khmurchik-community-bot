package handlers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const helpText = `Команды пользователей:
/report причина — ответом на сообщение отправить репорт админам.

Команды админов:
/mute 1h причина — мут ответом на сообщение.
/ban причина — бан ответом на сообщение.
/kick — кик ответом на сообщение.
/unmute — снять мут ответом на сообщение.
/balabol причина — выдать метку балабола.
/unbalabol — снять метку балабола.

Приветствие:
/greeting_status
/greeting_on
/greeting_off
/greeting_time 10:00 Europe/Minsk
/greeting_text текст сообщения`

func HelpHandler(bot *tgbotapi.BotAPI) func(tgbotapi.Message) {
	return func(msg tgbotapi.Message) {
		_, _ = bot.Send(tgbotapi.NewMessage(msg.Chat.ID, helpText))
	}
}
