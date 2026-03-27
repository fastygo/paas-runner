package dsl

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func ParseExtension(raw []byte) (Extension, error) {
	var extension Extension
	if err := yaml.Unmarshal(raw, &extension); err != nil {
		return extension, fmt.Errorf("parse extension: %w", err)
	}

	return extension, nil
}
