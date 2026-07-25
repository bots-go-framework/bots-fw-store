package botsfwmodels

import (
	"github.com/strongo/validation"
	"testing"
)

func TestWithBotIDs_Validate(t *testing.T) {
	type fields struct {
		BotIDs []string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Should pass for good records",
			fields: fields{
				BotIDs: []string{"bot1", "bot2", "bot3"},
			},
			wantErr: false,
		},
		{
			name: "Should pass for nil",
			fields: fields{
				BotIDs: nil,
			},
			wantErr: false,
		},
		{
			name: "Should pass for empty",
			fields: fields{
				BotIDs: []string{},
			},
			wantErr: false,
		},
		{
			// Regression: this used to be a copy of "Should pass for good
			// records" wired to wantErr: false, so a padded id silently
			// passed despite the test's own name. WithBotIDs.Validate only
			// rejected a fully blank/whitespace-only entry; it never checked
			// a non-blank entry for surrounding whitespace, unlike the
			// (already-fixed) equivalent check on ChatBaseData.BotUserIDs.
			name: "Should fail for non trimmed",
			fields: fields{
				BotIDs: []string{" bot1"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := WithBotIDs{
				BotIDs: tt.fields.BotIDs,
			}
			if err := v.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWithRequiredBotIDs_Validate(t *testing.T) {
	tests := []struct {
		name    string
		v       WithRequiredBotIDs
		wantErr bool
	}{
		{
			name: "Should pass",
			v: WithRequiredBotIDs{
				BotIDs: []string{"bot1"},
			},
			wantErr: false,
		},
		{
			name:    "Should fail on nil",
			v:       WithRequiredBotIDs{},
			wantErr: true,
		},
		{
			name:    "Should fail on empty",
			v:       WithRequiredBotIDs{BotIDs: []string{}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.v.Validate()
			if err == nil && tt.wantErr || err != nil && !tt.wantErr || tt.wantErr && !validation.IsBadFieldValueError(err) {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestWithRequiredBotIDs_EnsureBotID proves the self-heal seam behaves like
// BotBaseData.EnsureTimestamps: a legacy record missing BotIDs (or one
// created by a platform adapter's field setter, which never populates it -
// see bots-fw-telegram's SetBotUserFields) repairs itself on its next write
// instead of staying permanently invalid; it never duplicates an id it
// already has; and it never clobbers a different bot's id already present,
// since the same platform-user record can be shared by more than one bot.
func TestWithRequiredBotIDs_EnsureBotID(t *testing.T) {
	t.Run("appends to a record with no BotIDs", func(t *testing.T) {
		v := &WithRequiredBotIDs{}
		if err := v.Validate(); err == nil {
			t.Fatal("Validate() = nil, want an error before EnsureBotID is called")
		}
		v.EnsureBotID("bot-1")
		if got := v.BotIDs; len(got) != 1 || got[0] != "bot-1" {
			t.Fatalf("BotIDs = %v, want [bot-1]", got)
		}
		if err := v.Validate(); err != nil {
			t.Fatalf("Validate() = %v, want nil after EnsureBotID", err)
		}
	})

	t.Run("does not duplicate an id already present", func(t *testing.T) {
		v := &WithRequiredBotIDs{BotIDs: []string{"bot-1"}}
		v.EnsureBotID("bot-1")
		if got := v.BotIDs; len(got) != 1 || got[0] != "bot-1" {
			t.Fatalf("BotIDs = %v, want it to stay [bot-1] (no duplicate)", got)
		}
	})

	t.Run("appends without clobbering another bot's id", func(t *testing.T) {
		v := &WithRequiredBotIDs{BotIDs: []string{"bot-1"}}
		v.EnsureBotID("bot-2")
		if got := v.BotIDs; len(got) != 2 || got[0] != "bot-1" || got[1] != "bot-2" {
			t.Fatalf("BotIDs = %v, want [bot-1 bot-2]", got)
		}
	})

	t.Run("empty bot id is a no-op", func(t *testing.T) {
		v := &WithRequiredBotIDs{}
		v.EnsureBotID("")
		if len(v.BotIDs) != 0 {
			t.Fatalf("BotIDs = %v, want it to stay empty", v.BotIDs)
		}
	})

	t.Run("whitespace-only bot id is a no-op", func(t *testing.T) {
		v := &WithRequiredBotIDs{}
		v.EnsureBotID("   ")
		if len(v.BotIDs) != 0 {
			t.Fatalf("BotIDs = %v, want it to stay empty", v.BotIDs)
		}
	})
}
