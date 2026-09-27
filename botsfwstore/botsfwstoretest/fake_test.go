package botsfwstoretest

import (
	"context"
	"testing"
	"time"

	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/bots-go-framework/bots-fw-store/botsfwstore"
)

func TestFakeStateStore_AllMethods(t *testing.T) {
	ctx := context.Background()
	key := botsfwstore.WebhookUpdateKey{PlatformID: "p", BotID: "b", UpdateID: "u"}
	identity := botsfwstore.Identity{PlatformID: "p", BotID: "b", BotUserID: "u"}
	now := time.Now()

	// 1. nil receiver
	var nilStore *FakeStateStore
	if _, err := nilStore.ClaimWebhookUpdate(ctx, key, now); err == nil {
		t.Fatal("expected error on nilStore.ClaimWebhookUpdate")
	}
	if err := nilStore.CompleteWebhookUpdate(ctx, key, "lease"); err == nil {
		t.Fatal("expected error on nilStore.CompleteWebhookUpdate")
	}
	if err := nilStore.FailWebhookUpdate(ctx, key, "lease", botsfwstore.WebhookUpdateFailurePanic); err == nil {
		t.Fatal("expected error on nilStore.FailWebhookUpdate")
	}
	if _, err := nilStore.EnsureLinked(ctx, botsfwstore.LinkRequest{}); err == nil {
		t.Fatal("expected error on nilStore.EnsureLinked")
	}
	if _, err := nilStore.PlatformUser(ctx, identity, nil); err == nil {
		t.Fatal("expected error on nilStore.PlatformUser")
	}
	if _, err := nilStore.AppUser(ctx, "b", "u"); err == nil {
		t.Fatal("expected error on nilStore.AppUser")
	}
	if err := nilStore.SaveChat(ctx, identity, nil); err == nil {
		t.Fatal("expected error on nilStore.SaveChat")
	}
	if _, err := nilStore.SetPlatformUserAccessGranted(ctx, identity, nil, true); err == nil {
		t.Fatal("expected error on nilStore.SetPlatformUserAccessGranted")
	}

	// 2. unconfigured store
	store := &FakeStateStore{}
	if _, err := store.ClaimWebhookUpdate(ctx, key, now); err == nil {
		t.Fatal("expected error on unconfigured ClaimWebhookUpdate")
	}
	if err := store.CompleteWebhookUpdate(ctx, key, "lease"); err == nil {
		t.Fatal("expected error on unconfigured CompleteWebhookUpdate")
	}
	if err := store.FailWebhookUpdate(ctx, key, "lease", botsfwstore.WebhookUpdateFailurePanic); err == nil {
		t.Fatal("expected error on unconfigured FailWebhookUpdate")
	}
	if _, err := store.EnsureLinked(ctx, botsfwstore.LinkRequest{}); err == nil {
		t.Fatal("expected error on unconfigured EnsureLinked")
	}
	if _, err := store.PlatformUser(ctx, identity, nil); err == nil {
		t.Fatal("expected error on unconfigured PlatformUser")
	}
	if _, err := store.AppUser(ctx, "b", "u"); err == nil {
		t.Fatal("expected error on unconfigured AppUser")
	}
	if err := store.SaveChat(ctx, identity, nil); err == nil {
		t.Fatal("expected error on unconfigured SaveChat")
	}
	if _, err := store.SetPlatformUserAccessGranted(ctx, identity, nil, true); err == nil {
		t.Fatal("expected error on unconfigured SetPlatformUserAccessGranted")
	}

	// 3. configured functions
	store.ClaimWebhookUpdateFunc = func(context.Context, botsfwstore.WebhookUpdateKey, time.Time) (botsfwstore.WebhookUpdateClaim, error) {
		return botsfwstore.WebhookUpdateClaim{Status: botsfwstore.WebhookUpdateClaimAcquired, LeaseID: "l1"}, nil
	}
	claim, err := store.ClaimWebhookUpdate(ctx, key, now)
	if err != nil || claim.LeaseID != "l1" {
		t.Fatalf("ClaimWebhookUpdate failed: %v", err)
	}

	store.CompleteWebhookUpdateFunc = func(context.Context, botsfwstore.WebhookUpdateKey, string) error {
		return nil
	}
	if err := store.CompleteWebhookUpdate(ctx, key, "l1"); err != nil {
		t.Fatalf("CompleteWebhookUpdate failed: %v", err)
	}

	store.FailWebhookUpdateFunc = func(context.Context, botsfwstore.WebhookUpdateKey, string, botsfwstore.WebhookUpdateFailureCode) error {
		return nil
	}
	if err := store.FailWebhookUpdate(ctx, key, "l1", botsfwstore.WebhookUpdateFailurePanic); err != nil {
		t.Fatalf("FailWebhookUpdate failed: %v", err)
	}

	store.EnsureLinkedFunc = func(context.Context, botsfwstore.LinkRequest) (botsfwstore.LinkedIdentity, error) {
		return botsfwstore.LinkedIdentity{AppUser: botsfwstore.AppUser{ID: "app1"}}, nil
	}
	linked, err := store.EnsureLinked(ctx, botsfwstore.LinkRequest{})
	if err != nil || linked.AppUser.ID != "app1" {
		t.Fatalf("EnsureLinked failed: %v", err)
	}

	store.PlatformUserFunc = func(context.Context, botsfwstore.Identity, func() botsfwmodels.PlatformUserData) (botsfwstore.PlatformUser, error) {
		return botsfwstore.PlatformUser{ID: "pu1"}, nil
	}
	pu, err := store.PlatformUser(ctx, identity, nil)
	if err != nil || pu.ID != "pu1" {
		t.Fatalf("PlatformUser failed: %v", err)
	}

	store.AppUserFunc = func(context.Context, string, string) (botsfwstore.AppUser, error) {
		return botsfwstore.AppUser{ID: "au1"}, nil
	}
	au, err := store.AppUser(ctx, "b", "u")
	if err != nil || au.ID != "au1" {
		t.Fatalf("AppUser failed: %v", err)
	}

	store.SaveChatFunc = func(context.Context, botsfwstore.Identity, botsfwmodels.BotChatData) error {
		return nil
	}
	if err := store.SaveChat(ctx, identity, nil); err != nil {
		t.Fatalf("SaveChat failed: %v", err)
	}

	store.SetPlatformUserAccessGrantedFunc = func(context.Context, botsfwstore.Identity, func() botsfwmodels.PlatformUserData, bool) (botsfwstore.PlatformUser, error) {
		return botsfwstore.PlatformUser{ID: "pu_granted"}, nil
	}
	puGranted, err := store.SetPlatformUserAccessGranted(ctx, identity, nil, true)
	if err != nil || puGranted.ID != "pu_granted" {
		t.Fatalf("SetPlatformUserAccessGranted failed: %v", err)
	}
}
