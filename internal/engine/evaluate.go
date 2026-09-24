package engine

import (
	"fmt"
	"strconv"
	"strings"
)

// EvaluateRule applies a generic comparison. Field names, operators, and
// thresholds are supplied at runtime from PostgreSQL — the engine has no
// catalog of business rules.
func EvaluateRule(attributes map[string]any, field, operator, threshold string) (bool, error) {
	value, ok := lookup(attributes, field)
	if !ok {
		return false, nil
	}

	op := strings.ToLower(strings.TrimSpace(operator))
	switch op {
	case "gt", "gte", "lt", "lte":
		left, err := toFloat(value)
		if err != nil {
			return false, err
		}
		right, err := strconv.ParseFloat(strings.TrimSpace(threshold), 64)
		if err != nil {
			return false, fmt.Errorf("invalid numeric threshold %q", threshold)
		}
		switch op {
		case "gt":
			return left > right, nil
		case "gte":
			return left >= right, nil
		case "lt":
			return left < right, nil
		case "lte":
			return left <= right, nil
		}
	case "eq":
		return stringify(value) == strings.TrimSpace(threshold), nil
	case "neq":
		return stringify(value) != strings.TrimSpace(threshold), nil
	case "in":
		return containsCSV(threshold, stringify(value)), nil
	case "not_in":
		return !containsCSV(threshold, stringify(value)), nil
	case "contains":
		return strings.Contains(strings.ToLower(stringify(value)), strings.ToLower(threshold)), nil
	default:
		return false, fmt.Errorf("unsupported operator %q", operator)
	}
	return false, nil
}

func lookup(attributes map[string]any, field string) (any, bool) {
	if attributes == nil {
		return nil, false
	}
	if value, ok := attributes[field]; ok {
		return value, true
	}
	for key, value := range attributes {
		if strings.EqualFold(key, field) {
			return value, true
		}
	}
	return nil, false
}

func toFloat(value any) (float64, error) {
	switch typed := value.(type) {
	case float64:
		return typed, nil
	case float32:
		return float64(typed), nil
	case int:
		return float64(typed), nil
	case int64:
		return float64(typed), nil
	case jsonNumber:
		return typed.Float64()
	default:
		return strconv.ParseFloat(stringify(value), 64)
	}
}

type jsonNumber interface {
	Float64() (float64, error)
}

func stringify(value any) string {
	return strings.TrimSpace(fmt.Sprint(value))
}

func containsCSV(csv, candidate string) bool {
	for _, part := range strings.Split(csv, ",") {
		if strings.EqualFold(strings.TrimSpace(part), candidate) {
			return true
		}
	}
	return false
}
