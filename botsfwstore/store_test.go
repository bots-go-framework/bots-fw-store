package botsfwstore

import (
	"errors"
	"testing"

	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
)

func TestIdentityValidate(t *testing.T) {
	tests := []struct {
		name     string
		identity Identity
		wantErr  bool
	}{
		{
			name:     "valid",
			identity: Identity{PlatformID: "telegram", BotID: "bot", BotUserID: "user"},
			wantErr:  false,
		},
		{
			name:     "missing platform",
			identity: Identity{PlatformID: "", BotID: "bot", BotUserID: "user"},
			wantErr:  true,
		},
		{
			name:     "missing bot",
			identity: Identity{PlatformID: "telegram", BotID: "", BotUserID: "user"},
			wantErr:  true,
		},
		{
			name:     "missing bot user",
			identity: Identity{PlatformID: "telegram", BotID: "bot", BotUserID: ""},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.identity.Validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Identity.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLinkRequestValidate(t *testing.T) {
	validIdentity := Identity{PlatformID: "telegram", BotID: "bot", BotUserID: "user"}
	validIdentityWithChat := Identity{PlatformID: "telegram", BotID: "bot", BotUserID: "user", ChatID: "chat1"}

	dummyRead := func() botsfwmodels.PlatformUserData { return nil }
	dummyNewUser := func(string) (botsfwmodels.PlatformUserData, error) { return nil, nil }
	dummyNewChat := func(string, bool) (botsfwmodels.BotChatData, error) { return nil, nil }

	tests := []struct {
		name    string
		req     LinkRequest
		wantErr bool
	}{
		{
			name: "valid without chat",
			req: LinkRequest{
				Identity:             validIdentity,
				ReadPlatformUserData: dummyRead,
				NewPlatformUserData:  dummyNewUser,
			},
			wantErr: false,
		},
		{
			name: "valid with chat",
			req: LinkRequest{
				Identity:             validIdentityWithChat,
				ReadPlatformUserData: dummyRead,
				NewPlatformUserData:  dummyNewUser,
				NewChatData:          dummyNewChat,
			},
			wantErr: false,
		},
		{
			name: "invalid identity",
			req: LinkRequest{
				Identity:             Identity{},
				ReadPlatformUserData: dummyRead,
				NewPlatformUserData:  dummyNewUser,
			},
			wantErr: true,
		},
		{
			name: "missing read platform factory",
			req: LinkRequest{
				Identity:            validIdentity,
				NewPlatformUserData: dummyNewUser,
			},
			wantErr: true,
		},
		{
			name: "missing new platform factory",
			req: LinkRequest{
				Identity:             validIdentity,
				ReadPlatformUserData: dummyRead,
			},
			wantErr: true,
		},
		{
			name: "chat set but missing new chat factory",
			req: LinkRequest{
				Identity:             validIdentityWithChat,
				ReadPlatformUserData: dummyRead,
				NewPlatformUserData:  dummyNewUser,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("LinkRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRequireAppUserID(t *testing.T) {
	if err := RequireAppUserID(""); err == nil || !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound on empty app user ID, got %v", err)
	}
	if err := RequireAppUserID("app-123"); err != nil {
		t.Fatalf("expected nil for non-empty app user ID, got %v", err)
	}
}

func TestWebhookUpdateKeyValidate(t *testing.T) {
	tests := []struct {
		name    string
		key     WebhookUpdateKey
		wantErr bool
	}{
		{
			name:    "valid",
			key:     WebhookUpdateKey{PlatformID: "tg", BotID: "b", UpdateID: "u"},
			wantErr: false,
		},
		{
			name:    "empty platform",
			key:     WebhookUpdateKey{PlatformID: "", BotID: "b", UpdateID: "u"},
			wantErr: true,
		},
		{
			name:    "empty bot",
			key:     WebhookUpdateKey{PlatformID: "tg", BotID: "", UpdateID: "u"},
			wantErr: true,
		},
		{
			name:    "empty update",
			key:     WebhookUpdateKey{PlatformID: "tg", BotID: "b", UpdateID: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.key.Validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("WebhookUpdateKey.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNormalizeWebhookUpdateFailureCode(t *testing.T) {
	if code := NormalizeWebhookUpdateFailureCode(WebhookUpdateFailureProcessing); code != WebhookUpdateFailureProcessing {
		t.Fatalf("got %s, want %s", code, WebhookUpdateFailureProcessing)
	}
	if code := NormalizeWebhookUpdateFailureCode(WebhookUpdateFailurePanic); code != WebhookUpdateFailurePanic {
		t.Fatalf("got %s, want %s", code, WebhookUpdateFailurePanic)
	}
	if code := NormalizeWebhookUpdateFailureCode("custom_failure"); code != WebhookUpdateFailureUnknown {
		t.Fatalf("got %s, want %s", code, WebhookUpdateFailureUnknown)
	}
}

func TestWebhookUpdateClaimCanDispatch(t *testing.T) {
	c1 := WebhookUpdateClaim{Status: WebhookUpdateClaimAcquired, LeaseID: "lease1"}
	if !c1.CanDispatch() {
		t.Fatal("expected CanDispatch to be true for acquired claim with leaseID")
	}

	c2 := WebhookUpdateClaim{Status: WebhookUpdateClaimAcquired, LeaseID: ""}
	if c2.CanDispatch() {
		t.Fatal("expected CanDispatch to be false when leaseID is empty")
	}

	c3 := WebhookUpdateClaim{Status: WebhookUpdateClaimCompleted, LeaseID: "lease1"}
	if c3.CanDispatch() {
		t.Fatal("expected CanDispatch to be false when status is completed")
	}

	c4 := WebhookUpdateClaim{Status: WebhookUpdateClaimLeased, LeaseID: "lease1"}
	if c4.CanDispatch() {
		t.Fatal("expected CanDispatch to be false when status is leased")
	}
}
