package botsfwmodels

import (
	"github.com/strongo/validation"
	"time"
)

// BotBaseData holds properties common to all bot entities
type BotBaseData struct {
	//user.OwnedByUserWithID
	AppUserID string `json:"appUserID" firestore:"appUserID"` // intentionally indexed & do NOT omitempty (so we can find records with empty AppUserID)

	// AppUserIntID is a strongly typed ID of the user
	// Deprecated: use AppUserID instead, remove once all bots are migrated
	//AppUserIntID int64 `json:",omitempty" datastore:",omitempty" firestore:",omitempty"`

	DtCreated time.Time `json:"dtCreated" datastore:"dtCreated" firestore:"dtCreated"`
	DtUpdated time.Time `json:"dtUpdated" datastore:"dtUpdated" firestore:"dtUpdated"`

	// AccessGranted indicates if access to the bot has been granted
	AccessGranted bool `json:"isAccessGranted" dalgo:"isAccessGranted" datastore:"isAccessGranted" firestore:"isAccessGranted"`
}

// Validate returns error if data is invalid
func (e *BotBaseData) Validate() error {
	//if e.AppUserID == "" {
	//	return validation.NewErrRecordIsMissingRequiredField("appUserID")
	//}
	if e.DtCreated.IsZero() {
		return validation.NewErrRecordIsMissingRequiredField("dtCreated")
	}
	if e.DtUpdated.IsZero() {
		return validation.NewErrRecordIsMissingRequiredField("dtUpdated")
	}
	if e.DtUpdated.Before(e.DtCreated) {
		return validation.NewErrBadRecordFieldValue("dtUpdated", "is before dtCreated")
	}
	return nil
}

// SetUpdatedTime sets updated time
func (e *BotBaseData) SetUpdatedTime(v time.Time) {
	e.DtUpdated = v
}

// EnsureTimestamps stamps the record for persistence using a caller-supplied
// clock. It is the seam a botsfwstore.StateStore implementation must call
// before inserting or updating a PlatformUserData or BotChatData record, so
// that Validate never rejects it for missing DtCreated/DtUpdated.
//
// On the first call for a record (DtCreated still zero) it sets both
// DtCreated and DtUpdated to now. On every later call it only refreshes
// DtUpdated, so an update can never rewrite the original creation time. The
// two fields are always set from the same now on any given call, so
// DtUpdated can never end up before DtCreated as long as callers pass a
// non-decreasing clock across calls for the same record (e.g. time.Now, or a
// fixed/advancing clock in tests).
//
// Callers supply now explicitly rather than this method calling time.Now()
// itself, so a store implementation's clock stays injectable/overridable in
// tests instead of being hardcoded.
func (e *BotBaseData) EnsureTimestamps(now time.Time) {
	if e.DtCreated.IsZero() {
		e.DtCreated = now
	}
	e.DtUpdated = now
}

// GetAppUserID returns app user ID
func (e *BotBaseData) GetAppUserID() string {
	if e.AppUserID != "" {
		return e.AppUserID
	}
	//if e.AppUserIntID != 0 {
	//	return strconv.FormatInt(e.AppUserIntID, 10)
	//}
	return ""
}

// SetAppUserID sets app user ID
func (e *BotBaseData) SetAppUserID(s string) {
	e.AppUserID = s
}

// IsAccessGranted indicates if access to the bot has been granted
func (e *BotBaseData) IsAccessGranted() bool {
	return e.AccessGranted
}

// SetAccessGranted mark that access has been granted
func (e *BotBaseData) SetAccessGranted(value bool) bool {
	if e.AccessGranted != value {
		e.AccessGranted = value
		return true
	}
	return false
}
