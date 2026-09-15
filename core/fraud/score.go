package fraud

type RiskLevel string

const (
	RiskLow RiskLevel = "low"

	RiskMedium RiskLevel = "medium"

	RiskHigh RiskLevel = "high"

	RiskCritical RiskLevel = "critical"
)

type Score struct {
	Value int `json:"value"`

	Level RiskLevel `json:"level"`
}

func NewScore(
	value int,
) *Score {

	score := &Score{
		Value: value,
	}

	score.CalculateLevel()

	return score
}

func (s *Score) CalculateLevel() {

	switch {

	case s.Value >= 90:
		s.Level = RiskCritical

	case s.Value >= 70:
		s.Level = RiskHigh

	case s.Value >= 40:
		s.Level = RiskMedium

	default:
		s.Level = RiskLow
	}
}

func (s *Score) Add(
	weight int,
) {

	s.Value += weight

	s.CalculateLevel()
}

func (s *Score) IsLow() bool {

	return s.Level == RiskLow
}

func (s *Score) IsMedium() bool {

	return s.Level == RiskMedium
}

func (s *Score) IsHigh() bool {

	return s.Level == RiskHigh
}

func (s *Score) IsCritical() bool {

	return s.Level == RiskCritical
}

func (s *Score) RequiresAction() bool {

	return s.Level == RiskHigh ||
		s.Level == RiskCritical
}
