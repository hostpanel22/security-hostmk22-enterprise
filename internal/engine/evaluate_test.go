package engine

import "testing"

func TestEvaluateRuleNumericAndSet(t *testing.T) {
	attrs := map[string]any{
		"distance_km": 720.0,
		"country":     "US",
	}

	matched, err := EvaluateRule(attrs, "distance_km", "gt", "500")
	if err != nil || !matched {
		t.Fatalf("expected distance match, err=%v matched=%v", err, matched)
	}

	matched, err = EvaluateRule(attrs, "country", "in", "FR,DE")
	if err != nil {
		t.Fatal(err)
	}
	if matched {
		t.Fatal("US should not match FR,DE")
	}

	matched, err = EvaluateRule(attrs, "country", "neq", "FR")
	if err != nil || !matched {
		t.Fatalf("expected neq match, err=%v matched=%v", err, matched)
	}
}
