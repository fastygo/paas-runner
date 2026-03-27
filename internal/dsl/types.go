package dsl

type Input struct {
	Name     string   `yaml:"name"`
	Label    string   `yaml:"label,omitempty"`
	Type     string   `yaml:"type"`
	Default  *string  `yaml:"default,omitempty"`
	Required bool     `yaml:"required,omitempty"`
	Options  []string `yaml:"options,omitempty"`
}

type Step struct {
	ID          string            `yaml:"id,omitempty"`
	Run         string            `yaml:"run"`
	Description string            `yaml:"description,omitempty"`
	Local       bool              `yaml:"local,omitempty"`
	Capture     string            `yaml:"capture,omitempty"`
	When        string            `yaml:"when,omitempty"`
	IgnoreError bool              `yaml:"ignore_error,omitempty"`
	Workdir     string            `yaml:"workdir,omitempty"`
	Env         map[string]string `yaml:"env,omitempty"`
}

type Extension struct {
	ID          string  `yaml:"id"`
	Name        string  `yaml:"name"`
	Description string  `yaml:"description,omitempty"`
	Inputs      []Input `yaml:"inputs,omitempty"`
	Steps       []Step  `yaml:"steps"`
}
