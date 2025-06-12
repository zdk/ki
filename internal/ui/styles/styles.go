package styles

import "github.com/charmbracelet/lipgloss"

var (
	Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFDF5")).
		Background(lipgloss.Color("#25A065")).
		Padding(0, 1)

	Status = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#04B575")).
		Bold(true)

	Error = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF5F87")).
		Bold(true)

	Help = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262"))

	Focused = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFF")).
		Background(lipgloss.Color("#7C56F4")).
		Bold(true)

	Blurred = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262"))

	NoStyle = lipgloss.NewStyle()
)