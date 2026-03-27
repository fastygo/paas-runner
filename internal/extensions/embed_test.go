package extensions

import (
	"strings"
	"testing"

	"github.com/paas/paas-runner/internal/dsl"
)

func TestEmbeddedExtensionsValidate(t *testing.T) {
	names, err := List()
	if err != nil {
		t.Fatalf("List() failed: %v", err)
	}

	if len(names) == 0 {
		t.Fatal("expected embedded extensions to exist")
	}

	for _, name := range names {
		data, err := Read(name + ".yml")
		if err != nil {
			t.Fatalf("Read(%q) failed: %v", name, err)
		}

		extension, err := dsl.ParseExtension(data)
		if err != nil {
			t.Fatalf("ParseExtension(%q) failed: %v", name, err)
		}

		if strings.TrimSpace(extension.ID) == "" {
			t.Fatalf("extension %q has empty id", name)
		}

		if err := dsl.ValidateExtension(extension); err != nil {
			t.Fatalf("ValidateExtension(%q) failed: %v", name, err)
		}
	}
}
