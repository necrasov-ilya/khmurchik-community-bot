package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Router struct {
	cmdHandlers    map[string]func(tgbotapi.Message)
	unknownHandler func(tgbotapi.Message)
}

func NewRouter() *Router {
	return &Router{
		cmdHandlers: make(map[string]func(tgbotapi.Message)),
	}
}

func (r *Router) Register(cmd string, handler func(tgbotapi.Message)) {
	r.cmdHandlers[cmd] = handler
}

func (r *Router) SetUnknownHandler(handler func(tgbotapi.Message)) {
	r.unknownHandler = handler
}

func (r *Router) Handle(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}
	msg := *update.Message

	cmd := msg.Command()
	if cmd != "" {
		if handler, ok := r.cmdHandlers[cmd]; ok {
			handler(msg)
			return
		}
		if r.unknownHandler != nil {
			r.unknownHandler(msg)
		}
		return
	}
}
