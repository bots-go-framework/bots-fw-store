package botsfwmodels

import "testing"

func TestChatState(t *testing.T) {
	s := &chatState{}

	if s.IsChanged() {
		t.Fatal("expected IsChanged to be false initially")
	}
	if s.GetAwaitingReplyTo() != "" {
		t.Fatal("expected empty AwaitingReplyTo")
	}

	// SetAwaitingReplyTo trims leading slash
	s.SetAwaitingReplyTo("/step1")
	if s.GetAwaitingReplyTo() != "step1" {
		t.Fatalf("expected 'step1', got %q", s.GetAwaitingReplyTo())
	}
	if !s.IsChanged() {
		t.Fatal("expected IsChanged to be true after SetAwaitingReplyTo")
	}

	// IsAwaitingReplyTo without query
	if !s.IsAwaitingReplyTo("step1") {
		t.Fatal("expected IsAwaitingReplyTo('step1') to be true")
	}
	if s.IsAwaitingReplyTo("step2") {
		t.Fatal("expected IsAwaitingReplyTo('step2') to be false")
	}

	// PushStepToAwaitingReplyTo without query
	s.PushStepToAwaitingReplyTo("step2")
	if s.GetAwaitingReplyTo() != "step1/step2" {
		t.Fatalf("expected 'step1/step2', got %q", s.GetAwaitingReplyTo())
	}
	// Push same step - no-op
	s.PushStepToAwaitingReplyTo("step2")
	if s.GetAwaitingReplyTo() != "step1/step2" {
		t.Fatalf("expected 'step1/step2', got %q", s.GetAwaitingReplyTo())
	}

	// Suffix check in IsAwaitingReplyTo
	if !s.IsAwaitingReplyTo("step2") {
		t.Fatal("expected IsAwaitingReplyTo('step2') to be true for suffix")
	}

	// With query
	s.SetAwaitingReplyTo("step1/step2?foo=bar")
	if !s.IsAwaitingReplyTo("step2") {
		t.Fatal("expected IsAwaitingReplyTo('step2') with query")
	}
	// PushStepToAwaitingReplyTo with query
	s.PushStepToAwaitingReplyTo("step3")
	if s.GetAwaitingReplyTo() != "step1/step2/step3?foo=bar" {
		t.Fatalf("expected 'step1/step2/step3?foo=bar', got %q", s.GetAwaitingReplyTo())
	}
	// Push same step with query - no-op
	s.PushStepToAwaitingReplyTo("step3")
	if s.GetAwaitingReplyTo() != "step1/step2/step3?foo=bar" {
		t.Fatalf("expected 'step1/step2/step3?foo=bar', got %q", s.GetAwaitingReplyTo())
	}

	// PopStepsFromAwaitingReplyUpToSpecificParent with query
	// step2 is not the last step
	s.PopStepsFromAwaitingReplyUpToSpecificParent("step2")
	if s.GetAwaitingReplyTo() != "step1/step2?foo=bar" {
		t.Fatalf("expected 'step1/step2?foo=bar', got %q", s.GetAwaitingReplyTo())
	}
	// step2 is now the last step - popping up to step2 is a no-op (i == len(steps)-1)
	s.PopStepsFromAwaitingReplyUpToSpecificParent("step2")
	if s.GetAwaitingReplyTo() != "step1/step2?foo=bar" {
		t.Fatalf("expected 'step1/step2?foo=bar', got %q", s.GetAwaitingReplyTo())
	}

	// PopStepsFromAwaitingReplyUpToSpecificParent without query
	s.SetAwaitingReplyTo("a/b/c")
	s.PopStepsFromAwaitingReplyUpToSpecificParent("b")
	if s.GetAwaitingReplyTo() != "a/b" {
		t.Fatalf("expected 'a/b', got %q", s.GetAwaitingReplyTo())
	}
	// step not found in steps
	s.PopStepsFromAwaitingReplyUpToSpecificParent("unknown")
	if s.GetAwaitingReplyTo() != "a/b" {
		t.Fatalf("expected 'a/b', got %q", s.GetAwaitingReplyTo())
	}

	// AddWizardParam and GetWizardParam
	s.SetAwaitingReplyTo("step1")
	s.AddWizardParam("k1", "v1")
	if s.GetWizardParam("k1") != "v1" {
		t.Fatalf("expected wizard param 'v1', got %q", s.GetWizardParam("k1"))
	}
	s.AddWizardParam("k2", "v2")
	if s.GetWizardParam("k2") != "v2" {
		t.Fatalf("expected wizard param 'v2', got %q", s.GetWizardParam("k2"))
	}

	// GetWizardParam on invalid url
	s.AwaitingReplyTo = "http://[::1]:namedport"
	if s.GetWizardParam("k1") != "" {
		t.Fatal("expected empty string from GetWizardParam on invalid url")
	}

	// AddWizardParam panic on invalid url
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic from AddWizardParam on invalid url")
		}
	}()
	s.AddWizardParam("k3", "v3")
}
