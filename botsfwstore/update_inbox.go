package botsfwstore

import (
	"context"
	"errors"
	"time"
)

// WebhookUpdateKey identifies one provider delivery. The platform and bot are
// part of the key so provider update IDs cannot collide across bot webhooks.
type WebhookUpdateKey struct {
	PlatformID string
	BotID      string
	UpdateID   string
}

func (k WebhookUpdateKey) Validate() error {
	if k.PlatformID == "" || k.BotID == "" || k.UpdateID == "" {
		return errors.New("webhook update key requires platform ID, bot ID, and update ID")
	}
	return nil
}

type WebhookUpdateClaimStatus string

const (
	// WebhookUpdateClaimAcquired gives this caller the lease and permits dispatch.
	WebhookUpdateClaimAcquired WebhookUpdateClaimStatus = "acquired"
	// WebhookUpdateClaimCompleted means an earlier caller completed the update.
	WebhookUpdateClaimCompleted WebhookUpdateClaimStatus = "completed"
	// WebhookUpdateClaimLeased means another caller still owns a valid lease.
	WebhookUpdateClaimLeased WebhookUpdateClaimStatus = "leased"
)

type WebhookUpdateClaim struct {
	Status   WebhookUpdateClaimStatus
	LeaseID  string
	Attempts int
}

// WebhookUpdateFailureCode is a deliberately small operational classification.
// It must not contain an error message, provider payload, identifier, or other
// user-controlled value.
type WebhookUpdateFailureCode string

const (
	WebhookUpdateFailureUnknown    WebhookUpdateFailureCode = "unknown"
	WebhookUpdateFailureProcessing WebhookUpdateFailureCode = "processing_failed"
	WebhookUpdateFailurePanic      WebhookUpdateFailureCode = "panic"
)

// NormalizeWebhookUpdateFailureCode prevents implementations from retaining
// arbitrary caller input. Unknown values are intentionally coarsened.
func NormalizeWebhookUpdateFailureCode(code WebhookUpdateFailureCode) WebhookUpdateFailureCode {
	switch code {
	case WebhookUpdateFailureProcessing, WebhookUpdateFailurePanic:
		return code
	default:
		return WebhookUpdateFailureUnknown
	}
}

func (c WebhookUpdateClaim) CanDispatch() bool {
	return c.Status == WebhookUpdateClaimAcquired && c.LeaseID != ""
}

// WebhookUpdateInbox is a durable, leased inbox. Claim must be atomic: at
// most one live lease for a key may be returned at a time. Failed or expired
// leases may be claimed again; completed updates are never dispatched again.
// Implementations retain completed records for an operator-configured period.
type WebhookUpdateInbox interface {
	ClaimWebhookUpdate(ctx context.Context, key WebhookUpdateKey, leaseUntil time.Time) (WebhookUpdateClaim, error)
	CompleteWebhookUpdate(ctx context.Context, key WebhookUpdateKey, leaseID string) error
	// failureCode is an operational category, never a raw error message or
	// provider payload. Implementations must normalize it before persistence.
	FailWebhookUpdate(ctx context.Context, key WebhookUpdateKey, leaseID string, failureCode WebhookUpdateFailureCode) error
}
