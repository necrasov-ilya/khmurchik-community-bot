package modulesystem

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/evart2006/khmurchik-community-bot/internal/bot"
)

type Registry struct {
	modules []Module
}

func NewRegistry() *Registry {
	return &Registry{modules: make([]Module, 0)}
}

func (r *Registry) Register(m Module) {
	r.modules = append(r.modules, m)
}

func (r *Registry) StartAll(ctx context.Context) error {
	for _, m := range r.modules {
		if err := m.OnStart(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (r *Registry) StopAll(ctx context.Context) error {
	for i := len(r.modules) - 1; i >= 0; i-- {
		if err := r.modules[i].OnStop(ctx); err != nil {
		}
	}
	return nil
}

func (r *Registry) Modules() []Module {
	return r.modules
}

func (r *Registry) RegisterAll(bot *tgbotapi.BotAPI, router *bot.Router) error {
	for _, m := range r.modules {
		if err := m.Register(bot, router); err != nil {
			return err
		}
	}
	return nil
}
