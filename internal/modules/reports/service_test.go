package reports

import (
	"strings"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestBuildAdminMentions(t *testing.T) {
	admins := []tgbotapi.ChatMember{
		{User: &tgbotapi.User{ID: 1, UserName: "admin_one", FirstName: "Admin"}},
		{User: &tgbotapi.User{ID: 2, FirstName: "No", LastName: "Username"}},
		{User: &tgbotapi.User{ID: 3, UserName: "bot_admin", IsBot: true}},
	}

	got := BuildAdminMentions(admins)
	if !strings.Contains(got, "@admin_one") {
		t.Fatalf("mentions = %q, want username mention", got)
	}
	if !strings.Contains(got, `tg://user?id=2`) {
		t.Fatalf("mentions = %q, want tg user link", got)
	}
	if strings.Contains(got, "bot_admin") {
		t.Fatalf("mentions = %q, bot admins should be skipped", got)
	}
}

func TestBuildReportMessageEscapesUserInput(t *testing.T) {
	msg := tgbotapi.Message{
		From: &tgbotapi.User{ID: 10, UserName: "reporter"},
		ReplyToMessage: &tgbotapi.Message{
			From:      &tgbotapi.User{ID: 11, FirstName: "Bad <User>"},
			MessageID: 99,
		},
	}
	admins := []tgbotapi.ChatMember{{User: &tgbotapi.User{ID: 1, UserName: "admin"}}}

	got := BuildReportMessage(5, msg, admins, "<script>")
	if !strings.Contains(got, "Репорт #5") {
		t.Fatalf("message = %q, want report id", got)
	}
	if strings.Contains(got, "<script>") || strings.Contains(got, "Bad <User>") {
		t.Fatalf("message = %q, want escaped input", got)
	}
}
