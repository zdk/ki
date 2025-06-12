package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"ki/internal/cmd"
	"ki/internal/ui/commands"
	"ki/internal/ui/models"
	"ki/internal/ui/styles"
	"ki/internal/ui/views"
)

// App wraps the models.Model to provide the necessary methods
type App struct {
	model models.Model
}

// NewApp creates a new app instance
func NewApp() *App {
	// Create main menu items
	mainItems := []list.Item{
		models.NewItem("List Clusters", "View and delete KIND clusters", "clusters"),
		models.NewItem("Create Cluster", "Create a new KIND cluster", "create"),
		models.NewItem("Load Docker Image", "Load a Docker image into a KIND cluster", "load"),
		models.NewItem("Build Node Image", "Build a custom KIND node image from source", "build"),
		models.NewItem("Export Logs", "Export cluster logs for debugging", "logs"),
	}

	// Setup main menu list
	mainList := list.New(mainItems, list.NewDefaultDelegate(), 0, 0)
	mainList.Title = "Menu"
	mainList.SetShowStatusBar(false)
	mainList.SetFilteringEnabled(false)

	// Setup cluster list
	clusterList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	clusterList.Title = "KIND Clusters"
	clusterList.SetShowStatusBar(false)
	clusterList.SetFilteringEnabled(false)

	// Setup node list
	nodeList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	nodeList.Title = "Cluster Nodes"
	nodeList.SetShowStatusBar(false)
	nodeList.SetFilteringEnabled(false)

	// Setup text input
	ti := textinput.New()
	ti.Placeholder = "Enter value..."
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 50

	m := models.Model{
		CurrentView: models.MainMenuView,
		MainMenu:    mainList,
		ClusterList: clusterList,
		NodeList:    nodeList,
		TextInput:   ti,
		Help:        help.New(),
		Clusters:    []cmd.Cluster{},
		ShowHelp:    false,
	}

	return &App{model: m}
}

func (a *App) Init() tea.Cmd {
	return tea.Batch(
		commands.GetKindClusters(),
		textinput.Blink,
	)
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return a.handleWindowResize(msg)
	case models.ClustersMsg:
		return a.handleClustersMsg(msg)
	case models.NodesMsg:
		return a.handleNodesMsg(msg)
	case models.ClusterDetailMsg:
		return a.handleClusterDetailMsg(msg)
	case models.MessageMsg:
		return a.handleMessageMsg(msg)
	case tea.KeyMsg:
		return a.handleKeyMsg(msg)
	}
	return a, nil
}

func (a *App) View() string {
	if a.model.Quitting {
		return "\nThanks for using KIND Interactive! 🐳\n"
	}

	var content string

	// Header
	header := styles.Title.Render("KIND Interactive")

	// Message display
	message := ""
	if a.model.Message != "" {
		switch a.model.MessageType {
		case "success":
			message = styles.Status.Render("✓ " + a.model.Message)
		case "error":
			message = styles.Error.Render("✗ " + a.model.Message)
		default:
			message = a.model.Message
		}
	}

	// Content based on current view
	switch a.model.CurrentView {
	case models.MainMenuView:
		content = a.model.MainMenu.View()
	case models.ClusterListView:
		content = a.model.ClusterList.View()
	case models.ClusterDetailView:
		content = views.RenderClusterDetail(a.model.CurrentCluster)
	case models.NodeListView:
		content = a.model.NodeList.View()
	case models.DeleteConfirmView:
		content = views.RenderDeleteConfirmation(a.model.ClusterToDelete, a.model.DeleteConfirmChoice)
	case models.CreateClusterView, models.LoadImageView, models.BuildImageView, models.ExportLogsView:
		content = fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			a.model.InputPrompt,
			a.model.TextInput.View(),
			styles.Help.Render("Press Enter to confirm, Esc to cancel"),
		)
	}

	// Help
	helpView := ""
	if a.model.ShowHelp {
		helpView = "\n" + a.model.Help.View(models.Keys)
	}

	// Combine all parts
	parts := []string{header}
	if message != "" {
		parts = append(parts, message)
	}
	parts = append(parts, content)
	if helpView != "" {
		parts = append(parts, helpView)
	}

	// Footer with quick help based on current view
	footer := ""
	switch a.model.CurrentView {
	case models.MainMenuView:
		footer = styles.Help.Render("\nc create • ? help • q quit")
	case models.ClusterListView:
		footer = styles.Help.Render("\nenter/i info • n nodes • d delete • c create • l load • L logs • r refresh • esc back • ? help • q quit")
	case models.ClusterDetailView, models.NodeListView:
		footer = styles.Help.Render("\nesc back • ? help • q quit")
	case models.DeleteConfirmView:
		footer = styles.Help.Render("\n←/→/tab select • enter confirm • y yes • n no • esc cancel • q quit")
	default:
		footer = styles.Help.Render("\nenter confirm • esc cancel • ? help • q quit")
	}
	parts = append(parts, footer)

	return strings.Join(parts, "\n")
}

func (a *App) handleWindowResize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	a.model.Width = msg.Width
	a.model.Height = msg.Height

	// Update component sizes
	a.model.MainMenu.SetWidth(msg.Width)
	a.model.MainMenu.SetHeight(msg.Height - 8)
	a.model.ClusterList.SetWidth(msg.Width)
	a.model.ClusterList.SetHeight(msg.Height - 8)
	a.model.NodeList.SetWidth(msg.Width)
	a.model.NodeList.SetHeight(msg.Height - 8)
	a.model.Help.Width = msg.Width

	return a, nil
}

func (a *App) handleClustersMsg(msg models.ClustersMsg) (tea.Model, tea.Cmd) {
	a.model.Clusters = []cmd.Cluster(msg)

	// Update cluster list items
	items := make([]list.Item, len(a.model.Clusters))
	for i, cluster := range a.model.Clusters {
		nodeCount := len(cluster.Nodes)
		description := fmt.Sprintf("Status: %s | Nodes: %d", cluster.Status, nodeCount)
		if cluster.KubeVersion != "" {
			description += fmt.Sprintf(" | K8s: %s", cluster.KubeVersion)
		}

		items[i] = models.NewItem(cluster.Name, description, "select")
	}
	a.model.ClusterList.SetItems(items)

	return a, nil
}

func (a *App) handleNodesMsg(msg models.NodesMsg) (tea.Model, tea.Cmd) {
	nodes := []cmd.Node(msg)

	// Update node list items
	items := make([]list.Item, len(nodes))
	for i, node := range nodes {
		items[i] = models.NewItem(node.Name, fmt.Sprintf("Role: %s | Status: %s | Age: %s | IP: %s", node.Role, node.Status, node.Age, node.InternalIP), "select")
	}
	a.model.NodeList.SetItems(items)

	return a, nil
}

func (a *App) handleClusterDetailMsg(msg models.ClusterDetailMsg) (tea.Model, tea.Cmd) {
	cluster := cmd.Cluster(msg)
	a.model.CurrentCluster = &cluster
	return a, nil
}

func (a *App) handleMessageMsg(msg models.MessageMsg) (tea.Model, tea.Cmd) {
	a.model.Message = msg.Text
	a.model.MessageType = msg.MsgType

	var cmds []tea.Cmd

	// Auto-clear message after 5 seconds
	cmds = append(cmds, func() tea.Cmd {
		return tea.Tick(time.Second*5, func(time.Time) tea.Msg {
			return models.MessageMsg{Text: "", MsgType: ""}
		})
	}())

	// Refresh clusters after operations
	if msg.MsgType == "success" && strings.Contains(msg.Text, "Cluster") {
		cmds = append(cmds, commands.GetKindClusters())
	}

	return a, tea.Batch(cmds...)
}

func (a *App) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, models.Keys.Quit):
		a.model.Quitting = true
		return a, tea.Quit

	case key.Matches(msg, models.Keys.Help):
		a.model.ShowHelp = !a.model.ShowHelp
		return a, nil

	case key.Matches(msg, models.Keys.Back):
		if a.model.CurrentView != models.MainMenuView {
			if a.model.CurrentView == models.DeleteConfirmView {
				// Cancel deletion, go back to cluster list
				a.model.CurrentView = models.ClusterListView
				a.model.ClusterToDelete = ""
			} else {
				a.model.CurrentView = models.MainMenuView
				a.model.TextInput.SetValue("")
				a.model.InputAction = ""
				a.model.SelectedCluster = ""
			}
			return a, nil
		}
	}

	// Handle view-specific key presses
	switch a.model.CurrentView {
	case models.MainMenuView:
		return a.handleMainMenuKeys(msg)
	case models.ClusterListView:
		return a.handleClusterListKeys(msg)
	case models.ClusterDetailView:
		return a.handleClusterDetailKeys(msg)
	case models.NodeListView:
		return a.handleNodeListKeys(msg)
	case models.DeleteConfirmView:
		return a.handleDeleteConfirmKeys(msg)
	case models.CreateClusterView, models.LoadImageView, models.BuildImageView, models.ExportLogsView:
		return a.handleInputKeys(msg)
	}

	return a, nil
}