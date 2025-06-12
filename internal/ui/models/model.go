package models

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"ki/internal/cmd"
)

// Model represents the main application state
type Model struct {
	// Current view
	CurrentView ViewMode

	// Components
	MainMenu    list.Model
	ClusterList list.Model
	NodeList    list.Model
	TextInput   textinput.Model
	Help        help.Model

	// Data
	Clusters       []cmd.Cluster
	CurrentCluster *cmd.Cluster
	Message        string
	MessageType    string // "success", "error", "info"

	// State
	ShowHelp bool
	Quitting bool
	Loading  bool
	Width    int
	Height   int

	// Input context
	InputPrompt         string
	InputAction         string
	SelectedCluster     string
	ClusterToDelete     string
	DeleteConfirmChoice int // 0 = Yes, 1 = No
}

// Implement tea.Model interface
func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	return ""
}