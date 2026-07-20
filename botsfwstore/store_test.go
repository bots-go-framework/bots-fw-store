package botsfwstore

import "testing"

func TestIdentityValidate(t *testing.T) {
	identity := Identity{PlatformID: "telegram", BotID: "bot", BotUserID: "user"}
	if err := identity.Validate(); err != nil {
		t.Fatalf("Identity.Validate() returned an error: %v", err)
	}
}

func TestLinkRequestValidateRequiresFactories(t *testing.T) {
	request := LinkRequest{Identity: Identity{PlatformID: "telegram", BotID: "bot", BotUserID: "user"}}
	if err := request.Validate(); err == nil {
		t.Fatal("LinkRequest.Validate() returned nil for a request without factories")
	}
}
