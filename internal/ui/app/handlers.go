package app

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"ki/internal/ui/commands"
	"ki/internal/ui/models"
)

func (a *App) handleMainMenuKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch {
	case key.Matches(msg, models.Keys.Enter):
		selectedItem := a.model.MainMenu.SelectedItem()
		if item, ok := selectedItem.(models.Item); ok {
			switch item.Action {
			case "clusters":
				a.model.CurrentView = models.ClusterListView
				return a, commands.GetKindClusters()
			case "create":
				a.model.CurrentView = models.CreateClusterView
				a.model.InputPrompt = "Enter cluster name (leave empty for 'kind'):"
				a.model.InputAction = "create"
				a.model.TextInput.Placeholder = "cluster-name"
				a.model.TextInput.Focus()
				return a, nil
			case "load":
				a.model.CurrentView = models.LoadImageView
				a.model.InputPrompt = "Enter Docker image name:"
				a.model.InputAction = "load-image"
				a.model.TextInput.Placeholder = "nginx:latest"
				a.model.TextInput.Focus()
				return a, nil
			case "build":
				a.model.CurrentView = models.BuildImageView
				a.model.InputPrompt = "Enter Kubernetes source path (leave empty for default):"
				a.model.InputAction = "build"
				a.model.TextInput.Placeholder = "/path/to/kubernetes/source"
				a.model.TextInput.Focus()
				return a, nil
			case "logs":
				a.model.CurrentView = models.ExportLogsView
				a.model.InputPrompt = "Enter output directory (leave empty for current dir):"
				a.model.InputAction = "export-logs"
				a.model.TextInput.Placeholder = "./logs"
				a.model.TextInput.Focus()
				return a, nil
			}
		}
	case key.Matches(msg, models.Keys.Create):
		a.model.CurrentView = models.CreateClusterView
		a.model.InputPrompt = "Enter cluster name:"
		a.model.InputAction = "create"
		a.model.TextInput.Focus()
		return a, nil
	case key.Matches(msg, models.Keys.Refresh):
		return a, commands.GetKindClusters()
	}

	a.model.MainMenu, cmd = a.model.MainMenu.Update(msg)
	cmds = append(cmds, cmd)

	return a, tea.Batch(cmds...)
}

func (a *App) handleClusterListKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch {
	case key.Matches(msg, models.Keys.Enter), key.Matches(msg, models.Keys.Detail):
		selectedItem := a.model.ClusterList.SelectedItem()
		if item, ok := selectedItem.(models.Item); ok {
			a.model.CurrentView = models.ClusterDetailView
			return a, commands.GetClusterDetail(item.Title())
		}
	case key.Matches(msg, models.Keys.Nodes):
		selectedItem := a.model.ClusterList.SelectedItem()
		if item, ok := selectedItem.(models.Item); ok {
			a.model.CurrentView = models.NodeListView
			a.model.NodeList.Title = "Nodes - " + item.Title()
			return a, commands.GetClusterNodes(item.Title())
		}
	case key.Matches(msg, models.Keys.Delete):
		selectedItem := a.model.ClusterList.SelectedItem()
		if item, ok := selectedItem.(models.Item); ok {
			// Show confirmation dialog
			a.model.ClusterToDelete = item.Title()
			a.model.DeleteConfirmChoice = 0 // Default to "Yes"
			a.model.CurrentView = models.DeleteConfirmView
			return a, nil
		}
	case key.Matches(msg, models.Keys.Create):
		a.model.CurrentView = models.CreateClusterView
		a.model.InputPrompt = "Enter cluster name:"
		a.model.InputAction = "create"
		a.model.TextInput.Focus()
		return a, nil
	case key.Matches(msg, models.Keys.Load):
		selectedItem := a.model.ClusterList.SelectedItem()
		if item, ok := selectedItem.(models.Item); ok {
			a.model.SelectedCluster = item.Title()
			a.model.CurrentView = models.LoadImageView
			a.model.InputPrompt = "Enter Docker image name to load into '" + item.Title() + "':"
			a.model.InputAction = "load-image"
			a.model.TextInput.Placeholder = "nginx:latest"
			a.model.TextInput.Focus()
			return a, nil
		}
	case key.Matches(msg, models.Keys.Logs):
		selectedItem := a.model.ClusterList.SelectedItem()
		if item, ok := selectedItem.(models.Item); ok {
			a.model.SelectedCluster = item.Title()
			a.model.CurrentView = models.ExportLogsView
			a.model.InputPrompt = "Enter output directory for '" + item.Title() + "' logs:"
			a.model.InputAction = "export-logs"
			a.model.TextInput.Placeholder = "./logs"
			a.model.TextInput.Focus()
			return a, nil
		}
	case key.Matches(msg, models.Keys.Refresh):
		return a, commands.GetKindClusters()
	}

	a.model.ClusterList, cmd = a.model.ClusterList.Update(msg)
	cmds = append(cmds, cmd)

	return a, tea.Batch(cmds...)
}

func (a *App) handleClusterDetailKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	a.model.ClusterList, cmd = a.model.ClusterList.Update(msg)
	return a, cmd
}

func (a *App) handleNodeListKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	a.model.NodeList, cmd = a.model.NodeList.Update(msg)
	return a, cmd
}

func (a *App) handleDeleteConfirmKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, models.Keys.Left), key.Matches(msg, models.Keys.Right), key.Matches(msg, models.Keys.Tab):
		// Toggle between Yes (0) and No (1)
		a.model.DeleteConfirmChoice = 1 - a.model.DeleteConfirmChoice
		return a, nil
	case key.Matches(msg, models.Keys.Enter):
		// Execute based on current selection
		if a.model.DeleteConfirmChoice == 0 {
			// Yes - delete cluster
			clusterName := a.model.ClusterToDelete
			a.model.ClusterToDelete = ""
			a.model.DeleteConfirmChoice = 0
			a.model.CurrentView = models.ClusterListView
			return a, commands.DeleteKindCluster(clusterName)
		} else {
			// No - cancel deletion
			a.model.ClusterToDelete = ""
			a.model.DeleteConfirmChoice = 0
			a.model.CurrentView = models.ClusterListView
			return a, nil
		}
	case key.Matches(msg, models.Keys.Yes):
		// Direct 'y' key - delete cluster
		clusterName := a.model.ClusterToDelete
		a.model.ClusterToDelete = ""
		a.model.DeleteConfirmChoice = 0
		a.model.CurrentView = models.ClusterListView
		return a, commands.DeleteKindCluster(clusterName)
	case key.Matches(msg, models.Keys.No):
		// Direct 'n' key - cancel deletion
		a.model.ClusterToDelete = ""
		a.model.DeleteConfirmChoice = 0
		a.model.CurrentView = models.ClusterListView
		return a, nil
	}

	return a, nil
}

func (a *App) handleInputKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg.String() {
	case "enter":
		inputValue := strings.TrimSpace(a.model.TextInput.Value())
		a.model.TextInput.SetValue("")
		a.model.CurrentView = models.MainMenuView

		switch a.model.InputAction {
		case "create":
			clusterName := inputValue
			if clusterName == "" {
				clusterName = "kind"
			}
			return a, commands.CreateKindCluster(clusterName)

		case "load-image":
			if inputValue == "" {
				return a, func() tea.Msg {
					return models.MessageMsg{
						Text:    "Image name cannot be empty",
						MsgType: "error",
					}
				}
			}
			return a, commands.LoadDockerImage(inputValue, a.model.SelectedCluster)

		case "build":
			sourcePath := inputValue
			return a, commands.BuildNodeImage(sourcePath)

		case "export-logs":
			outputPath := inputValue
			return a, commands.ExportKindLogs(a.model.SelectedCluster, outputPath)
		}
	}

	a.model.TextInput, cmd = a.model.TextInput.Update(msg)
	cmds = append(cmds, cmd)

	return a, tea.Batch(cmds...)
}