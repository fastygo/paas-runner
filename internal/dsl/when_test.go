package dsl

import "testing"

func TestParseWhen(t *testing.T) {
	tests := []struct {
		input string
		err   bool
		kind  string
	}{
		{"", false, whenAny},
		{"INPUT_TAG", false, whenVar},
		{"not INPUT_SKIP", false, whenNot},
		{"STEP_COUNT == \"2\"", false, whenEquals},
		{"STEP_COUNT != \"\"", false, whenNotEqual},
		{"docker ps", true, ""},
	}

	for _, tc := range tests {
		_, err := ParseWhen(tc.input)
		if (err != nil) != tc.err {
			t.Fatalf("parse when %q err=%v expected error=%v", tc.input, err, tc.err)
		}

		if err == nil {
			parsed, _ := ParseWhen(tc.input)
			if parsed.Kind != tc.kind {
				t.Fatalf("parse when %q kind=%q expected %q", tc.input, parsed.Kind, tc.kind)
			}
		}
	}
}

func TestEvaluateWhen(t *testing.T) {
	env := map[string]string{
		"VAR":   "1",
		"FALSE": "false",
		"ZERO":  "0",
		"EMPTY": "",
	}

	checks := []struct {
		expr string
		want bool
	}{
		{"", true},
		{"VAR", true},
		{"not VAR", false},
		{"FALSE", false},
		{"ZERO", false},
		{"EMPTY", false},
		{"VAR == \"1\"", true},
		{"VAR == \"2\"", false},
		{"EMPTY != \"\"", false},
	}

	for _, tc := range checks {
		got, err := EvaluateWhen(tc.expr, env)
		if err != nil {
			t.Fatalf("evaluate when %q failed: %v", tc.expr, err)
		}
		if got != tc.want {
			t.Fatalf("evaluate when %q got=%v want=%v", tc.expr, got, tc.want)
		}
	}
}
