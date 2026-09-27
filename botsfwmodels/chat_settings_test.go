package botsfwmodels

import "testing"

func TestChatSettings(t *testing.T) {
	s := &chatSettings{}

	if s.GetPreferredLanguage() != "" {
		t.Fatalf("expected empty preferred language, got %q", s.GetPreferredLanguage())
	}

	s.SetPreferredLanguage("en")
	if s.GetPreferredLanguage() != "en" {
		t.Fatalf("expected 'en', got %q", s.GetPreferredLanguage())
	}

	// Add empty language code
	if s.AddClientLanguage("") {
		t.Fatal("expected false when adding empty language code")
	}

	// Add "root" language code
	if s.AddClientLanguage("root") {
		t.Fatal("expected false when adding 'root' language code")
	}

	// Add valid language code when PreferredLanguage is already set
	s.AddClientLanguage("fr")
	if len(s.LanguageCodes) != 1 || s.LanguageCodes[0] != "fr" {
		t.Fatalf("expected ['fr'], got %v", s.LanguageCodes)
	}

	// Add duplicate language code
	if s.AddClientLanguage("fr") {
		t.Fatal("expected false when adding duplicate language code")
	}

	// Add valid language code when PreferredLanguage is empty
	s2 := &chatSettings{}
	s2.AddClientLanguage("de")
	if s2.GetPreferredLanguage() != "de" {
		t.Fatalf("expected preferred language 'de', got %q", s2.GetPreferredLanguage())
	}
	if len(s2.LanguageCodes) != 1 || s2.LanguageCodes[0] != "de" {
		t.Fatalf("expected ['de'], got %v", s2.LanguageCodes)
	}
}
