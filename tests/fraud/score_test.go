package fraud_test

import (
	"testing"

	"github.com/azezaz/HostMK-Security-Framework/core/fraud"
)

func TestNewScoreLow(t *testing.T) {
	score := fraud.NewScore(10)

	if score.Value != 10 {
		t.Fatalf("expected 10 got %d", score.Value)
	}

	if score.Level != fraud.RiskLow {
		t.Fatalf("expected %s got %s", fraud.RiskLow, score.Level)
	}
}

func TestNewScoreMedium(t *testing.T) {
	score := fraud.NewScore(40)

	if score.Level != fraud.RiskMedium {
		t.Fatalf("expected %s got %s", fraud.RiskMedium, score.Level)
	}
}

func TestNewScoreHigh(t *testing.T) {
	score := fraud.NewScore(70)

	if score.Level != fraud.RiskHigh {
		t.Fatalf("expected %s got %s", fraud.RiskHigh, score.Level)
	}
}

func TestNewScoreCritical(t *testing.T) {
	score := fraud.NewScore(90)

	if score.Level != fraud.RiskCritical {
		t.Fatalf("expected %s got %s", fraud.RiskCritical, score.Level)
	}
}

func TestCalculateLevelLow(t *testing.T) {
	score := &fraud.Score{
		Value: 20,
	}

	score.CalculateLevel()

	if score.Level != fraud.RiskLow {
		t.Fatalf("expected %s got %s", fraud.RiskLow, score.Level)
	}
}

func TestCalculateLevelMedium(t *testing.T) {
	score := &fraud.Score{
		Value: 50,
	}

	score.CalculateLevel()

	if score.Level != fraud.RiskMedium {
		t.Fatalf("expected %s got %s", fraud.RiskMedium, score.Level)
	}
}

func TestCalculateLevelHigh(t *testing.T) {
	score := &fraud.Score{
		Value: 80,
	}

	score.CalculateLevel()

	if score.Level != fraud.RiskHigh {
		t.Fatalf("expected %s got %s", fraud.RiskHigh, score.Level)
	}
}

func TestCalculateLevelCritical(t *testing.T) {
	score := &fraud.Score{
		Value: 95,
	}

	score.CalculateLevel()

	if score.Level != fraud.RiskCritical {
		t.Fatalf("expected %s got %s", fraud.RiskCritical, score.Level)
	}
}

func TestAddKeepsLowLevel(t *testing.T) {
	score := fraud.NewScore(10)

	score.Add(20)

	if score.Value != 30 {
		t.Fatalf("expected 30 got %d", score.Value)
	}

	if score.Level != fraud.RiskLow {
		t.Fatalf("expected %s got %s", fraud.RiskLow, score.Level)
	}
}

func TestAddChangesToMedium(t *testing.T) {
	score := fraud.NewScore(20)

	score.Add(25)

	if score.Level != fraud.RiskMedium {
		t.Fatalf("expected %s got %s", fraud.RiskMedium, score.Level)
	}
}

func TestAddChangesToHigh(t *testing.T) {
	score := fraud.NewScore(50)

	score.Add(25)

	if score.Level != fraud.RiskHigh {
		t.Fatalf("expected %s got %s", fraud.RiskHigh, score.Level)
	}
}

func TestAddChangesToCritical(t *testing.T) {
	score := fraud.NewScore(80)

	score.Add(20)

	if score.Level != fraud.RiskCritical {
		t.Fatalf("expected %s got %s", fraud.RiskCritical, score.Level)
	}
}

func TestIsLow(t *testing.T) {
	score := fraud.NewScore(10)

	if !score.IsLow() {
		t.Fatal("expected low")
	}
}

func TestIsMedium(t *testing.T) {
	score := fraud.NewScore(50)

	if !score.IsMedium() {
		t.Fatal("expected medium")
	}
}

func TestIsHigh(t *testing.T) {
	score := fraud.NewScore(80)

	if !score.IsHigh() {
		t.Fatal("expected high")
	}
}

func TestIsCritical(t *testing.T) {
	score := fraud.NewScore(100)

	if !score.IsCritical() {
		t.Fatal("expected critical")
	}
}

func TestRequiresActionLow(t *testing.T) {
	score := fraud.NewScore(10)

	if score.RequiresAction() {
		t.Fatal("low risk should not require action")
	}
}

func TestRequiresActionMedium(t *testing.T) {
	score := fraud.NewScore(50)

	if score.RequiresAction() {
		t.Fatal("medium risk should not require action")
	}
}

func TestRequiresActionHigh(t *testing.T) {
	score := fraud.NewScore(75)

	if !score.RequiresAction() {
		t.Fatal("high risk should require action")
	}
}

func TestRequiresActionCritical(t *testing.T) {
	score := fraud.NewScore(100)

	if !score.RequiresAction() {
		t.Fatal("critical risk should require action")
	}
}

func TestScoreThresholdBoundaries(t *testing.T) {
	cases := []struct {
		value    int
		expected fraud.RiskLevel
	}{
		{0, fraud.RiskLow},
		{39, fraud.RiskLow},
		{40, fraud.RiskMedium},
		{69, fraud.RiskMedium},
		{70, fraud.RiskHigh},
		{89, fraud.RiskHigh},
		{90, fraud.RiskCritical},
		{150, fraud.RiskCritical},
	}

	for _, tc := range cases {
		score := fraud.NewScore(tc.value)

		if score.Level != tc.expected {
			t.Fatalf(
				"value=%d expected=%s got=%s",
				tc.value,
				tc.expected,
				score.Level,
			)
		}
	}
}
