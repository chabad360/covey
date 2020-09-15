package ui

type InputType int

const (
	// Text is a text input.
	Text InputType = iota

	// Number is a integer input.
	Number

	// List is a single select input.
	List

	// Multiple is a multi-select input.
	Multiple

	// Boolean is a checkbox input.
	Boolean
)

// Form is the form that is used for plugins to specify their options.
type Form struct {
	Inputs []Input `json:"inputs"`
}

// Input is a Form input. The response is tested by go-playground/validate to ensure the input is valid.
type Input struct {
	Name        string    `json:"name"`
	Type        InputType `json:"type"`
	Default     string    `json:"default,omitempty"`
	Label       string    `json:"label"`
	Options     []string  `json:"options,omitempty"`
	Required    bool      `json:"required"`
	Description string    `json:"description,omitempty"`
}
