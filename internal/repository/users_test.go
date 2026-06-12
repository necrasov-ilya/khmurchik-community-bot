package repository

import "testing"

func TestNormalizeUsername(t *testing.T) {
	got := NormalizeUsername(" @SomeUser ")
	if got != "someuser" {
		t.Fatalf("NormalizeUsername() = %q, want %q", got, "someuser")
	}
}
