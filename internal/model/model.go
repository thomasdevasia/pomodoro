package model

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
)

type Model struct {
	Title            string
	Description      string
	InputTitle       textinput.Model
	InputDescription textinput.Model
	// Step indicates which input step the UI is on (title, description, done)
	Step Step
}

// creates a initial model
func New() *Model {
	ti := textinput.New()
	ti.Placeholder = "Like Restudying for exams, Focus time"
	ti.Focus()
	ti.CharLimit = 64
	ti.Width = 50

	di := textinput.New()
	di.Placeholder = "A short note (optional)"
	di.CharLimit = 140
	di.Width = 50

	return &Model{
		Title:            "",
		Description:      "",
		InputTitle:       ti,
		InputDescription: di,
		Step:             StepTitle,
	}
}

// Step represents the current input step in the UI.
type Step int

const (
	StepTitle Step = iota
	StepDescription
	StepDone
)

// SetTitleFromInput sets Model.Title from the text input value. If the
// trimmed input is empty, it uses the provided defaultTitle instead.
func (m *Model) SetTitleFromInput(defaultTitle string) {
	v := strings.TrimSpace(m.InputTitle.Value())
	if v == "" {
		m.Title = defaultTitle
	} else {
		m.Title = v
	}
}

// SetDescriptionFromInput sets Model.Description from the description input value.
// If the trimmed input is empty, it uses the provided defaultDesc instead.
func (m *Model) SetDescriptionFromInput(defaultDesc string) {
	v := strings.TrimSpace(m.InputDescription.Value())
	if v == "" {
		m.Description = defaultDesc
	} else {
		m.Description = v
	}
}
