package dsl

import "testing"

func TestSubstituteVariables(t *testing.T) {
	env := map[string]string{
		"INPUT_NAME": "app",
		"EMPTY":      "",
	}

	value, err := SubstituteVariables("echo ${INPUT_NAME}", env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if value != "echo app" {
		t.Fatalf("unexpected substitution result %q", value)
	}

	if _, err := SubstituteVariables("echo ${MISSING}", env); err == nil {
		t.Fatal("expected missing variable error")
	}

	if _, err := SubstituteVariables("echo", map[string]string{}); err != nil {
		t.Fatalf("did not expect error for no substitutions: %v", err)
	}

	value, err = SubstituteVariables("echo ${EMPTY}", env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if value != "echo " {
		t.Fatalf("unexpected empty substitution result %q", value)
	}
}
