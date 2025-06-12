package views

import (
	"fmt"
	"strings"

	"ki/internal/ui/styles"
)

// RenderDeleteConfirmation renders the delete confirmation dialog
func RenderDeleteConfirmation(clusterToDelete string, deleteConfirmChoice int) string {
	var content strings.Builder

	// Warning header
	content.WriteString(styles.Error.Render("⚠️  DELETE CLUSTER CONFIRMATION"))
	content.WriteString("\n\n")

	// Cluster info
	content.WriteString("You are about to delete the following cluster:\n\n")
	content.WriteString(styles.Status.Render(fmt.Sprintf("Cluster Name: %s", clusterToDelete)))
	content.WriteString("\n\n")

	// Warning message
	content.WriteString(styles.Error.Render("⚠️  WARNING: This action cannot be undone!"))
	content.WriteString("\n")
	content.WriteString("All pods, services, and data in this cluster will be permanently lost.\n\n")

	// Confirmation prompt
	content.WriteString("Are you sure you want to delete this cluster?\n\n")

	// Interactive buttons with selection highlighting
	if deleteConfirmChoice == 0 {
		// Yes is selected
		content.WriteString(styles.Focused.Render("  [Y] Yes, delete the cluster  "))
		content.WriteString("  ")
		content.WriteString(styles.Blurred.Render("  [N] No, cancel  "))
	} else {
		// No is selected
		content.WriteString(styles.Blurred.Render("  [Y] Yes, delete the cluster  "))
		content.WriteString("  ")
		content.WriteString(styles.Focused.Render("  [N] No, cancel  "))
	}
	content.WriteString("\n\n")

	content.WriteString(styles.Help.Render("Use ←/→/Tab to select, Enter to confirm, or press Y/N directly"))

	return content.String()
}