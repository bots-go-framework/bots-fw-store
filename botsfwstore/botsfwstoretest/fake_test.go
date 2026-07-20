package botsfwstoretest

import (
	"context"
	"testing"

	"github.com/bots-go-framework/bots-fw-store/botsfwstore"
)

func TestFakeStateStoreRequiresExplicitBehavior(t *testing.T) {
	store := &FakeStateStore{}
	if _, err := store.AppUser(context.Background(), "bot", "user"); err == nil {
		t.Fatal("AppUser() error = nil, want an unconfigured-operation error")
	}

	store.AppUserFunc = func(context.Context, string, string) (botsfwstore.AppUser, error) {
		return botsfwstore.AppUser{ID: "user"}, nil
	}
	if user, err := store.AppUser(context.Background(), "bot", "user"); err != nil || user.ID != "user" {
		t.Fatalf("AppUser() = %#v, %v", user, err)
	}
}
