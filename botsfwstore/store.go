// Package botsfwstore defines persistence-neutral state contracts used by bots-fw.
//
// Implementations may use DALgo, SQL, an in-memory map, or a remote service. The
// contracts intentionally do not expose database connections, transactions,
// records, or storage keys.
package botsfwstore

import (
	"context"
	"errors"
	"fmt"

	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
)

// ErrNotFound reports that a requested bot-state item does not exist.
var ErrNotFound = errors.New("bot state not found")

// ErrIdentityConflict reports that the same platform identity was resolved to
// different application users. This normally means an AppUserStore violated
// its idempotency contract or identity state was changed concurrently.
var ErrIdentityConflict = errors.New("bot identity conflict")

// Identity identifies a platform user and, when present, the chat through which
// the bot is interacting with that user. All IDs are platform-provided opaque
// strings. Implementations must not interpret them as numeric IDs.
type Identity struct {
	PlatformID string
	BotID      string
	BotUserID  string
	ChatID     string

	FirstName    string
	LastName     string
	Username     string
	LanguageCode string
}

// Validate checks the identity fields required for bot identity persistence.
func (v Identity) Validate() error {
	if v.PlatformID == "" {
		return errors.New("identity platform ID is required")
	}
	if v.BotID == "" {
		return errors.New("identity bot ID is required")
	}
	if v.BotUserID == "" {
		return errors.New("identity bot user ID is required")
	}
	return nil
}

// PlatformUser is the persistence-neutral view of a platform user.
// It intentionally contains no storage record or key.
type PlatformUser struct {
	ID   string
	Data botsfwmodels.PlatformUserData
}

// AppUser is the persistence-neutral view of an application user.
// Application-specific data remains owned by the consumer application.
type AppUser struct {
	ID   string
	Data botsfwmodels.AppUserData
}

// LinkedIdentity is the state made available to a webhook after identity
// resolution. ChatData is nil for inputs that do not belong to a chat.
type LinkedIdentity struct {
	AppUser      AppUser
	PlatformUser PlatformUser
	ChatData     botsfwmodels.BotChatData
}

// LinkRequest describes the one framework-owned identity use case: make sure a
// platform user, application-user link, and (when ChatID is present) bot chat
// exist. Implementations keep their own persistence atomic; the framework never
// receives their transaction handle.
//
// Factories are called only when an item must be created. They must be free of
// external side effects because a persistence adapter may call them again when
// its database retries a transaction. They keep the framework's
// platform-specific field population out of store implementations.
type LinkRequest struct {
	Identity Identity

	// ReadPlatformUserData returns an empty platform-user value for decoding an
	// existing record. It is separate from NewPlatformUserData because the latter
	// is only valid once the application-user ID is known.
	ReadPlatformUserData func() botsfwmodels.PlatformUserData
	NewPlatformUserData  func(appUserID string) (botsfwmodels.PlatformUserData, error)
	NewChatData          func(appUserID string, accessGranted bool) (botsfwmodels.BotChatData, error)
}

// Validate checks that a link request contains the factories necessary to create
// missing identity state.
func (v LinkRequest) Validate() error {
	if err := v.Identity.Validate(); err != nil {
		return err
	}
	if v.ReadPlatformUserData == nil {
		return errors.New("platform-user reader factory is required")
	}
	if v.NewPlatformUserData == nil {
		return errors.New("new platform-user data factory is required")
	}
	if v.Identity.ChatID != "" && v.NewChatData == nil {
		return errors.New("new chat-data factory is required when chat ID is set")
	}
	return nil
}

// StateStore provides the framework's persistence use cases. It is deliberately
// narrow: command handlers get no generic database access through this contract.
type StateStore interface {
	// EnsureLinked resolves or creates the identity state described by request.
	// A durable implementation must atomically persist its own database changes
	// before returning; router dispatch happens after this call returns.
	EnsureLinked(ctx context.Context, request LinkRequest) (LinkedIdentity, error)

	// PlatformUser returns a platform user without exposing its storage record.
	PlatformUser(ctx context.Context, identity Identity, newData func() botsfwmodels.PlatformUserData) (PlatformUser, error)

	// AppUser returns application-user data associated with a bot. Consumers own
	// the concrete data type and its application persistence.
	AppUser(ctx context.Context, botID, appUserID string) (AppUser, error)

	// SaveChat persists chat state. It is intentionally separate from application
	// business transactions and must not send external messages from a retryable
	// persistence callback.
	SaveChat(ctx context.Context, identity Identity, data botsfwmodels.BotChatData) error

	// SetPlatformUserAccessGranted changes the framework's access flag and returns
	// the updated neutral platform-user view.
	SetPlatformUserAccessGranted(ctx context.Context, identity Identity, newData func() botsfwmodels.PlatformUserData, value bool) (PlatformUser, error)
}

// RequireAppUserID validates an app-user ID returned by a store implementation.
// It is exported for adapters so all implementations report malformed linked
// state consistently.
func RequireAppUserID(id string) error {
	if id == "" {
		return fmt.Errorf("%w: linked app user ID is empty", ErrNotFound)
	}
	return nil
}
