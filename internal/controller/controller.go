package controller

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Title       string
	Description string
	Duration    int
	Timer       timer.Model
	Step        Step
	form        *huh.Form
	help        help.Model
	keys        KeyMap
	progress    progress.Model
}

type KeyMap struct {
	Quit key.Binding
}

var Keys = KeyMap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}

// ShortHelp returns a short set of keybindings shown in the mini help view.
// It implements the `help.KeyMap` (bubbles) interface.
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit}
}

// FullHelp returns keybindings for the expanded help view.
// It also implements the `help.KeyMap` interface.
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Quit}}
}

// Step represents the current input step in the UI.
type Step int

const (
	StepForm Step = iota
	StepFormDone
)

func createDurationOptions(number int) []huh.Option[int] {
	options := make([]huh.Option[int], number)
	for i := 1; i <= number; i++ {
		label := fmt.Sprintf("%d minutes", i)
		value := i
		options[i-1] = huh.NewOption(label, value)
	}
	return options
}

func New() *Model {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter Title for Pomodoro").
				Placeholder("Like 'Study Session'").CharLimit(50).Key("Title"),
			huh.NewText().Title("Give Description").CharLimit(200).Key("Description"),
			huh.NewSelect[int]().
				Title("Select Duration (minutes)").
				OptionsFunc(
					func() []huh.Option[int] {
						return createDurationOptions(60)
					},
					"15",
				).
				Height(5).
				Key("DurationSelect"),
		),
	)

	return &Model{
		Step: StepForm,
		form: form,
		help: help.New(),
		keys: Keys,
	}
}

func (m Model) Init() tea.Cmd {
	return m.form.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// ...

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	// When the form completes, assign values to the model and advance the step.
	if m.Step == StepForm && m.form.State == huh.StateCompleted {
		title := strings.TrimSpace(m.form.GetString("Title"))
		if title == "" {
			title = "Focus Time"
		}
		m.Title = title
		m.Description = m.form.GetString("Description")
		m.Duration = m.form.GetInt("DurationSelect")
		m.Step = StepFormDone
		// Create timer using the selected duration (minutes).
		// Convert minutes to a time.Duration so the timer reflects the form selection.
		timeout := time.Duration(m.Duration) * time.Minute
		// timeout := time.Second * 5
		m.Timer = timer.NewWithInterval(timeout, time.Second)
		// Initialize/start the timer and merge its command with the form update command.
		cmd = tea.Batch(cmd, m.Timer.Init())
		p := progress.New(progress.WithDefaultGradient())
		// p.Width = 40
		m.progress = p
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		}
	case timer.TickMsg:
		var cmd tea.Cmd
		m.Timer, cmd = m.Timer.Update(msg)
		// Update progress bar based on elapsed time.
		if m.Timer.Timeout > 0 {
			minutesLeft := m.Timer.Timeout.Minutes()
			minutesTotal := (time.Duration(m.Duration) * time.Minute).Minutes()
			percent := 1 - (minutesLeft / minutesTotal)
			cmd = tea.Batch(cmd, m.progress.SetPercent(percent))
		} else {
			cmd = tea.Batch(cmd, m.progress.SetPercent(1.0))
			// cmd = tea.Batch(cmd, tea.Quit)
		}
		return m, cmd
	case progress.FrameMsg:
		// Let the progress model handle frame/animation messages.
		pm, cmd := m.progress.Update(msg)
		m.progress = pm.(progress.Model)
		return m, cmd
	case timer.StartStopMsg:
		var cmd tea.Cmd
		m.Timer, cmd = m.Timer.Update(msg)
		// Enable Quit only when the timer is NOT running.
		m.keys.Quit.SetEnabled(!m.Timer.Running())
		return m, cmd
	case timer.TimeoutMsg:
		// m.quitting = true
		// use notify-send or similar to send a notification here
		cmd = tea.Batch(cmd, tea.ExecProcess(exec.Command("spd-say", fmt.Sprintf("Hey, The timer has ran out for %s", m.Title)), nil))
		cmd = tea.Batch(cmd, tea.ExecProcess(exec.Command("notify-send", "-i", "dialog-information", fmt.Sprintf("%s", m.Title), "Timer has Finished"), nil))
		// cmd = tea.Batch(cmd, tea.Quit)
		return m, cmd
	case tea.QuitMsg:
		return m, cmd
	}
	return m, cmd
}

var (
	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7D56F4")).Padding(2, 12)
	middleStyle = lipgloss.NewStyle().Padding(1, 3)
)

func (m Model) View() string {
	helpView := m.help.View(m.keys)
	if m.form.State == huh.StateCompleted {
		return titleStyle.Render(m.Title) + middleStyle.Render("\nTimer Started!\n") + middleStyle.Render(m.Timer.View()) + "\n\n" + middleStyle.Render(m.progress.View()) + "\n\n" + middleStyle.Render(helpView)
	}
	return m.form.View() + "\n" + helpView
}
