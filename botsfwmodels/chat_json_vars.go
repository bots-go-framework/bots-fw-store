package botsfwmodels

import (
	"encoding/json"
	"fmt"
)

// SetJSONVar marshals value to JSON and stores it under key using SetVar.
// Returns an error if marshalling fails.
func SetJSONVar(chat BotChatData, key string, value any) error {
	b, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("SetJSONVar(%q): failed to marshal value: %w", key, err)
	}
	chat.SetVar(key, string(b))
	return nil
}

// GetJSONVar retrieves the JSON string stored under key and unmarshals it into into.
// Returns nil without touching into when the key is absent or empty.
// Callers must pass a non-nil pointer as into, just like json.Unmarshal.
func GetJSONVar(chat BotChatData, key string, into any) error {
	raw := chat.GetVar(key)
	if raw == "" {
		return nil
	}
	if err := json.Unmarshal([]byte(raw), into); err != nil {
		return fmt.Errorf("GetJSONVar(%q): failed to unmarshal value: %w", key, err)
	}
	return nil
}
