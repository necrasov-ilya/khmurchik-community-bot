package app

import (
	"context"

	"github.com/evart2006/khmurchik-community-bot/internal/bot"
	"github.com/evart2006/khmurchik-community-bot/internal/config"
	"github.com/evart2006/khmurchik-community-bot/internal/handlers"
	"github.com/evart2006/khmurchik-community-bot/internal/modules/greeting"
	"github.com/evart2006/khmurchik-community-bot/internal/modules/moderation"
	"github.com/evart2006/khmurchik-community-bot/internal/modules/reports"
	"github.com/evart2006/khmurchik-community-bot/internal/modulesystem"
	"github.com/evart2006/khmurchik-community-bot/internal/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type App struct {
	cfg    *config.Config
	db     *pgxpool.Pool
	logger *zap.Logger
	bot    *tgbotapi.BotAPI
	router *bot.Router
	reg    *modulesystem.Registry
}

func New(cfg *config.Config, db *pgxpool.Pool, logger *zap.Logger) *App {
	return &App{cfg: cfg, db: db, logger: logger}
}

func (a *App) Start(ctx context.Context) error {
	b, err := tgbotapi.NewBotAPI(a.cfg.Bot.Token)
	if err != nil {
		return err
	}
	a.bot = b

	a.router = bot.NewRouter()
	a.router.Register("start", handlers.StartHandler(a.bot))
	a.router.Register("help", handlers.HelpHandler(a.bot))
	a.router.SetUnknownHandler(handlers.UnknownHandler(a.bot))

	a.reg = modulesystem.NewRegistry()

	greetingMod := greeting.NewModule(a.bot, a.db, a.cfg.Defaults, a.logger)
	a.reg.Register(greetingMod)

	moderationMod := moderation.NewModule(a.bot, a.db, a.logger)
	a.reg.Register(moderationMod)

	reportsMod := reports.NewModule(a.bot, a.db, a.logger)
	a.reg.Register(reportsMod)

	if err := a.reg.RegisterAll(a.bot, a.router); err != nil {
		return err
	}

	if err := a.registerBotCommands(); err != nil {
		a.logger.Warn("failed to register bot commands", zap.Error(err))
	}

	if err := a.reg.StartAll(ctx); err != nil {
		return err
	}

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = int(a.cfg.Server.PollTimeout)
	updateConfig.AllowedUpdates = a.cfg.Server.AllowedUpdates
	updates := a.bot.GetUpdatesChan(updateConfig)

	userRepo := repository.NewUserRepository(a.db)
	chatRepo := repository.NewChatRepository(a.db, a.cfg.Defaults)

	go func() {
		for update := range updates {
			if update.Message != nil && update.Message.From != nil {
				if err := userRepo.UpsertUser(update.Message.From); err != nil {
					a.logger.Error("upsert user failed", zap.Error(err), zap.Int64("user_id", update.Message.From.ID))
				}
				if err := chatRepo.UpsertChat(*update.Message.Chat); err != nil {
					a.logger.Error("upsert chat failed", zap.Error(err), zap.Int64("chat_id", update.Message.Chat.ID))
				}
				if update.Message.ReplyToMessage != nil && update.Message.ReplyToMessage.From != nil {
					if err := userRepo.UpsertUser(update.Message.ReplyToMessage.From); err != nil {
						a.logger.Error("upsert reply user failed", zap.Error(err), zap.Int64("user_id", update.Message.ReplyToMessage.From.ID))
					}
				}
			}
			if update.MyChatMember != nil {
				if err := chatRepo.UpsertChat(update.MyChatMember.Chat); err != nil {
					a.logger.Error("upsert my_chat_member chat failed", zap.Error(err), zap.Int64("chat_id", update.MyChatMember.Chat.ID))
				}
				a.handleMyChatMember(update.MyChatMember)
			}
			a.router.Handle(update)
		}
	}()

	<-ctx.Done()

	a.reg.StopAll(context.Background())
	a.bot.StopReceivingUpdates()

	return nil
}

func (a *App) registerBotCommands() error {
	commands := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{Command: "help", Description: "Показать команды"},
		tgbotapi.BotCommand{Command: "report", Description: "Пожаловаться админам на сообщение"},
		tgbotapi.BotCommand{Command: "mute", Description: "Замутить пользователя"},
		tgbotapi.BotCommand{Command: "ban", Description: "Забанить пользователя"},
		tgbotapi.BotCommand{Command: "kick", Description: "Кикнуть пользователя"},
		tgbotapi.BotCommand{Command: "unmute", Description: "Снять мут"},
		tgbotapi.BotCommand{Command: "balabol", Description: "Выдать метку балабола"},
		tgbotapi.BotCommand{Command: "unbalabol", Description: "Снять метку балабола"},
		tgbotapi.BotCommand{Command: "greeting_on", Description: "Включить утреннее приветствие"},
		tgbotapi.BotCommand{Command: "greeting_off", Description: "Выключить утреннее приветствие"},
		tgbotapi.BotCommand{Command: "greeting_time", Description: "Настроить время приветствия"},
		tgbotapi.BotCommand{Command: "greeting_text", Description: "Настроить текст приветствия"},
		tgbotapi.BotCommand{Command: "greeting_status", Description: "Показать настройки приветствия"},
	)
	_, err := a.bot.Request(commands)
	return err
}

func (a *App) handleMyChatMember(update *tgbotapi.ChatMemberUpdated) {
	if update.Chat.Type != "group" && update.Chat.Type != "supergroup" {
		return
	}
	if update.NewChatMember.User == nil || update.NewChatMember.User.ID != a.bot.Self.ID {
		return
	}
	if update.NewChatMember.Status != "member" && update.NewChatMember.Status != "administrator" {
		return
	}

	msg := tgbotapi.NewMessage(update.Chat.ID, "Я в чате. Для модерации нужны права администратора с возможностью ограничивать пользователей.\n\nКоманды: /help\nПриветствие: /greeting_status, /greeting_time 10:00 Europe/Minsk, /greeting_text текст")
	if _, err := a.bot.Send(msg); err != nil {
		a.logger.Warn("failed to send onboarding message", zap.Error(err), zap.Int64("chat_id", update.Chat.ID))
	}
}
