package botsfwmodels

import "testing"

func TestChatVars(t *testing.T) {
	v := &chatVars{}

	// When Vars is nil
	if v.GetVar("k1") != "" {
		t.Fatal("expected empty string when Vars is nil")
	}
	if v.HasChangedVars() {
		t.Fatal("expected no changed vars when Vars is nil")
	}
	v.DelVar("k1") // should be no-op when nil

	// SetVar on nil map
	v.SetVar("k1", "v1")
	if v.GetVar("k1") != "v1" {
		t.Fatalf("expected 'v1', got %q", v.GetVar("k1"))
	}
	if !v.HasChangedVars() {
		t.Fatal("expected HasChangedVars to be true after SetVar")
	}

	// SetVar same value - no-op
	v.SetVar("k1", "v1")

	// SetVar different value
	v.SetVar("k2", "v2")
	if v.GetVar("k2") != "v2" {
		t.Fatalf("expected 'v2', got %q", v.GetVar("k2"))
	}

	// DelVar non-existent key
	v.DelVar("non-existent")

	// DelVar existing key (k1 is in changed)
	v.DelVar("k1")
	if v.GetVar("k1") != "" {
		t.Fatal("expected k1 to be deleted")
	}
	if !v.HasChangedVars() {
		t.Fatal("expected HasChangedVars to be true with deleted keys")
	}

	// DelVar k1 again - should be no-op
	v.DelVar("k1")

	// SetVar on a key that was in deleted
	v.SetVar("k1", "v1_restored")
	if v.GetVar("k1") != "v1_restored" {
		t.Fatalf("expected restored value, got %q", v.GetVar("k1"))
	}
}
