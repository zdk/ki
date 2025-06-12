package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"ki/internal/cmd"
	"ki/internal/ui/models"
)

// GetKindClusters fetches all KIND clusters
func GetKindClusters() tea.Cmd {
	return func() tea.Msg {
		clusters, err := cmd.Commands.GetClusters()
		if err != nil {
			return models.MessageMsg{
				Text:    err.Error(),
				MsgType: "error",
			}
		}
		return models.ClustersMsg(clusters)
	}
}

// GetClusterNodes fetches nodes for a specific cluster
func GetClusterNodes(clusterName string) tea.Cmd {
	return func() tea.Msg {
		nodes, err := cmd.Commands.GetClusterNodes(clusterName)
		if err != nil {
			return models.MessageMsg{
				Text:    err.Error(),
				MsgType: "error",
			}
		}
		return models.NodesMsg(nodes)
	}
}

// GetClusterDetail fetches detailed information about a cluster
func GetClusterDetail(clusterName string) tea.Cmd {
	return func() tea.Msg {
		cluster, err := cmd.Commands.GetClusterDetail(clusterName)
		if err != nil {
			return models.MessageMsg{
				Text:    err.Error(),
				MsgType: "error",
			}
		}
		return models.ClusterDetailMsg(cluster)
	}
}

// CreateKindCluster creates a new KIND cluster
func CreateKindCluster(name string) tea.Cmd {
	return func() tea.Msg {
		if err := cmd.Commands.CreateCluster(name); err != nil {
			return models.MessageMsg{
				Text:    err.Error(),
				MsgType: "error",
			}
		}
		return models.MessageMsg{
			Text:    fmt.Sprintf("Cluster '%s' created successfully!", name),
			MsgType: "success",
		}
	}
}

// DeleteKindCluster deletes a KIND cluster
func DeleteKindCluster(name string) tea.Cmd {
	return func() tea.Msg {
		if err := cmd.Commands.DeleteCluster(name); err != nil {
			return models.MessageMsg{
				Text:    err.Error(),
				MsgType: "error",
			}
		}
		return models.MessageMsg{
			Text:    fmt.Sprintf("Cluster '%s' deleted successfully!", name),
			MsgType: "success",
		}
	}
}

// LoadDockerImage loads a Docker image into a KIND cluster
func LoadDockerImage(imageName, clusterName string) tea.Cmd {
	return func() tea.Msg {
		if err := cmd.Commands.LoadDockerImage(imageName, clusterName); err != nil {
			return models.MessageMsg{
				Text:    err.Error(),
				MsgType: "error",
			}
		}
		return models.MessageMsg{
			Text:    fmt.Sprintf("Image '%s' loaded successfully!", imageName),
			MsgType: "success",
		}
	}
}

// BuildNodeImage builds a KIND node image from source
func BuildNodeImage(sourcePath string) tea.Cmd {
	return func() tea.Msg {
		if err := cmd.Commands.BuildNodeImage(sourcePath); err != nil {
			return models.MessageMsg{
				Text:    err.Error(),
				MsgType: "error",
			}
		}
		return models.MessageMsg{
			Text:    "Node image built successfully!",
			MsgType: "success",
		}
	}
}

// ExportKindLogs exports cluster logs
func ExportKindLogs(clusterName, outputPath string) tea.Cmd {
	return func() tea.Msg {
		// Expand ~ to home directory
		if outputPath != "" && strings.HasPrefix(outputPath, "~/") {
			home, _ := os.UserHomeDir()
			outputPath = filepath.Join(home, outputPath[2:])
		}

		if err := cmd.Commands.ExportLogs(clusterName, outputPath); err != nil {
			return models.MessageMsg{
				Text:    err.Error(),
				MsgType: "error",
			}
		}

		exportPath := outputPath
		if exportPath == "" {
			exportPath = "current directory"
		}

		return models.MessageMsg{
			Text:    fmt.Sprintf("Logs exported to %s successfully!", exportPath),
			MsgType: "success",
		}
	}
}