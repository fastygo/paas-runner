package dsl

import (
	"fmt"
	"regexp"
	"strings"
)

var templateVarPattern = regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)\}`)

func SubstituteVariables(raw string, env map[string]string) (string, error) {
	missing := make([]string, 0)

	replaced := templateVarPattern.ReplaceAllStringFunc(raw, func(match string) string {
		groups := templateVarPattern.FindStringSubmatch(match)
		if len(groups) != 2 {
			return match
		}

		key := groups[1]
		value, ok := env[key]
		if !ok {
			missing = appendUnique(missing, key)
			return match
		}

		return value
	})

	if len(missing) > 0 {
		return "", &UndefinedVariableError{Template: raw, Variables: missing}
	}

	return replaced, nil
}

func substituteAll(raw map[string]string, env map[string]string) (map[string]string, error) {
	result := make(map[string]string, len(raw))
	for key, value := range raw {
		replaced, err := SubstituteVariables(value, env)
		if err != nil {
			return nil, err
		}
		result[key] = replaced
	}

	return result, nil
}

func extractRefs(raw string) []string {
	refs := templateVarPattern.FindAllStringSubmatch(raw, -1)
	seen := make(map[string]struct{})
	ordered := make([]string, 0)

	for _, entry := range refs {
		if len(entry) < 2 {
			continue
		}

		if _, ok := seen[entry[1]]; ok {
			continue
		}

		seen[entry[1]] = struct{}{}
		ordered = append(ordered, entry[1])
	}

	return ordered
}

func appendUnique(values []string, value string) []string {
	for _, v := range values {
		if v == value {
			return values
		}
	}

	return append(values, value)
}

func MaskableRefs(raw string) string {
	return strings.Join(extractRefs(raw), ",")
}

type UndefinedVariableError struct {
	Template  string
	Variables []string
}

func (e *UndefinedVariableError) Error() string {
	return fmt.Sprintf("undefined variable(s) in %q: %v", e.Template, e.Variables)
}
