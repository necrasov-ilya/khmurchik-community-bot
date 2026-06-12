package moderation

import (
	"context"
	"errors"
	"time"

	"github.com/evart2006/khmurchik-community-bot/internal/middleware"
	"github.com/evart2006/khmurchik-community-bot/internal/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Service struct {
	api      *tgbotapi.BotAPI
	db       *pgxpool.Pool
	logger   *zap.Logger
	executor *Executor
	repo     *Repository
	userRepo *repository.UserRepository
	markRepo *repository.UserMarkRepository
	checker  *middleware.AdminChecker
}

func NewService(api *tgbotapi.BotAPI, db *pgxpool.Pool, logger *zap.Logger, checker *middleware.AdminChecker) *Service {
	return &Service{
		api:      api,
		db:       db,
		logger:   logger,
		executor: NewExecutor(api),
		repo:     NewRepository(db, logger),
		userRepo: repository.NewUserRepository(db),
		markRepo: repository.NewUserMarkRepository(db),
		checker:  checker,
	}
}

func (s *Service) logAction(adminID, targetID, chatID int64, action, reason string, duration time.Duration) {
	var expiresAt *time.Time
	if action == "mute" {
		exp := time.Now().Add(duration)
		expiresAt = &exp
	}
	log := ModerationLog{
		UserID:     targetID,
		AdminID:    adminID,
		ChatID:     chatID,
		ActionType: action,
		Reason:     reason,
		ExpiresAt:  expiresAt,
	}
	go func() {
		if err := s.repo.LogAction(context.Background(), log); err != nil {
			s.logger.Error("failed to log moderation action", zap.Error(err))
		}
	}()
}

func (s *Service) ensureRestrictAllowed(chatID, adminID, targetID int64) error {
	if targetID == s.api.Self.ID {
		return errors.New("бот не может модерировать сам себя")
	}

	adminOK, err := s.checker.CheckCanRestrict(chatID, adminID)
	if err != nil {
		return err
	}
	if !adminOK {
		return errors.New("у вас нет прав ограничивать пользователей")
	}

	botOK, err := s.checker.CheckCanRestrict(chatID, s.api.Self.ID)
	if err != nil {
		return err
	}
	if !botOK {
		return errors.New("у бота нет права ограничивать пользователей")
	}

	member, err := s.checker.GetMember(chatID, targetID)
	if err != nil {
		return err
	}
	if middleware.IsAdmin(member) {
		return errors.New("нельзя модерировать владельца или администратора")
	}

	return nil
}

func (s *Service) ensureAdmin(chatID, adminID int64) error {
	ok, err := s.checker.Check(chatID, adminID)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("у вас нет прав для этой команды")
	}
	return nil
}
