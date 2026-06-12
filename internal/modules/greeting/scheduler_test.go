package greeting

import (
	"testing"
	"time"

	"github.com/evart2006/khmurchik-community-bot/internal/repository"
)

func TestShouldSendGreetingAtConfiguredLocalTime(t *testing.T) {
	now := time.Date(2026, 6, 12, 7, 0, 0, 0, time.UTC) // 10:00 Europe/Minsk.
	chat := repository.ChatSettings{
		GreetingTime:     "10:00",
		GreetingTimezone: "Europe/Minsk",
	}

	should, _, err := ShouldSendGreeting(chat, now)
	if err != nil {
		t.Fatalf("ShouldSendGreeting() error = %v", err)
	}
	if !should {
		t.Fatal("ShouldSendGreeting() = false, want true")
	}
}

func TestShouldSendGreetingSkipsSameLocalDate(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Minsk")
	if err != nil {
		t.Fatal(err)
	}
	sent := time.Date(2026, 6, 12, 10, 0, 0, 0, loc)
	now := time.Date(2026, 6, 12, 7, 0, 0, 0, time.UTC)
	chat := repository.ChatSettings{
		GreetingTime:         "10:00",
		GreetingTimezone:     "Europe/Minsk",
		LastGreetingSentDate: &sent,
	}

	should, _, err := ShouldSendGreeting(chat, now)
	if err != nil {
		t.Fatalf("ShouldSendGreeting() error = %v", err)
	}
	if should {
		t.Fatal("ShouldSendGreeting() = true, want false")
	}
}
