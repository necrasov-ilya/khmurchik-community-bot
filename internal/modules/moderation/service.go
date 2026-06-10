package moderation

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/evart2006/khmurchik-community-bot/internal/middleware"
	"github.com/evart2006/khmurchik-community-bot/internal/repository"
	"go.uber.org/zap"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Service struct {
	api          *tgbotapi.BotAPI
	db           *pgxpool.Pool
	logger       *zap.Logger
	executor     *Executor
	repo         *Repository
	userRepo     *repository.UserRepository
	checker      *middleware.AdminChecker
	targetChatID int64
}

func NewService(api *tgbotapi.BotAPI, db *pgxpool.Pool, logger *zap.Logger, checker *middleware.AdminChecker, targetChatID int64) *Service {
	return &Service{
		api:          api,
		db:           db,
		logger:       logger,
		executor:     NewExecutor(api),
		repo:         NewRepository(db, logger),
		userRepo:     repository.NewUserRepository(db),
		checker:      checker,
		targetChatID: targetChatID,
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
