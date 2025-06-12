package main

import (
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"ki/internal/ui/app"
)


func main() {
	// Check if kind is installed
	if _, err := exec.LookPath("kind"); err != nil {
		fmt.Fprintf(os.Stderr, "Error: 'kind' command not found. Please install KIND first.\n")
		fmt.Fprintf(os.Stderr, "Visit: https://kind.sigs.k8s.io/docs/user/quick-start/#installation\n")
		os.Exit(1)
	}

	p := tea.NewProgram(app.NewApp(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
