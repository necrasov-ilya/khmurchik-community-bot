package moderation

import (
	"fmt"
	"strings"
	"time"

	"github.com/evart2006/khmurchik-community-bot/internal/repository"
	"github.com/evart2006/khmurchik-community-bot/internal/timeutil"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func parseTarget(args string, replyTo *tgbotapi.Message, userRepo *repository.UserRepository) (int64, error) {
	if replyTo != nil && replyTo.From != nil {
		return replyTo.From.ID, nil
	}

	args = strings.TrimSpace(args)
	if strings.HasPrefix(args, "@") {
		username := strings.ToLower(args)
		id, err := userRepo.FindByUsername(username)
		if err != nil {
			return 0, fmt.Errorf("user %s not found in database", args)
		}
		return id, nil
	}

	return 0, fmt.Errorf("no target: reply to message or provide @username")
}

func parseMuteArgs(args string, replyTo *tgbotapi.Message, userRepo *repository.UserRepository) (targetID int64, duration time.Duration, reason string, err error) {
	args = strings.TrimSpace(args)

	if args != "" && strings.HasPrefix(args, "@") {
		parts := strings.Fields(args)
		if len(parts) < 2 {
			return 0, 0, "", fmt.Errorf("/mute requires @username and duration: /mute @user 1d \"reason\"")
		}
		username := strings.ToLower(parts[0])
		id, err := userRepo.FindByUsername(username)
		if err != nil {
			return 0, 0, "", fmt.Errorf("user %s not found in database", parts[0])
		}
		duration, err = timeutil.ParseDuration(parts[1])
		if err != nil {
			return 0, 0, "", err
		}
		if len(parts) > 2 {
			reason = strings.Join(parts[2:], " ")
		}
		return id, duration, reason, nil
	}

	if replyTo == nil || replyTo.From == nil {
		return 0, 0, "", fmt.Errorf("/mute requires reply to a message or @username")
	}
	parts := strings.SplitN(args, " ", 2)
	if len(parts) < 1 || args == "" {
		return 0, 0, "", fmt.Errorf("/mute requires duration: /mute 1d \"reason\"")
	}
	duration, err = timeutil.ParseDuration(parts[0])
	if err != nil {
		return 0, 0, "", err
	}
	if len(parts) > 1 {
		reason = parts[1]
	}
	return replyTo.From.ID, duration, reason, nil
}

func parseBanArgs(args string, replyTo *tgbotapi.Message, userRepo *repository.UserRepository) (int64, string, error) {
	args = strings.TrimSpace(args)
	targetID, err := parseTarget(args, replyTo, userRepo)
	if err != nil {
		return 0, "", err
	}
	reason := strings.TrimSpace(args)
	if strings.HasPrefix(reason, "@") {
		parts := strings.Fields(reason)
		if len(parts) > 1 {
			reason = strings.Join(parts[1:], " ")
		} else {
			reason = ""
		}
	}
	return targetID, reason, nil
}

func parseKickArgs(args string, replyTo *tgbotapi.Message, userRepo *repository.UserRepository) (int64, error) {
	return parseTarget(args, replyTo, userRepo)
}

func parseUnmuteArgs(args string, replyTo *tgbotapi.Message, userRepo *repository.UserRepository) (int64, error) {
	return parseTarget(args, replyTo, userRepo)
}
