package botsfwmodels

import (
	"testing"
)

func TestBotChatEntity_PopStepsFromAwaitingReplyUpToSpecificParent(t *testing.T) {
	chatEntity := ChatBaseData{}

	chatEntity.AwaitingReplyTo = "step1/step2/step3"
	chatEntity.PopStepsFromAwaitingReplyUpToSpecificParent("step2")
	if chatEntity.AwaitingReplyTo != "step1/step2" {
		t.Errorf("Failed to remove last step3. AwaitingReplyTo: %s", chatEntity.AwaitingReplyTo)
	}

	chatEntity.AwaitingReplyTo = "step1/step2"
	chatEntity.PopStepsFromAwaitingReplyUpToSpecificParent("step2")
	if chatEntity.AwaitingReplyTo != "step1/step2" {
		t.Errorf("Failed to remove last step3. AwaitingReplyTo: %s", chatEntity.AwaitingReplyTo)
	}
}
