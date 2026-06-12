package moderation

import (
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestParseMuteArgsReplyFirst(t *testing.T) {
	reply := &tgbotapi.Message{From: &tgbotapi.User{ID: 42}}

	targetID, duration, reason, err := parseMuteArgs("2h шумит в чате", reply, nil)
	if err != nil {
		t.Fatalf("parseMuteArgs() error = %v", err)
	}
	if targetID != 42 {
		t.Fatalf("targetID = %d, want 42", targetID)
	}
	if duration != 2*time.Hour {
		t.Fatalf("duration = %v, want 2h", duration)
	}
	if reason != "шумит в чате" {
		t.Fatalf("reason = %q", reason)
	}
}

func TestParseBanArgsReplyKeepsReason(t *testing.T) {
	reply := &tgbotapi.Message{From: &tgbotapi.User{ID: 77}}

	targetID, reason, err := parseBanArgs("очень плохое поведение", reply, nil)
	if err != nil {
		t.Fatalf("parseBanArgs() error = %v", err)
	}
	if targetID != 77 {
		t.Fatalf("targetID = %d, want 77", targetID)
	}
	if reason != "очень плохое поведение" {
		t.Fatalf("reason = %q", reason)
	}
}
