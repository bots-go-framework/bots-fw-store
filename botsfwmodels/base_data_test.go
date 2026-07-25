package botsfwmodels

import (
	"testing"
	"time"
)

// TestBotBaseData_Validate pins the timestamp rules enforced before a record
// may be persisted: DtCreated and DtUpdated are both required, and DtUpdated
// may never precede DtCreated. This is exactly the check that blocks the
// first write for a brand-new platform user or bot chat when nothing has
// stamped the record yet - see TestBotBaseData_EnsureTimestamps for the seam
// meant to satisfy it before a store implementation inserts the record.
func TestBotBaseData_Validate(t *testing.T) {
	created := time.Date(2026, 7, 25, 9, 0, 0, 0, time.UTC)

	t.Run("zero DtCreated is rejected", func(t *testing.T) {
		data := &BotBaseData{DtUpdated: created}
		if err := data.Validate(); err == nil {
			t.Fatal("Validate() = nil, want an error when DtCreated is zero")
		}
	})

	t.Run("zero DtUpdated is rejected", func(t *testing.T) {
		data := &BotBaseData{DtCreated: created}
		if err := data.Validate(); err == nil {
			t.Fatal("Validate() = nil, want an error when DtUpdated is zero")
		}
	})

	t.Run("DtUpdated before DtCreated is rejected", func(t *testing.T) {
		data := &BotBaseData{DtCreated: created, DtUpdated: created.Add(-time.Second)}
		if err := data.Validate(); err == nil {
			t.Fatal("Validate() = nil, want an error when DtUpdated precedes DtCreated")
		}
	})

	t.Run("equal timestamps are accepted", func(t *testing.T) {
		data := &BotBaseData{DtCreated: created, DtUpdated: created}
		if err := data.Validate(); err != nil {
			t.Fatalf("Validate() = %v, want nil when DtCreated == DtUpdated", err)
		}
	})

	t.Run("DtUpdated after DtCreated is accepted", func(t *testing.T) {
		data := &BotBaseData{DtCreated: created, DtUpdated: created.Add(time.Second)}
		if err := data.Validate(); err != nil {
			t.Fatalf("Validate() = %v, want nil when DtUpdated follows DtCreated", err)
		}
	})
}

// TestBotBaseData_EnsureTimestamps proves the seam a botsfwstore.StateStore
// implementation must call before persisting a freshly-made PlatformUserData
// or BotChatData record. Nothing else stamps these fields: platform field
// setters (e.g. bots-fw-telegram's SetBotUserFields/SetBotChatFields) only
// ever populate platform-specific fields such as FirstName/LastName, so
// without this seam DtCreated/DtUpdated stay zero and Validate rejects the
// record - the bug this closes. The method takes "now" as a parameter rather
// than calling time.Now() itself so a store implementation's clock is
// injectable/overridable in tests, matching the existing SetUpdatedTime seam.
func TestBotBaseData_EnsureTimestamps(t *testing.T) {
	now := time.Date(2026, 7, 25, 9, 0, 0, 0, time.UTC)

	t.Run("a fresh record is unwritable before stamping", func(t *testing.T) {
		data := &BotBaseData{}
		if err := data.Validate(); err == nil {
			t.Fatal("Validate() = nil, want an error before EnsureTimestamps is called")
		}
	})

	t.Run("stamping a fresh record sets both fields to now and makes it valid", func(t *testing.T) {
		data := &BotBaseData{}
		data.EnsureTimestamps(now)
		if !data.DtCreated.Equal(now) {
			t.Errorf("DtCreated = %v, want %v", data.DtCreated, now)
		}
		if !data.DtUpdated.Equal(now) {
			t.Errorf("DtUpdated = %v, want %v", data.DtUpdated, now)
		}
		if err := data.Validate(); err != nil {
			t.Fatalf("Validate() = %v, want nil after EnsureTimestamps", err)
		}
	})

	t.Run("re-stamping an existing record preserves DtCreated and advances DtUpdated", func(t *testing.T) {
		data := &BotBaseData{}
		data.EnsureTimestamps(now)

		later := now.Add(time.Hour)
		data.EnsureTimestamps(later)

		if !data.DtCreated.Equal(now) {
			t.Errorf("DtCreated = %v, want it to stay %v (unchanged by re-stamping)", data.DtCreated, now)
		}
		if !data.DtUpdated.Equal(later) {
			t.Errorf("DtUpdated = %v, want %v", data.DtUpdated, later)
		}
		if err := data.Validate(); err != nil {
			t.Fatalf("Validate() = %v, want nil after re-stamping", err)
		}
	})

	t.Run("re-stamping never leaves DtUpdated before DtCreated", func(t *testing.T) {
		// Guards the invariant structurally: DtCreated is only ever set once
		// (on the first call, from zero), and every call - first or later -
		// sets DtUpdated to that same "now", so the two can never diverge in
		// a way Validate would reject as long as callers pass a
		// non-decreasing clock.
		data := &BotBaseData{}
		data.EnsureTimestamps(now)
		if err := data.Validate(); err != nil {
			t.Fatalf("Validate() = %v, want nil after first stamp", err)
		}
		data.EnsureTimestamps(now.Add(time.Minute))
		if err := data.Validate(); err != nil {
			t.Fatalf("Validate() = %v, want nil after second stamp", err)
		}
	})
}
