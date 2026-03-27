package dsl

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

var allowedInputTypes = map[string]struct{}{
	"text":     {},
	"select":   {},
	"confirm":  {},
	"password": {},
}

func ValidateExtension(ext Extension) error {
	if ext.ID == "" {
		return ValidationError{Message: "extension id is required"}
	}

	if len(ext.Steps) == 0 {
		return ValidationError{Message: fmt.Sprintf("extension %q has no steps", ext.ID)}
	}

	seenSteps := map[string]struct{}{}
	for index, step := range ext.Steps {
		if err := ValidateStep(ext, index, step, seenSteps); err != nil {
			return err
		}
	}

	for _, input := range ext.Inputs {
		if err := ValidateInput(input); err != nil {
			return err
		}
	}

	return nil
}

func ValidateStep(ext Extension, index int, step Step, seen map[string]struct{}) error {
	if step.Run == "" {
		return ValidationError{Message: fmt.Sprintf("extension %q step %d: run is required", ext.ID, index+1)}
	}

	if step.Capture != "" {
		if !isVariableName("STEP_" + sanitizeIdentifier(step.Capture)) {
			return ValidationError{Message: fmt.Sprintf("extension %q step %d: invalid capture %q", ext.ID, index+1, step.Capture)}
		}
	}

	if _, err := ParseWhen(step.When); err != nil {
		return ValidationError{Message: fmt.Sprintf("extension %q step %d: %v", ext.ID, index+1, err)}
	}

	for _, key := range ExtractVariableRefs(step.Run) {
		if !isVariableName(key) {
			return ValidationError{Message: fmt.Sprintf("extension %q step %d: invalid variable %q", ext.ID, index+1, key)}
		}
	}

	for _, value := range step.Env {
		for _, key := range ExtractVariableRefs(value) {
			if !isVariableName(key) {
				return ValidationError{Message: fmt.Sprintf("extension %q step %d: invalid variable %q", ext.ID, index+1, key)}
			}
		}
	}

	if step.ID != "" {
		if _, ok := seen[step.ID]; ok {
			return ValidationError{Message: fmt.Sprintf("extension %q has duplicate step id %q", ext.ID, step.ID)}
		}
		seen[step.ID] = struct{}{}
	}

	return nil
}

func ValidateInput(input Input) error {
	if input.Name == "" {
		return ValidationError{Message: "input name is required"}
	}

	envName := NormalizeInputEnvKey(input.Name)
	if !isVariableName(envName) {
		return ValidationError{Message: fmt.Sprintf("invalid input name %q", input.Name)}
	}

	if _, ok := allowedInputTypes[input.Type]; !ok {
		return ValidationError{Message: fmt.Sprintf("input %q has unsupported type %q", input.Name, input.Type)}
	}

	if input.Type == "select" && len(input.Options) == 0 {
		return ValidationError{Message: fmt.Sprintf("input %q is select but has no options", input.Name)}
	}

	return nil
}

func ExtractVariableRefs(raw string) []string {
	refs := templateVarPattern.FindAllStringSubmatch(raw, -1)
	uniq := make(map[string]struct{})
	out := make([]string, 0)

	for _, m := range refs {
		if len(m) < 2 {
			continue
		}

		if _, ok := uniq[m[1]]; ok {
			continue
		}

		uniq[m[1]] = struct{}{}
		out = append(out, m[1])
	}

	sort.Strings(out)

	return out
}

func BuildBaseEnv(projectDefaults map[string]string, serverEnv map[string]string, processEnv map[string]string, inputs []Input, cliInputs map[string]string, stepCaptures map[string]string) (map[string]string, error) {
	env := make(map[string]string)

	for key, value := range projectDefaults {
		env[key] = value
	}

	for key, value := range serverEnv {
		env[key] = value
	}

	for key, value := range processEnv {
		env[key] = value
	}

	for _, input := range inputs {
		key := NormalizeInputEnvKey(input.Name)

		if input.Default != nil {
			env[key] = *input.Default
		}

		if value, ok := processEnv[key]; ok {
			env[key] = value
		}

		if value, ok := cliInputs[key]; ok {
			env[key] = value
		}

		if input.Required && env[key] == "" {
			return nil, &MissingRequiredInputError{Input: input.Name}
		}
	}

	for key, value := range stepCaptures {
		env[key] = value
	}

	return env, nil
}

func sanitizeIdentifier(raw string) string {
	if raw == "" {
		return ""
	}

	out := make([]rune, 0, len(raw))
	for _, ch := range raw {
		if ch >= 'a' && ch <= 'z' {
			ch = ch - ('a' - 'A')
		}

		if ch == '-' {
			ch = '_'
		}

		if (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_' {
			out = append(out, ch)
		}
	}

	return string(out)
}

func InputDefaults(inputs []Input) map[string]string {
	env := make(map[string]string)

	for _, input := range inputs {
		key := NormalizeInputEnvKey(input.Name)
		if input.Default != nil {
			env[key] = *input.Default
		}
	}

	return env
}

func ProcessEnvironment() map[string]string {
	env := make(map[string]string)

	for _, item := range os.Environ() {
		if parts := strings.SplitN(item, "=", 2); len(parts) == 2 {
			env[parts[0]] = parts[1]
		}
	}

	return env
}

func BuildInputEnvFromCLI(values map[string]string) map[string]string {
	env := make(map[string]string)

	for key, value := range values {
		env[NormalizeInputEnvKey(key)] = value
	}

	return env
}

func NormalizeInputEnvKey(raw string) string {
	return "INPUT_" + sanitizeIdentifier(raw)
}

type MissingRequiredInputError struct {
	Input string
}

func (e *MissingRequiredInputError) Error() string {
	return fmt.Sprintf("missing required input %q", e.Input)
}
