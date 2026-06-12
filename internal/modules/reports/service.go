package reports

import (
	"context"
	"fmt"
	"html"
	"strings"

	"github.com/evart2006/khmurchik-community-bot/internal/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Service struct {
	api        *tgbotapi.BotAPI
	logger     *zap.Logger
	reportRepo *repository.ReportRepository
	markRepo   *repository.UserMarkRepository
}

func NewService(api *tgbotapi.BotAPI, db *pgxpool.Pool, logger *zap.Logger) *Service {
	return &Service{
		api:        api,
		logger:     logger,
		reportRepo: repository.NewReportRepository(db),
		markRepo:   repository.NewUserMarkRepository(db),
	}
}

func (s *Service) ReportHandler(bot *tgbotapi.BotAPI) func(tgbotapi.Message) {
	return func(msg tgbotapi.Message) {
		if msg.From == nil {
			return
		}
		if msg.Chat.Type != "group" && msg.Chat.Type != "supergroup" {
			s.reply(msg, "Репорты работают только в группах.")
			return
		}
		if msg.ReplyToMessage == nil || msg.ReplyToMessage.From == nil {
			s.reply(msg, "Ответьте командой /report на сообщение, которое нужно показать админам.")
			return
		}

		blocked, err := s.markRepo.HasMark(context.Background(), msg.Chat.ID, msg.From.ID, repository.MarkBalabol)
		if err != nil {
			s.logger.Error("failed to check balabol mark", zap.Error(err), zap.Int64("chat_id", msg.Chat.ID), zap.Int64("user_id", msg.From.ID))
			s.reply(msg, "Не смог проверить метки. Попробуйте позже.")
			return
		}
		if blocked {
			s.reply(msg, "Вы не можете отправить репорт, так как у вас метка балабола.")
			return
		}

		reason := strings.TrimSpace(msg.CommandArguments())
		reportID, err := s.reportRepo.Create(context.Background(), repository.Report{
			ChatID:            msg.Chat.ID,
			ReporterID:        msg.From.ID,
			ReportedUserID:    msg.ReplyToMessage.From.ID,
			ReportedMessageID: msg.ReplyToMessage.MessageID,
			Reason:            reason,
		})
		if err != nil {
			s.logger.Error("failed to create report", zap.Error(err), zap.Int64("chat_id", msg.Chat.ID))
			s.reply(msg, "Не смог сохранить репорт. Попробуйте позже.")
			return
		}

		admins, err := s.api.GetChatAdministrators(tgbotapi.ChatAdministratorsConfig{
			ChatConfig: tgbotapi.ChatConfig{ChatID: msg.Chat.ID},
		})
		if err != nil {
			s.logger.Error("failed to get chat administrators", zap.Error(err), zap.Int64("chat_id", msg.Chat.ID))
			s.reply(msg, "Репорт сохранён, но я не смог получить список админов.")
			return
		}

		reportText := BuildReportMessage(reportID, msg, admins, reason)
		adminMsg := tgbotapi.NewMessage(msg.Chat.ID, reportText)
		adminMsg.ParseMode = tgbotapi.ModeHTML
		adminMsg.DisableWebPagePreview = true
		adminMsg.ReplyToMessageID = msg.ReplyToMessage.MessageID
		if _, err := bot.Send(adminMsg); err != nil {
			s.logger.Error("failed to send report to chat", zap.Error(err), zap.Int64("chat_id", msg.Chat.ID))
			s.reply(msg, "Репорт сохранён, но я не смог отправить сообщение админам.")
			return
		}

		s.reply(msg, "Репорт отправлен админам, ожидайте.")
	}
}

func BuildReportMessage(reportID int64, msg tgbotapi.Message, admins []tgbotapi.ChatMember, reason string) string {
	mentions := BuildAdminMentions(admins)
	if mentions == "" {
		mentions = "админы"
	}

	reporter := UserLabel(msg.From)
	reported := "неизвестный пользователь"
	if msg.ReplyToMessage != nil && msg.ReplyToMessage.From != nil {
		reported = UserLabel(msg.ReplyToMessage.From)
	}
	if reason == "" {
		reason = "не указана"
	}

	return fmt.Sprintf(
		"🚨 Репорт #%d для %s\n\nОт: %s\nНа: %s\nПричина: %s",
		reportID,
		mentions,
		html.EscapeString(reporter),
		html.EscapeString(reported),
		html.EscapeString(reason),
	)
}

func BuildAdminMentions(admins []tgbotapi.ChatMember) string {
	mentions := make([]string, 0, len(admins))
	for _, admin := range admins {
		if admin.User == nil || admin.User.IsBot {
			continue
		}
		if admin.User.UserName != "" {
			mentions = append(mentions, "@"+html.EscapeString(admin.User.UserName))
			continue
		}

		name := strings.TrimSpace(strings.Join([]string{admin.User.FirstName, admin.User.LastName}, " "))
		if name == "" {
			name = fmt.Sprintf("admin-%d", admin.User.ID)
		}
		mentions = append(mentions, fmt.Sprintf(`<a href="tg://user?id=%d">%s</a>`, admin.User.ID, html.EscapeString(name)))
	}
	return strings.Join(mentions, " ")
}

func UserLabel(user *tgbotapi.User) string {
	if user == nil {
		return "unknown"
	}
	if user.UserName != "" {
		return "@" + user.UserName
	}
	name := strings.TrimSpace(strings.Join([]string{user.FirstName, user.LastName}, " "))
	if name == "" {
		return fmt.Sprintf("id:%d", user.ID)
	}
	return fmt.Sprintf("%s (id:%d)", name, user.ID)
}

func (s *Service) reply(msg tgbotapi.Message, text string) {
	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	reply.ReplyToMessageID = msg.MessageID
	if _, err := s.api.Send(reply); err != nil {
		s.logger.Warn("send report reply failed", zap.Error(err), zap.Int64("chat_id", msg.Chat.ID))
	}
}
