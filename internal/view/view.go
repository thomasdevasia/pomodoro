package view

import (
	"fmt"
	"strings"

	"github.com/thomasdevasia/pomodoro/internal/model"
)

func Render(m *model.Model) string {
	if m.Step == model.StepTitle {
		return fmt.Sprintf("Give a Title to your Pomodoro:\n%s\n(press Enter to edit description)\n", m.InputTitle.View())
	}

	// Trim any trailing newlines that the textinput widget might include
	titleView := strings.TrimRight(m.InputTitle.View(), "\r\n")
	descView := strings.TrimRight(m.InputDescription.View(), "\r\n")

	return fmt.Sprintf("Give a Title to your Pomodoro:\n%s\nAdd a short Description (optional):\n%s\n(press Enter to submit)\n", titleView, descView)
}
