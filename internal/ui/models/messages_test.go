package models

import (
	"testing"

	"ki/internal/cmd"
)

func TestClustersMsg(t *testing.T) {
	tests := []struct {
		name     string
		clusters []cmd.Cluster
	}{
		{
			name:     "empty clusters",
			clusters: []cmd.Cluster{},
		},
		{
			name: "single cluster",
			clusters: []cmd.Cluster{
				{Name: "test-cluster", Status: "running"},
			},
		},
		{
			name: "multiple clusters",
			clusters: []cmd.Cluster{
				{Name: "cluster1", Status: "running"},
				{Name: "cluster2", Status: "stopped"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := ClustersMsg(tt.clusters)
			
			// Verify type conversion works
			clusters := []cmd.Cluster(msg)
			if len(clusters) != len(tt.clusters) {
				t.Errorf("Expected %d clusters, got %d", len(tt.clusters), len(clusters))
			}
			
			// Verify content is preserved
			for i, cluster := range clusters {
				if cluster.Name != tt.clusters[i].Name {
					t.Errorf("Cluster[%d].Name = %s, want %s", i, cluster.Name, tt.clusters[i].Name)
				}
				if cluster.Status != tt.clusters[i].Status {
					t.Errorf("Cluster[%d].Status = %s, want %s", i, cluster.Status, tt.clusters[i].Status)
				}
			}
		})
	}
}

func TestNodesMsg(t *testing.T) {
	tests := []struct {
		name  string
		nodes []cmd.Node
	}{
		{
			name:  "empty nodes",
			nodes: []cmd.Node{},
		},
		{
			name: "single node",
			nodes: []cmd.Node{
				{Name: "node1", Status: "Ready", Role: "control-plane"},
			},
		},
		{
			name: "multiple nodes",
			nodes: []cmd.Node{
				{Name: "control-plane", Status: "Ready", Role: "control-plane"},
				{Name: "worker1", Status: "Ready", Role: "worker"},
				{Name: "worker2", Status: "NotReady", Role: "worker"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := NodesMsg(tt.nodes)
			
			// Verify type conversion works
			nodes := []cmd.Node(msg)
			if len(nodes) != len(tt.nodes) {
				t.Errorf("Expected %d nodes, got %d", len(tt.nodes), len(nodes))
			}
			
			// Verify content is preserved
			for i, node := range nodes {
				if node.Name != tt.nodes[i].Name {
					t.Errorf("Node[%d].Name = %s, want %s", i, node.Name, tt.nodes[i].Name)
				}
				if node.Status != tt.nodes[i].Status {
					t.Errorf("Node[%d].Status = %s, want %s", i, node.Status, tt.nodes[i].Status)
				}
				if node.Role != tt.nodes[i].Role {
					t.Errorf("Node[%d].Role = %s, want %s", i, node.Role, tt.nodes[i].Role)
				}
			}
		})
	}
}

func TestClusterDetailMsg(t *testing.T) {
	tests := []struct {
		name    string
		cluster cmd.Cluster
	}{
		{
			name: "basic cluster",
			cluster: cmd.Cluster{
				Name:   "test-cluster",
				Status: "running",
			},
		},
		{
			name: "cluster with nodes",
			cluster: cmd.Cluster{
				Name:        "production",
				Status:      "running",
				KubeVersion: "v1.25.3",
				Nodes: []cmd.Node{
					{Name: "node1", Status: "Ready", Role: "control-plane"},
					{Name: "node2", Status: "Ready", Role: "worker"},
				},
			},
		},
		{
			name: "cluster with all fields",
			cluster: cmd.Cluster{
				Name:        "full-cluster",
				Status:      "running",
				KubeVersion: "v1.26.0",
				Nodes: []cmd.Node{
					{
						Name:       "master-node",
						Status:     "Ready",
						Role:       "control-plane",
						Age:        "5d",
						Version:    "v1.26.0",
						InternalIP: "192.168.1.10",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := ClusterDetailMsg(tt.cluster)
			
			// Verify type conversion works
			cluster := cmd.Cluster(msg)
			
			if cluster.Name != tt.cluster.Name {
				t.Errorf("Cluster.Name = %s, want %s", cluster.Name, tt.cluster.Name)
			}
			if cluster.Status != tt.cluster.Status {
				t.Errorf("Cluster.Status = %s, want %s", cluster.Status, tt.cluster.Status)
			}
			if cluster.KubeVersion != tt.cluster.KubeVersion {
				t.Errorf("Cluster.KubeVersion = %s, want %s", cluster.KubeVersion, tt.cluster.KubeVersion)
			}
			if len(cluster.Nodes) != len(tt.cluster.Nodes) {
				t.Errorf("Expected %d nodes, got %d", len(tt.cluster.Nodes), len(cluster.Nodes))
			}
		})
	}
}

func TestMessageMsg(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		msgType string
	}{
		{
			name:    "success message",
			text:    "Operation completed successfully",
			msgType: "success",
		},
		{
			name:    "error message",
			text:    "Failed to perform operation",
			msgType: "error",
		},
		{
			name:    "info message",
			text:    "Information about the system",
			msgType: "info",
		},
		{
			name:    "empty message",
			text:    "",
			msgType: "",
		},
		{
			name:    "custom message type",
			text:    "Warning message",
			msgType: "warning",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MessageMsg{
				Text:    tt.text,
				MsgType: tt.msgType,
			}
			
			if msg.Text != tt.text {
				t.Errorf("MessageMsg.Text = %s, want %s", msg.Text, tt.text)
			}
			if msg.MsgType != tt.msgType {
				t.Errorf("MessageMsg.MsgType = %s, want %s", msg.MsgType, tt.msgType)
			}
		})
	}
}

func TestMessageMsgFields(t *testing.T) {
	// Test that MessageMsg struct has the expected fields and can be used properly
	msg := MessageMsg{}
	
	// Test that fields can be set
	msg.Text = "Test message"
	msg.MsgType = "test"
	
	if msg.Text != "Test message" {
		t.Error("Failed to set Text field")
	}
	if msg.MsgType != "test" {
		t.Error("Failed to set MsgType field")
	}
	
	// Test struct literal creation
	msg2 := MessageMsg{
		Text:    "Another message",
		MsgType: "another",
	}
	
	if msg2.Text != "Another message" {
		t.Error("Failed to create struct with Text field")
	}
	if msg2.MsgType != "another" {
		t.Error("Failed to create struct with MsgType field")
	}
}

func TestMessageTypes(t *testing.T) {
	// Test common message types used throughout the application
	messageTypes := []string{"success", "error", "info"}
	
	for _, msgType := range messageTypes {
		t.Run("message_type_"+msgType, func(t *testing.T) {
			msg := MessageMsg{
				Text:    "Test message for " + msgType,
				MsgType: msgType,
			}
			
			if msg.MsgType != msgType {
				t.Errorf("Expected message type %s, got %s", msgType, msg.MsgType)
			}
		})
	}
}

func TestMessageConversions(t *testing.T) {
	// Test that message types can be properly used in type assertions
	var msg interface{}
	
	// Test ClustersMsg
	clusters := []cmd.Cluster{{Name: "test", Status: "running"}}
	msg = ClustersMsg(clusters)
	
	if clustersMsg, ok := msg.(ClustersMsg); !ok {
		t.Error("ClustersMsg type assertion failed")
	} else if len(clustersMsg) != 1 {
		t.Error("ClustersMsg content incorrect")
	}
	
	// Test NodesMsg
	nodes := []cmd.Node{{Name: "node1", Status: "Ready"}}
	msg = NodesMsg(nodes)
	
	if nodesMsg, ok := msg.(NodesMsg); !ok {
		t.Error("NodesMsg type assertion failed")
	} else if len(nodesMsg) != 1 {
		t.Error("NodesMsg content incorrect")
	}
	
	// Test ClusterDetailMsg
	cluster := cmd.Cluster{Name: "detail-test", Status: "running"}
	msg = ClusterDetailMsg(cluster)
	
	if clusterMsg, ok := msg.(ClusterDetailMsg); !ok {
		t.Error("ClusterDetailMsg type assertion failed")
	} else if cmd.Cluster(clusterMsg).Name != "detail-test" {
		t.Error("ClusterDetailMsg content incorrect")
	}
	
	// Test MessageMsg
	message := MessageMsg{Text: "test", MsgType: "success"}
	msg = message
	
	if msgMsg, ok := msg.(MessageMsg); !ok {
		t.Error("MessageMsg type assertion failed")
	} else if msgMsg.Text != "test" {
		t.Error("MessageMsg content incorrect")
	}
}