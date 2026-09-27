package botsfwmodels

import (
	"testing"
	"time"
)

func TestChatData_NewChatIDPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on empty botChatID")
		}
	}()
	NewChatID("bot1", "")
}

func TestChatJSONVars_MarshalError(t *testing.T) {
	chat := &ChatBaseData{}
	// Channels cannot be marshalled to JSON
	ch := make(chan int)
	err := SetJSONVar(chat, "ch", ch)
	if err == nil {
		t.Fatal("expected error from SetJSONVar when marshalling channel")
	}
}

func TestPlatformUserBaseDbo_BaseData(t *testing.T) {
	p := &PlatformUserBaseDbo{
		FirstName: "John",
	}
	if p.BaseData() != p {
		t.Fatal("expected BaseData() to return receiver pointer")
	}
}

func TestWithBotIDs_EmptyItemValidate(t *testing.T) {
	w := WithBotIDs{
		BotIDs: []string{""},
	}
	err := w.Validate()
	if err == nil {
		t.Fatal("expected error on empty botID in slice")
	}
}

func TestChatKey_PanicsAndValidation(t *testing.T) {
	// panic on empty botID
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic on empty botID")
			}
		}()
		NewChatKey("", "chat1")
	}()

	// panic on empty chatID
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic on empty chatID")
			}
		}()
		NewChatKey("bot1", "")
	}()

	// Validate missing BotID
	k1 := ChatKey{BotID: "", ChatID: "c1"}
	if err := k1.Validate(); err == nil {
		t.Fatal("expected validation error on missing BotID")
	}

	// Validate missing ChatID
	k2 := ChatKey{BotID: "b1", ChatID: ""}
	if err := k2.Validate(); err == nil {
		t.Fatal("expected validation error on missing ChatID")
	}
}

func TestWithBotUserIDs_PanicsAndDuplicates(t *testing.T) {
	w := &WithBotUserIDs{}

	// panic on empty userID
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic on empty userID")
			}
		}()
		w.SetBotUserID("tg", "bot", "")
	}()

	// panic on untrimmed userID
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic on untrimmed userID")
			}
		}()
		w.SetBotUserID("tg", "bot", " user ")
	}()

	// Add user and duplicate
	w.SetBotUserID("tg", "bot", "u1")
	if len(w.BotUserIDs) != 1 {
		t.Fatalf("expected 1 bot user ID, got %d", len(w.BotUserIDs))
	}
	// duplicate should return without appending
	w.SetBotUserID("tg", "bot", "u1")
	if len(w.BotUserIDs) != 1 {
		t.Fatalf("expected still 1 bot user ID, got %d", len(w.BotUserIDs))
	}

	// Validate with empty bot user ID
	wInvalid := &WithBotUserIDs{
		BotUserIDs: []string{"   "},
	}
	if err := wInvalid.Validate(); err == nil {
		t.Fatal("expected validation error on empty botUserID")
	}
}

func TestChatBaseData_MethodsAndValidation(t *testing.T) {
	chat := &ChatBaseData{}

	// Base()
	if chat.Base() != chat {
		t.Fatal("expected Base() to return receiver pointer")
	}

	// IsChanged()
	if chat.IsChanged() {
		t.Fatal("expected IsChanged to be false initially")
	}

	// GroupChat
	if chat.IsGroupChat() {
		t.Fatal("expected IsGroupChat to be false initially")
	}
	chat.SetIsGroupChat(true)
	if !chat.IsGroupChat() {
		t.Fatal("expected IsGroupChat to be true after SetIsGroupChat")
	}
	if !chat.IsChanged() {
		t.Fatal("expected IsChanged to be true after SetIsGroupChat")
	}

	// SetBotUserID panic
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic on SetBotUserID")
			}
		}()
		chat.SetBotUserID("user1")
	}()

	// SetDtLastInteraction
	now := time.Now()
	chat.SetDtLastInteraction(now)
	if chat.DtLastInteraction != now {
		t.Fatalf("expected DtLastInteraction %v, got %v", now, chat.DtLastInteraction)
	}
	if chat.InteractionsCount != 1 {
		t.Fatalf("expected InteractionsCount 1, got %d", chat.InteractionsCount)
	}

	// SetDtUpdateToNow
	chat.SetDtUpdateToNow()
	if chat.DtUpdated.IsZero() {
		t.Fatal("expected DtUpdated to be set by SetDtUpdateToNow")
	}

	// Validate negative InteractionsCount
	chat.DtCreated = now
	chat.DtUpdated = now
	chat.BotUserIDs = []string{"u1"}
	chat.InteractionsCount = -1
	if err := chat.Validate(); err == nil {
		t.Fatal("expected validation error on negative InteractionsCount")
	}
}

func TestBotBaseData_Methods(t *testing.T) {
	b := &BotBaseData{}

	now := time.Now()
	b.SetUpdatedTime(now)
	if b.DtUpdated != now {
		t.Fatalf("expected DtUpdated %v, got %v", now, b.DtUpdated)
	}

	if b.GetAppUserID() != "" {
		t.Fatal("expected empty AppUserID initially")
	}
	b.SetAppUserID("app1")
	if b.GetAppUserID() != "app1" {
		t.Fatalf("expected 'app1', got %q", b.GetAppUserID())
	}

	if b.IsAccessGranted() {
		t.Fatal("expected IsAccessGranted to be false initially")
	}

	// SetAccessGranted changes value -> returns true
	if !b.SetAccessGranted(true) {
		t.Fatal("expected SetAccessGranted(true) to return true when changed")
	}
	if !b.IsAccessGranted() {
		t.Fatal("expected IsAccessGranted to be true")
	}

	// SetAccessGranted same value -> returns false
	if b.SetAccessGranted(true) {
		t.Fatal("expected SetAccessGranted(true) to return false when unchanged")
	}
}
