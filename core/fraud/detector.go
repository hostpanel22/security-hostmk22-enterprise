package fraud

type DetectionResult struct {
	Score *Score `json:"score"`

	MatchedRules []*Rule `json:"matched_rules"`

	Violations []*Violation `json:"violations"`
}

type Detector struct {
	rules []*Rule
}

func NewDetector() *Detector {

	return &Detector{
		rules: make(
			[]*Rule,
			0,
		),
	}
}

func (d *Detector) AddRule(
	rule *Rule,
) error {

	if rule == nil {
		return ErrInvalidRule
	}

	if err := rule.Validate(); err != nil {
		return err
	}

	d.rules = append(
		d.rules,
		rule,
	)

	return nil
}

func (d *Detector) Rules() []*Rule {

	result := make(
		[]*Rule,
		0,
		len(d.rules),
	)

	for _, rule := range d.rules {

		result = append(
			result,
			rule.Clone(),
		)
	}

	return result
}

func (d *Detector) Detect(
	signal *Signal,
) (*DetectionResult, error) {

	if signal == nil {
		return nil, ErrInvalidSignal
	}

	if err := signal.Validate(); err != nil {
		return nil, err
	}

	result := &DetectionResult{
		Score: NewScore(0),

		MatchedRules: make(
			[]*Rule,
			0,
		),

		Violations: make(
			[]*Violation,
			0,
		),
	}

	for _, rule := range d.rules {

		if rule == nil {
			continue
		}

		if !rule.Match(
			signal,
		) {
			continue
		}

		result.MatchedRules = append(
			result.MatchedRules,
			rule.Clone(),
		)

		result.Score.Add(
			rule.Weight,
		)

		violation := NewViolation(
			signal.TenantID,
			signal.UserID,
			signal.ID,
			rule.ID,
			rule.Name,
			result.Score.Value,
			result.Score.Level,
		)

		result.Violations = append(
			result.Violations,
			violation,
		)
	}

	return result, nil
}

func (d *Detector) RuleCount() int {

	return len(
		d.rules,
	)
}

func (d *Detector) ClearRules() {

	d.rules = make(
		[]*Rule,
		0,
	)
}
