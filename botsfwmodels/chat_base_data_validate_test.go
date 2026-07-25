package botsfwmodels

import (
	"strings"
	"testing"
	"time"
)

// TestChatBaseData_Validate_botUserIDs pins the botUserIDs checks. The
// leading/trailing-space case is here because the comparison was inverted for
// roughly eighteen months: `TrimSpace(id) == id` (i.e. the id is CLEAN) raised
// "has leading or trailing spaces", so Validate rejected every well-formed id
// and could never have rejected a padded one. Nothing caught it because the
// framework's own tests wire dalgo2memory, which does not call Validate on
// insert.
func TestChatBaseData_Validate_botUserIDs(t *testing.T) {
	// A minimally valid record, so each case below varies only BotUserIDs.
	// DtCreated/DtUpdated must be set here too: ChatBaseData.Validate now
	// delegates to BotBaseData.Validate (see
	// TestChatBaseData_Validate_delegatesToBotBaseData below), so a zero
	// BotBaseData would fail every case in this table for a reason unrelated
	// to what it is actually testing.
	validTimestamp := time.Date(2026, 7, 25, 9, 0, 0, 0, time.UTC)
	newChat := func(botUserIDs ...string) *ChatBaseData {
		return &ChatBaseData{
			BotBaseData: BotBaseData{DtCreated: validTimestamp, DtUpdated: validTimestamp},
			BotUserIDs:  botUserIDs,
		}
	}

	t.Run("clean id is accepted", func(t *testing.T) {
		// The regression: this is the ordinary case — a real Telegram user id
		// — and it used to fail.
		if err := newChat("123456789").Validate(); err != nil {
			t.Fatalf("Validate() = %v, want nil for a clean botUserID", err)
		}
	})

	t.Run("several clean ids are accepted", func(t *testing.T) {
		if err := newChat("1", "22", "333").Validate(); err != nil {
			t.Fatalf("Validate() = %v, want nil", err)
		}
	})

	t.Run("no ids is rejected", func(t *testing.T) {
		if err := newChat().Validate(); err == nil {
			t.Fatal("Validate() = nil, want an error when botUserIDs is empty")
		}
	})

	t.Run("empty id is rejected", func(t *testing.T) {
		err := newChat("").Validate()
		if err == nil {
			t.Fatal("Validate() = nil, want an error for an empty botUserID")
		}
		if !strings.Contains(err.Error(), "is empty string") {
			t.Errorf("Validate() = %v, want it to name the empty string", err)
		}
	})

	for _, padded := range []string{" 123", "123 ", "\t123", "123\n", " 123 "} {
		t.Run("padded id is rejected: "+strings.ReplaceAll(padded, "\n", `\n`), func(t *testing.T) {
			err := newChat(padded).Validate()
			if err == nil {
				t.Fatalf("Validate() = nil, want an error for botUserID %q", padded)
			}
			if !strings.Contains(err.Error(), "leading or trailing spaces") {
				t.Errorf("Validate() = %v, want it to name the padding", err)
			}
		})
	}
}

// TestChatBaseData_Validate_dtForbidden keeps the neighbouring check honest,
// so the fix above cannot be "verified" by a test file that only ever
// exercises one branch.
func TestChatBaseData_Validate_dtForbidden(t *testing.T) {
	created := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	forbidden := time.Date(2026, 7, 25, 12, 0, 0, 0, time.UTC)

	chat := &ChatBaseData{
		BotBaseData:     BotBaseData{DtCreated: created, DtUpdated: created},
		BotUserIDs:      []string{"123"},
		DtForbidden:     forbidden,
		DtForbiddenLast: forbidden.Add(-time.Hour),
	}
	if err := chat.Validate(); err == nil {
		t.Fatal("Validate() = nil, want an error when DtForbiddenLast precedes DtForbidden")
	}

	chat.DtForbiddenLast = forbidden.Add(time.Hour)
	if err := chat.Validate(); err != nil {
		t.Fatalf("Validate() = %v, want nil when DtForbiddenLast follows DtForbidden", err)
	}
}

// TestChatBaseData_Validate_delegatesToBotBaseData pins the fix for a gap in
// the same shape as botUserIDs above but for the whole embedded BotBaseData:
// ChatBaseData.Validate checked its own fields (BotUserIDs, DtForbidden vs
// DtForbiddenLast, InteractionsCount) but never called BotBaseData.Validate,
// unlike PlatformUserBaseDbo.Validate, which already delegates to it. So a
// chat record's DtCreated/DtUpdated were never actually checked by Validate,
// even though EnsureLinked relies on BotBaseData.Validate's rules (via
// BotBaseData.EnsureTimestamps) to decide whether a record is safe to
// persist.
func TestChatBaseData_Validate_delegatesToBotBaseData(t *testing.T) {
	chat := &ChatBaseData{BotUserIDs: []string{"123"}} // zero DtCreated/DtUpdated
	err := chat.Validate()
	if err == nil {
		t.Fatal("Validate() = nil, want an error because BotBaseData is invalid (zero DtCreated/DtUpdated)")
	}
	if !strings.Contains(err.Error(), "dtCreated") {
		t.Errorf("Validate() = %v, want it to name dtCreated", err)
	}

	chat.BotBaseData.EnsureTimestamps(time.Date(2026, 7, 25, 9, 0, 0, 0, time.UTC))
	if err := chat.Validate(); err != nil {
		t.Fatalf("Validate() = %v, want nil once BotBaseData is stamped", err)
	}
}
