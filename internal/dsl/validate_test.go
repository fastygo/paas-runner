package dsl

import "testing"

func TestValidateWhenSyntax(t *testing.T) {
	ext := Extension{
		ID:    "test",
		Steps: []Step{{Run: "echo hi", When: "STEP_COUNT != \"0\""}},
	}

	if err := ValidateExtension(ext); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestValidateUnknownWhen(t *testing.T) {
	ext := Extension{
		ID:    "test",
		Steps: []Step{{Run: "echo hi", When: "docker ps"}},
	}

	if err := ValidateExtension(ext); err == nil {
		t.Fatal("expected validation error")
	}
}
