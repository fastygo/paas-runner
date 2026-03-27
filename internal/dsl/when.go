package dsl

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	whenAny      = "any"
	whenVar      = "var"
	whenNot      = "not"
	whenEquals   = "eq"
	whenNotEqual = "neq"
)

type WhenCondition struct {
	Kind     string
	Var      string
	Literal  string
	Operator string
}

func ParseWhen(raw string) (WhenCondition, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return WhenCondition{Kind: whenAny}, nil
	}

	if strings.HasPrefix(raw, "not ") {
		rest := strings.TrimSpace(strings.TrimPrefix(raw, "not"))
		if rest == "" {
			return WhenCondition{}, fmt.Errorf("invalid when expression: %q", raw)
		}
		if !isVariableName(rest) {
			return WhenCondition{}, fmt.Errorf("invalid variable in when expression: %q", raw)
		}
		return WhenCondition{Kind: whenNot, Var: rest}, nil
	}

	if strings.Contains(raw, "!=") {
		parts := strings.SplitN(raw, "!=", 2)
		if len(parts) != 2 {
			return WhenCondition{}, fmt.Errorf("invalid when expression: %q", raw)
		}

		variable := strings.TrimSpace(parts[0])
		literal := strings.TrimSpace(parts[1])

		if !isVariableName(variable) {
			return WhenCondition{}, fmt.Errorf("invalid variable in when expression: %q", raw)
		}

		value, err := parseQuotedLiteral(literal)
		if err != nil {
			return WhenCondition{}, fmt.Errorf("invalid when expression: %w", err)
		}

		return WhenCondition{Kind: whenNotEqual, Var: variable, Literal: value, Operator: "!="}, nil
	}

	if strings.Contains(raw, "==") {
		parts := strings.SplitN(raw, "==", 2)
		if len(parts) != 2 {
			return WhenCondition{}, fmt.Errorf("invalid when expression: %q", raw)
		}

		variable := strings.TrimSpace(parts[0])
		literal := strings.TrimSpace(parts[1])

		if !isVariableName(variable) {
			return WhenCondition{}, fmt.Errorf("invalid variable in when expression: %q", raw)
		}

		value, err := parseQuotedLiteral(literal)
		if err != nil {
			return WhenCondition{}, fmt.Errorf("invalid when expression: %w", err)
		}

		return WhenCondition{Kind: whenEquals, Var: variable, Literal: value, Operator: "=="}, nil
	}

	if isVariableName(raw) {
		return WhenCondition{Kind: whenVar, Var: raw}, nil
	}

	return WhenCondition{}, fmt.Errorf("invalid when expression: %q", raw)
}

func parseQuotedLiteral(raw string) (string, error) {
	if !strings.HasPrefix(raw, "\"") || !strings.HasSuffix(raw, "\"") || len(raw) < 2 {
		return "", fmt.Errorf("literal must be double quoted")
	}

	value, err := strconv.Unquote(raw)
	if err != nil {
		return "", err
	}

	return value, nil
}

func EvaluateWhen(raw string, env map[string]string) (bool, error) {
	cond, err := ParseWhen(raw)
	if err != nil {
		return false, err
	}

	return evalCondition(cond, env), nil
}

func evalCondition(cond WhenCondition, env map[string]string) bool {
	value := env[cond.Var]

	switch cond.Kind {
	case whenAny:
		return true
	case whenVar:
		return isTruthy(value)
	case whenNot:
		return !isTruthy(value)
	case whenEquals:
		return value == cond.Literal
	case whenNotEqual:
		return value != cond.Literal
	default:
		return false
	}
}

func isTruthy(value string) bool {
	return value != "" && value != "0" && value != "false"
}

func isVariableName(raw string) bool {
	if raw == "" {
		return false
	}

	if !(raw[0] == '_' || (raw[0] >= 'A' && raw[0] <= 'Z')) {
		return false
	}

	for i := 1; i < len(raw); i++ {
		ch := raw[i]
		if ch == '_' || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') {
			continue
		}
		return false
	}

	return true
}
