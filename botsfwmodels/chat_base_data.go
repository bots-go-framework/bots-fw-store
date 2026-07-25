package botsfwmodels

import (
	"fmt"
	"github.com/strongo/validation"
	"strings"
	"time"
)

var _ BotChatData = (*ChatBaseData)(nil)

// ChatBaseData hold common properties for bot chat entities not specific to any platform
type ChatBaseData struct {
	changed bool // if true needs to be saved

	//ChatKey // BotID & ChatID

	//// BotUserID is an ID of a bot user who owns this chat
	//BotUserID string // We want it to be indexed and not to omitempty, to find chats without an assigned bot user.

	// BotUserIDs keeps ids of bot users who are members of a group chat or a single user ID if it's a chat with bot
	BotUserIDs []string `firestore:"botUserIDs,omitempty"`

	BotBaseData
	chatState
	chatSettings
	chatVars

	// AppUserIntIDs is kept for legacy reasons
	// Deprecated: replace with `AppUserIDs []string`
	//AppUserIntIDs []int64 `dalgo:",omitempty" datastore:",omitempty" firestore:",omitempty"`

	// IsGroup indicates if bot is added/used in a group chat
	IsGroup bool `firestore:"isGroup,omitempty"`

	Type string `firestore:"type,omitempty"` // TODO: document what is it & why needed or remove

	// Title stores a title of a chat if bot platforms supports named chats
	Title string `firestore:"title,omitempty"`

	// AnalyticsClientIDs stores IDs of analytics clients. For example {"GA": "1234567890.1234567890"}
	AnalyticsClientIDs map[string]string `firestore:"analyticsClientIDs,omitempty"`

	// DtLastInteraction must be set through SetDtLastInteraction() as it also increments InteractionsCount
	DtLastInteraction time.Time `firestore:"dtLastInteraction"`

	// InteractionsCount is a number of interactions with a bot in this chat
	InteractionsCount int `firestore:"interactionsCount,omitempty"`

	// DtForbidden is a time when bot was forbidden to interact with a chat
	DtForbidden time.Time `firestore:"dtForbidden,omitempty"`

	// DtForbiddenLast needs documentation on intended usage. TODO: Consider removing
	DtForbiddenLast time.Time `dalgo:",omitempty" datastore:",omitempty" firestore:",omitempty"`
}

func (e *ChatBaseData) Validate() error {
	//if err := e.ChatKey.Validate(); err != nil {
	//	return err
	//}
	if len(e.BotUserIDs) == 0 {
		return validation.NewErrRecordIsMissingRequiredField("botUserIDs")
	}
	for i, botUserID := range e.BotUserIDs {
		if botUserID == "" {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("botUserIDs[%d]", i), "is empty string")
		}
		if strings.TrimSpace(botUserID) != botUserID {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("botUserIDs[%d]", i), "has leading or trailing spaces")
		}
	}
	if e.DtForbiddenLast.Before(e.DtForbidden) {
		return validation.NewErrBadRecordFieldValue("DtForbidden", fmt.Sprintf("DtForbiddenLast(%v) is before DtForbidden(%v)", e.DtForbiddenLast, e.DtForbidden))
	}
	if e.InteractionsCount < 0 {
		return validation.NewErrBadRecordFieldValue("InteractionsCount", fmt.Sprintf("is less than zero: %d", e.InteractionsCount))
	}
	return nil
}

func (e *ChatBaseData) Base() *ChatBaseData {
	return e
}

// IsChanged indicates that chat data has been changed and record needs to be saved
func (e *ChatBaseData) IsChanged() bool {
	return e.changed || e.chatState.changed
}

//func (e *ChatBaseData) Key() ChatKey {
//	return NewChatKey(e.BotID, e.ChatID)
//}

// GetBotID returns bot ID
//func (e *ChatBaseData) GetBotID() string {
//	return e.BotID
//}

//// SetBotID sets bot ID - TODO: consider removing?
//func (e *ChatBaseData) SetBotID(botID string) {
//	e.changed = true
//	e.BotID = botID
//}

// GetChatID returns chat ID
//func (e *ChatBaseData) GetChatID() string {
//	return e.ChatID
//}

//// SetChatID sets chat ID
//func (e *ChatBaseData) SetChatID(chatID string) {
//	e.changed = true
//	e.ChatID = chatID
//}

// IsGroupChat indicates if it is a group chat
func (e *ChatBaseData) IsGroupChat() bool {
	return e.IsGroup
}

// SetIsGroupChat marks chat as a group one
func (e *ChatBaseData) SetIsGroupChat(v bool) {
	e.changed = true
	e.IsGroup = v
}

// func (e *ChatBaseData) GetBotUserIntID() int {
// 	panic("Should be overwritten in subclass")
// }
//
// func (e *ChatBaseData) GetBotUserStringID() string {
// 	panic("Should be overwritten in subclass")
// }

// SetBotUserID sets bot user ID
func (e *ChatBaseData) SetBotUserID(id interface{}) {
	e.changed = true
	panic(fmt.Sprintf("Should be overwritten in subclass, got: %T=%v", id, id))
}

// SetDtLastInteraction sets date time of last interaction
func (e *ChatBaseData) SetDtLastInteraction(v time.Time) {
	e.changed = true
	e.DtLastInteraction = v
	e.InteractionsCount++
}

// GetGaClientID returns Google Analytics client UUID
// TODO: random implementation should not be here in this module so we do not have dep on random?
//func (e *ChatBaseData) GetGaClientID() string {
//	if len(e.GaClientID) == 0 {
//		e.GaClientID = []byte(random.ID(32))
//	}
//	return string(e.GaClientID)
//}

// SetDtUpdateToNow mark entity updated with now
func (e *ChatBaseData) SetDtUpdateToNow() {
	e.changed = true
	e.DtUpdated = time.Now()
}
