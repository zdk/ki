package models

import "testing"

func TestViewMode(t *testing.T) {
	tests := []struct {
		name     string
		mode     ViewMode
		expected ViewMode
	}{
		{
			name:     "main menu view",
			mode:     MainMenuView,
			expected: 0,
		},
		{
			name:     "cluster list view",
			mode:     ClusterListView,
			expected: 1,
		},
		{
			name:     "cluster detail view",
			mode:     ClusterDetailView,
			expected: 2,
		},
		{
			name:     "node list view",
			mode:     NodeListView,
			expected: 3,
		},
		{
			name:     "create cluster view",
			mode:     CreateClusterView,
			expected: 4,
		},
		{
			name:     "load image view",
			mode:     LoadImageView,
			expected: 5,
		},
		{
			name:     "build image view",
			mode:     BuildImageView,
			expected: 6,
		},
		{
			name:     "export logs view",
			mode:     ExportLogsView,
			expected: 7,
		},
		{
			name:     "delete confirm view",
			mode:     DeleteConfirmView,
			expected: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mode != tt.expected {
				t.Errorf("ViewMode %s = %d, want %d", tt.name, tt.mode, tt.expected)
			}
		})
	}
}

func TestViewModeValues(t *testing.T) {
	// Ensure view modes have distinct values
	modes := map[ViewMode]string{
		MainMenuView:      "MainMenuView",
		ClusterListView:   "ClusterListView",
		ClusterDetailView: "ClusterDetailView",
		NodeListView:      "NodeListView",
		CreateClusterView: "CreateClusterView",
		LoadImageView:     "LoadImageView",
		BuildImageView:    "BuildImageView",
		ExportLogsView:    "ExportLogsView",
		DeleteConfirmView: "DeleteConfirmView",
	}

	// Check for duplicate values
	seen := make(map[ViewMode]bool)
	for mode, name := range modes {
		if seen[mode] {
			t.Errorf("Duplicate ViewMode value %d for %s", mode, name)
		}
		seen[mode] = true
	}

	// Verify we have the expected number of modes
	expectedCount := 9
	if len(modes) != expectedCount {
		t.Errorf("Expected %d view modes, got %d", expectedCount, len(modes))
	}
}