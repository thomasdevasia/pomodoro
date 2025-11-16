package controller

import (
	"fmt"

	"github.com/thomasdevasia/pomodoro/internal/model"
	"github.com/thomasdevasia/pomodoro/internal/view"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Controller struct {
	m *model.Model
}

func New() *Controller {
	return &Controller{m: model.New()}
}

// Init is part of tea.Model
func (c *Controller) Init() tea.Cmd {
	return textinput.Blink
}

// Update is part of tea.Model
func (c *Controller) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// If we're on the title step, set title and move to description step
			if c.m.Step == model.StepTitle {
				c.m.SetTitleFromInput("Focus Time")
				c.m.InputTitle.Blur()
				c.m.InputDescription.Focus()
				c.m.Step = model.StepDescription
				return c, nil
			}

			// If we're on the description step, set description and finish
			if c.m.Step == model.StepDescription {
				c.m.SetDescriptionFromInput("No description")
				c.m.Step = model.StepDone
				// Ensure output starts on a fresh line (avoid leftover padding/cursor)
				fmt.Printf("\nPomodoro Title: %s\nDescription: %s\n", c.m.Title, c.m.Description)
				return c, tea.Quit
			}
		}
	}

	// Update both inputs (only the focused one will react)
	var cmd1, cmd2 tea.Cmd
	c.m.InputTitle, cmd1 = c.m.InputTitle.Update(msg)
	c.m.InputDescription, cmd2 = c.m.InputDescription.Update(msg)
	return c, tea.Batch(cmd1, cmd2)
}

// View is part of tea.Model
func (c *Controller) View() string {
	return view.Render(c.m)
}
