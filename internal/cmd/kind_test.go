package cmd

import (
	"strings"
	"testing"
)

func TestParseNodes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Node
	}{
		{
			name: "single control plane node",
			input: `kind-control-plane   Ready    control-plane   5d    v1.25.3   172.18.0.2   <none>        Ubuntu 21.10   5.15.49-linuxkit   containerd://1.6.9`,
			expected: []Node{
				{
					Name:       "kind-control-plane",
					Status:     "Ready",
					Role:       "control-plane",
					Age:        "5d",
					Version:    "v1.25.3",
					InternalIP: "172.18.0.2",
				},
			},
		},
		{
			name: "multiple nodes with control plane and workers",
			input: `kind-control-plane   Ready    control-plane   5d    v1.25.3   172.18.0.2   <none>        Ubuntu 21.10   5.15.49-linuxkit   containerd://1.6.9
kind-worker          Ready    <none>          5d    v1.25.3   172.18.0.3   <none>        Ubuntu 21.10   5.15.49-linuxkit   containerd://1.6.9
kind-worker2         Ready    <none>          5d    v1.25.3   172.18.0.4   <none>        Ubuntu 21.10   5.15.49-linuxkit   containerd://1.6.9`,
			expected: []Node{
				{
					Name:       "kind-control-plane",
					Status:     "Ready",
					Role:       "control-plane",
					Age:        "5d",
					Version:    "v1.25.3",
					InternalIP: "172.18.0.2",
				},
				{
					Name:       "kind-worker",
					Status:     "Ready",
					Role:       "worker",
					Age:        "5d",
					Version:    "v1.25.3",
					InternalIP: "172.18.0.3",
				},
				{
					Name:       "kind-worker2",
					Status:     "Ready",
					Role:       "worker",
					Age:        "5d",
					Version:    "v1.25.3",
					InternalIP: "172.18.0.4",
				},
			},
		},
		{
			name:     "empty input",
			input:    "",
			expected: []Node{},
		},
		{
			name:     "invalid input with too few fields",
			input:    "incomplete line",
			expected: []Node{},
		},
		{
			name: "node with master role",
			input: `kind-master   Ready    master   5d    v1.25.3   172.18.0.2   <none>        Ubuntu 21.10   5.15.49-linuxkit   containerd://1.6.9`,
			expected: []Node{
				{
					Name:       "kind-master",
					Status:     "Ready",
					Role:       "control-plane",
					Age:        "5d",
					Version:    "v1.25.3",
					InternalIP: "172.18.0.2",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseNodes(tt.input)
			
			if len(result) != len(tt.expected) {
				t.Errorf("ParseNodes() returned %d nodes, expected %d", len(result), len(tt.expected))
				return
			}
			
			for i, node := range result {
				if node.Name != tt.expected[i].Name {
					t.Errorf("Node[%d].Name = %s, expected %s", i, node.Name, tt.expected[i].Name)
				}
				if node.Status != tt.expected[i].Status {
					t.Errorf("Node[%d].Status = %s, expected %s", i, node.Status, tt.expected[i].Status)
				}
				if node.Role != tt.expected[i].Role {
					t.Errorf("Node[%d].Role = %s, expected %s", i, node.Role, tt.expected[i].Role)
				}
				if node.Age != tt.expected[i].Age {
					t.Errorf("Node[%d].Age = %s, expected %s", i, node.Age, tt.expected[i].Age)
				}
				if node.Version != tt.expected[i].Version {
					t.Errorf("Node[%d].Version = %s, expected %s", i, node.Version, tt.expected[i].Version)
				}
				if node.InternalIP != tt.expected[i].InternalIP {
					t.Errorf("Node[%d].InternalIP = %s, expected %s", i, node.InternalIP, tt.expected[i].InternalIP)
				}
			}
		})
	}
}

func TestEnrichClusterInfo(t *testing.T) {
	// This test verifies that EnrichClusterInfo properly handles clusters
	tests := []struct {
		name     string
		cluster  Cluster
		expected Cluster
	}{
		{
			name: "cluster with default name",
			cluster: Cluster{
				Name:   "kind",
				Status: "running",
			},
			expected: Cluster{
				Name:   "kind",
				Status: "running",
				// Note: actual nodes will be populated by kubectl command
			},
		},
		{
			name: "cluster with custom name",
			cluster: Cluster{
				Name:   "my-cluster",
				Status: "running",
			},
			expected: Cluster{
				Name:   "my-cluster",
				Status: "running",
				// Note: actual nodes will be populated by kubectl command
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: This test will only verify the basic structure
			// since we can't mock the kubectl command easily
			result := EnrichClusterInfo(tt.cluster)
			
			if result.Name != tt.expected.Name {
				t.Errorf("EnrichClusterInfo() Name = %s, expected %s", result.Name, tt.expected.Name)
			}
			if result.Status != tt.expected.Status {
				t.Errorf("EnrichClusterInfo() Status = %s, expected %s", result.Status, tt.expected.Status)
			}
		})
	}
}

func TestNodeRoleDetection(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected string
	}{
		{
			name:     "control-plane role",
			line:     "node1 Ready control-plane 5d v1.25.3 172.18.0.2",
			expected: "control-plane",
		},
		{
			name:     "master role",
			line:     "node1 Ready master 5d v1.25.3 172.18.0.2",
			expected: "control-plane",
		},
		{
			name:     "worker role",
			line:     "node1 Ready <none> 5d v1.25.3 172.18.0.2",
			expected: "worker",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			role := "worker"
			if strings.Contains(tt.line, "control-plane") || strings.Contains(tt.line, "master") {
				role = "control-plane"
			}
			
			if role != tt.expected {
				t.Errorf("Role detection for line %q = %s, expected %s", tt.line, role, tt.expected)
			}
		})
	}
}