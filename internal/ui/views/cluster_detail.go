package views

import (
	"fmt"
	"strings"

	"ki/internal/cmd"
	"ki/internal/ui/styles"
)

// RenderClusterDetail renders the cluster detail view
func RenderClusterDetail(cluster *cmd.Cluster) string {
	if cluster == nil {
		return "Loading cluster details..."
	}

	var details strings.Builder

	// Cluster header
	details.WriteString(styles.Title.Render(fmt.Sprintf("Cluster: %s", cluster.Name)))
	details.WriteString("\n\n")

	// Basic info
	details.WriteString(styles.Status.Render("Basic Information:"))
	details.WriteString("\n")
	details.WriteString(fmt.Sprintf("• Name: %s\n", cluster.Name))
	details.WriteString(fmt.Sprintf("• Status: %s\n", cluster.Status))
	if cluster.KubeVersion != "" {
		details.WriteString(fmt.Sprintf("• Kubernetes Version: %s\n", cluster.KubeVersion))
	}
	details.WriteString(fmt.Sprintf("• Node Count: %d\n", len(cluster.Nodes)))

	// Nodes section
	if len(cluster.Nodes) > 0 {
		details.WriteString("\n")
		details.WriteString(styles.Status.Render("Nodes:"))
		details.WriteString("\n")

		// Header
		details.WriteString(fmt.Sprintf("%-25s %-15s %-10s %-10s %-15s\n",
			"NAME", "ROLE", "STATUS", "AGE", "INTERNAL-IP"))
		details.WriteString(strings.Repeat("-", 80))
		details.WriteString("\n")

		// Node rows
		for _, node := range cluster.Nodes {
			details.WriteString(fmt.Sprintf("%-25s %-15s %-10s %-10s %-15s\n",
				node.Name, node.Role, node.Status, node.Age, node.InternalIP))
		}
	}

	details.WriteString("\n")
	details.WriteString(styles.Help.Render("Press 'esc' to go back to cluster list"))

	return details.String()
}