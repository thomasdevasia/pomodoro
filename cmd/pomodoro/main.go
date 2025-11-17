package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/thomasdevasia/pomodoro/internal/controller"
)

func main() {
	p := tea.NewProgram(controller.New(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}
