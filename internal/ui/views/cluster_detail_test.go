package views

import (
	"strings"
	"testing"

	"ki/internal/cmd"
)

func TestRenderClusterDetail(t *testing.T) {
	tests := []struct {
		name     string
		cluster  *cmd.Cluster
		contains []string
	}{
		{
			name:     "nil cluster",
			cluster:  nil,
			contains: []string{"Loading cluster details..."},
		},
		{
			name: "cluster without nodes",
			cluster: &cmd.Cluster{
				Name:        "test-cluster",
				Status:      "running",
				KubeVersion: "v1.25.3",
				Nodes:       []cmd.Node{},
			},
			contains: []string{
				"Cluster: test-cluster",
				"Basic Information:",
				"Name: test-cluster",
				"Status: running",
				"Kubernetes Version: v1.25.3",
				"Node Count: 0",
			},
		},
		{
			name: "cluster with nodes",
			cluster: &cmd.Cluster{
				Name:        "kind",
				Status:      "running",
				KubeVersion: "v1.25.3",
				Nodes: []cmd.Node{
					{
						Name:       "kind-control-plane",
						Role:       "control-plane",
						Status:     "Ready",
						Age:        "5d",
						Version:    "v1.25.3",
						InternalIP: "172.18.0.2",
					},
					{
						Name:       "kind-worker",
						Role:       "worker",
						Status:     "Ready",
						Age:        "5d",
						Version:    "v1.25.3",
						InternalIP: "172.18.0.3",
					},
				},
			},
			contains: []string{
				"Cluster: kind",
				"Basic Information:",
				"Name: kind",
				"Status: running",
				"Kubernetes Version: v1.25.3",
				"Node Count: 2",
				"Nodes:",
				"NAME",
				"ROLE",
				"STATUS",
				"AGE",
				"INTERNAL-IP",
				"kind-control-plane",
				"control-plane",
				"Ready",
				"5d",
				"172.18.0.2",
				"kind-worker",
				"worker",
				"172.18.0.3",
			},
		},
		{
			name: "cluster without kube version",
			cluster: &cmd.Cluster{
				Name:        "test",
				Status:      "running",
				KubeVersion: "",
				Nodes:       []cmd.Node{},
			},
			contains: []string{
				"Cluster: test",
				"Name: test",
				"Status: running",
				"Node Count: 0",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderClusterDetail(tt.cluster)

			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("RenderClusterDetail() should contain %q, but it doesn't.\nGot:\n%s", expected, result)
				}
			}

			// Check that Kubernetes Version is not shown when empty
			if tt.cluster != nil && tt.cluster.KubeVersion == "" {
				if strings.Contains(result, "Kubernetes Version:") {
					t.Error("RenderClusterDetail() should not show Kubernetes Version when it's empty")
				}
			}

			// Always check for the help text at the bottom
			if tt.cluster != nil && !strings.Contains(result, "Press 'esc' to go back") {
				t.Error("RenderClusterDetail() should contain help text")
			}
		})
	}
}

func TestRenderClusterDetailFormatting(t *testing.T) {
	cluster := &cmd.Cluster{
		Name:   "test-cluster",
		Status: "running",
		Nodes: []cmd.Node{
			{
				Name:       "very-long-node-name-that-should-fit",
				Role:       "control-plane",
				Status:     "Ready",
				Age:        "10d15h",
				Version:    "v1.25.3",
				InternalIP: "192.168.100.100",
			},
		},
	}

	result := RenderClusterDetail(cluster)

	// Check for proper table alignment
	lines := strings.Split(result, "\n")
	var separatorLine, dataLine string
	
	for i, line := range lines {
		if strings.Contains(line, "NAME") && strings.Contains(line, "ROLE") {
			if i+1 < len(lines) {
				separatorLine = lines[i+1]
			}
			if i+2 < len(lines) {
				dataLine = lines[i+2]
			}
			break
		}
	}

	// Verify separator line has dashes
	if !strings.Contains(separatorLine, "---") {
		t.Error("Separator line should contain dashes")
	}

	// Verify data is present in the data line
	if !strings.Contains(dataLine, "very-long-node-name-that-should-fit") {
		t.Error("Data line should contain the node name")
	}
}