// Package botsfwstoretest provides persistence-neutral test doubles for store consumers.
package botsfwstoretest

import (
	"context"
	"fmt"

	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/bots-go-framework/bots-fw-store/botsfwstore"
)

// FakeStateStore is a configurable StateStore test double. A missing function
// reports an error so tests cannot accidentally rely on zero-value persistence
// behavior.
type FakeStateStore struct {
	EnsureLinkedFunc                 func(context.Context, botsfwstore.LinkRequest) (botsfwstore.LinkedIdentity, error)
	PlatformUserFunc                 func(context.Context, botsfwstore.Identity, func() botsfwmodels.PlatformUserData) (botsfwstore.PlatformUser, error)
	AppUserFunc                      func(context.Context, string, string) (botsfwstore.AppUser, error)
	SaveChatFunc                     func(context.Context, botsfwstore.Identity, botsfwmodels.BotChatData) error
	SetPlatformUserAccessGrantedFunc func(context.Context, botsfwstore.Identity, func() botsfwmodels.PlatformUserData, bool) (botsfwstore.PlatformUser, error)
}

var _ botsfwstore.StateStore = (*FakeStateStore)(nil)

func (s *FakeStateStore) EnsureLinked(ctx context.Context, request botsfwstore.LinkRequest) (botsfwstore.LinkedIdentity, error) {
	if s != nil && s.EnsureLinkedFunc != nil {
		return s.EnsureLinkedFunc(ctx, request)
	}
	return botsfwstore.LinkedIdentity{}, fmt.Errorf("FakeStateStore.EnsureLinked is not configured")
}

func (s *FakeStateStore) PlatformUser(ctx context.Context, identity botsfwstore.Identity, newData func() botsfwmodels.PlatformUserData) (botsfwstore.PlatformUser, error) {
	if s != nil && s.PlatformUserFunc != nil {
		return s.PlatformUserFunc(ctx, identity, newData)
	}
	return botsfwstore.PlatformUser{}, fmt.Errorf("FakeStateStore.PlatformUser is not configured")
}

func (s *FakeStateStore) AppUser(ctx context.Context, botID, appUserID string) (botsfwstore.AppUser, error) {
	if s != nil && s.AppUserFunc != nil {
		return s.AppUserFunc(ctx, botID, appUserID)
	}
	return botsfwstore.AppUser{}, fmt.Errorf("FakeStateStore.AppUser is not configured")
}

func (s *FakeStateStore) SaveChat(ctx context.Context, identity botsfwstore.Identity, data botsfwmodels.BotChatData) error {
	if s != nil && s.SaveChatFunc != nil {
		return s.SaveChatFunc(ctx, identity, data)
	}
	return fmt.Errorf("FakeStateStore.SaveChat is not configured")
}

func (s *FakeStateStore) SetPlatformUserAccessGranted(ctx context.Context, identity botsfwstore.Identity, newData func() botsfwmodels.PlatformUserData, value bool) (botsfwstore.PlatformUser, error) {
	if s != nil && s.SetPlatformUserAccessGrantedFunc != nil {
		return s.SetPlatformUserAccessGrantedFunc(ctx, identity, newData, value)
	}
	return botsfwstore.PlatformUser{}, fmt.Errorf("FakeStateStore.SetPlatformUserAccessGranted is not configured")
}
