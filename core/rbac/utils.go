package rbac

import (
	"strings"
)

func normalize(
	value string,
) string {

	return strings.ToLower(
		strings.TrimSpace(
			value,
		),
	)
}

func containsString(
	values []string,
	target string,
) bool {

	target = normalize(
		target,
	)

	for _, value := range values {

		if normalize(
			value,
		) == target {

			return true
		}
	}

	return false
}

func removeString(
	values []string,
	target string,
) []string {

	target = normalize(
		target,
	)

	result := make(
		[]string,
		0,
		len(values),
	)

	for _, value := range values {

		if normalize(
			value,
		) == target {

			continue
		}

		result = append(
			result,
			value,
		)
	}

	return result
}
