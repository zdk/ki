package views

import (
	"strings"
	"testing"
)

func TestRenderDeleteConfirmation(t *testing.T) {
	tests := []struct {
		name                string
		clusterToDelete     string
		deleteConfirmChoice int
		contains            []string
		checkFocused        string // "yes" or "no"
	}{
		{
			name:                "confirm delete with yes selected",
			clusterToDelete:     "test-cluster",
			deleteConfirmChoice: 0,
			contains: []string{
				"DELETE CLUSTER CONFIRMATION",
				"You are about to delete the following cluster:",
				"Cluster Name: test-cluster",
				"WARNING: This action cannot be undone!",
				"All pods, services, and data in this cluster will be permanently lost.",
				"Are you sure you want to delete this cluster?",
				"[Y] Yes, delete the cluster",
				"[N] No, cancel",
				"Use ←/→/Tab to select, Enter to confirm, or press Y/N directly",
			},
			checkFocused: "yes",
		},
		{
			name:                "confirm delete with no selected",
			clusterToDelete:     "production-cluster",
			deleteConfirmChoice: 1,
			contains: []string{
				"DELETE CLUSTER CONFIRMATION",
				"Cluster Name: production-cluster",
				"WARNING: This action cannot be undone!",
				"[Y] Yes, delete the cluster",
				"[N] No, cancel",
			},
			checkFocused: "no",
		},
		{
			name:                "empty cluster name",
			clusterToDelete:     "",
			deleteConfirmChoice: 0,
			contains: []string{
				"DELETE CLUSTER CONFIRMATION",
				"Cluster Name:",
				"WARNING: This action cannot be undone!",
			},
			checkFocused: "yes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderDeleteConfirmation(tt.clusterToDelete, tt.deleteConfirmChoice)

			// Check that all expected strings are present
			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("RenderDeleteConfirmation() should contain %q, but it doesn't.\nGot:\n%s", expected, result)
				}
			}

			// Check button focus state
			if tt.checkFocused == "yes" && tt.deleteConfirmChoice == 0 {
				// When Yes is selected, it should be styled differently than No
				// We can't easily test the actual styling, but we can verify the structure
				if !strings.Contains(result, "[Y] Yes, delete the cluster") {
					t.Error("Yes option should be present when selected")
				}
			} else if tt.checkFocused == "no" && tt.deleteConfirmChoice == 1 {
				// When No is selected
				if !strings.Contains(result, "[N] No, cancel") {
					t.Error("No option should be present when selected")
				}
			}
		})
	}
}

func TestRenderDeleteConfirmationWarnings(t *testing.T) {
	result := RenderDeleteConfirmation("critical-cluster", 0)

	// Check for multiple warning indicators
	warningCount := strings.Count(result, "⚠️")
	if warningCount < 2 {
		t.Errorf("Expected at least 2 warning symbols, got %d", warningCount)
	}

	// Check for strong warning language
	strongWarnings := []string{
		"cannot be undone",
		"permanently lost",
		"WARNING",
	}

	for _, warning := range strongWarnings {
		if !strings.Contains(result, warning) {
			t.Errorf("Delete confirmation should contain strong warning: %q", warning)
		}
	}
}

func TestRenderDeleteConfirmationLayout(t *testing.T) {
	result := RenderDeleteConfirmation("test", 0)
	lines := strings.Split(result, "\n")

	// Check that we have a reasonable number of lines
	if len(lines) < 10 {
		t.Errorf("Delete confirmation should have at least 10 lines, got %d", len(lines))
	}

	// Check for proper spacing (empty lines for readability)
	emptyLineCount := 0
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			emptyLineCount++
		}
	}

	if emptyLineCount < 2 {
		t.Error("Delete confirmation should have empty lines for better readability")
	}
}