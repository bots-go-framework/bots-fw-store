package botsfwmodels

import (
	"fmt"
	"github.com/strongo/validation"
	"slices"
	"strings"
)

type WithBotIDs struct {
	BotIDs []string `json:"botIDs,omitempty" dalgo:"botIDs,omitempty,noindex" firestore:"botIDs,omitempty"`
}

func (v WithBotIDs) Validate() error {
	for i, botID := range v.BotIDs {
		if botID == "" {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("botIDs[%d]", i), "is empty")
		}
		if strings.TrimSpace(botID) != botID {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("botIDs[%d]", i), "has leading or trailing spaces")
		}
	}
	return nil
}

type WithRequiredBotIDs WithBotIDs

func (v WithRequiredBotIDs) Validate() error {
	if len(v.BotIDs) == 0 {
		return validation.NewErrRecordIsMissingRequiredField("botIDs")
	}
	return WithBotIDs(v).Validate()
}

// EnsureBotID appends botID to BotIDs unless it is already present. It is the
// self-heal seam a botsfwstore.StateStore implementation must call before
// inserting or updating a PlatformUserBaseDbo record - mirroring
// BotBaseData.EnsureTimestamps - so that Validate never rejects a record for
// missing BotIDs.
//
// Nothing else populates BotIDs: platform adapters (e.g. bots-fw-telegram's
// SetBotUserFields) only ever set platform-specific fields such as
// FirstName/LastName, and a record written before this field existed, or by
// an application flow that pre-provisions a platform user, has none at all.
// Calling EnsureBotID on every write means such a record repairs itself
// instead of failing WithRequiredBotIDs.Validate forever.
//
// It never removes or duplicates an existing id, so calling it again for the
// same bot is a no-op, and calling it for a different bot appends alongside
// whatever is already there - the same platform-user record can be shared by
// more than one bot. An empty botID is a no-op rather than a panic: this is a
// self-heal seam meant to run unconditionally on every write, not a strict
// setter a caller must get right.
func (v *WithRequiredBotIDs) EnsureBotID(botID string) {
	if strings.TrimSpace(botID) == "" {
		return
	}
	if slices.Contains(v.BotIDs, botID) {
		return
	}
	v.BotIDs = append(v.BotIDs, botID)
}
