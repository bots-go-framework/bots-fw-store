package botsfwmodels

import (
	"strings"
	"testing"
	"time"
)

// TestPlatformUserBaseDbo_Validate_botIDs pins the fix for a required field
// nothing ever populated: PlatformUserBaseDbo embeds WithRequiredBotIDs,
// whose Validate requires a non-empty BotIDs, but no platform adapter's
// field setter sets it (e.g. bots-fw-telegram's SetBotUserFields receives a
// botID parameter and never uses it), and EnsureLinked itself did not either
// before WithRequiredBotIDs.EnsureBotID was wired in at the store layer
// (bots-fw-store-dalgo). Every brand-new platform user therefore failed
// Validate on its very first write against any dalgo backend that validates
// on insert (a real Firestore-backed one does; the framework's own tests use
// dalgo2memory, which does not - see that package's doc comment).
func TestPlatformUserBaseDbo_Validate_botIDs(t *testing.T) {
	now := time.Date(2026, 7, 25, 9, 0, 0, 0, time.UTC)

	t.Run("a fresh record fails on both timestamps and botIDs", func(t *testing.T) {
		data := &PlatformUserBaseDbo{}
		if err := data.Validate(); err == nil {
			t.Fatal("Validate() = nil, want an error for a fresh record")
		}
	})

	t.Run("stamping timestamps alone is not enough - botIDs is still missing", func(t *testing.T) {
		data := &PlatformUserBaseDbo{}
		data.EnsureTimestamps(now)
		err := data.Validate()
		if err == nil {
			t.Fatal("Validate() = nil, want an error because BotIDs is still empty")
		}
		if !strings.Contains(err.Error(), "botIDs") {
			t.Errorf("Validate() = %v, want it to name botIDs", err)
		}
	})

	t.Run("stamping timestamps and ensuring the bot id makes the record valid", func(t *testing.T) {
		data := &PlatformUserBaseDbo{}
		data.EnsureTimestamps(now)
		data.EnsureBotID("bot-1")
		if err := data.Validate(); err != nil {
			t.Fatalf("Validate() = %v, want nil once both timestamps and botIDs are set", err)
		}
	})
}
