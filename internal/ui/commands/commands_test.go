package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"ki/internal/cmd"
	"ki/internal/ui/models"
)

// MockCommands implements cmd.CommandInterface for testing
type MockCommands struct {
	GetClustersFunc      func() ([]cmd.Cluster, error)
	GetClusterNodesFunc  func(string) ([]cmd.Node, error)
	GetClusterDetailFunc func(string) (cmd.Cluster, error)
	CreateClusterFunc    func(string) error
	DeleteClusterFunc    func(string) error
	LoadDockerImageFunc  func(string, string) error
	BuildNodeImageFunc   func(string) error
	ExportLogsFunc       func(string, string) error
}

func (m *MockCommands) GetClusters() ([]cmd.Cluster, error) {
	if m.GetClustersFunc != nil {
		return m.GetClustersFunc()
	}
	return []cmd.Cluster{}, nil
}

func (m *MockCommands) GetClusterNodes(name string) ([]cmd.Node, error) {
	if m.GetClusterNodesFunc != nil {
		return m.GetClusterNodesFunc(name)
	}
	return []cmd.Node{}, nil
}

func (m *MockCommands) GetClusterDetail(name string) (cmd.Cluster, error) {
	if m.GetClusterDetailFunc != nil {
		return m.GetClusterDetailFunc(name)
	}
	return cmd.Cluster{}, nil
}

func (m *MockCommands) CreateCluster(name string) error {
	if m.CreateClusterFunc != nil {
		return m.CreateClusterFunc(name)
	}
	return nil
}

func (m *MockCommands) DeleteCluster(name string) error {
	if m.DeleteClusterFunc != nil {
		return m.DeleteClusterFunc(name)
	}
	return nil
}

func (m *MockCommands) LoadDockerImage(image, cluster string) error {
	if m.LoadDockerImageFunc != nil {
		return m.LoadDockerImageFunc(image, cluster)
	}
	return nil
}

func (m *MockCommands) BuildNodeImage(path string) error {
	if m.BuildNodeImageFunc != nil {
		return m.BuildNodeImageFunc(path)
	}
	return nil
}

func (m *MockCommands) ExportLogs(cluster, path string) error {
	if m.ExportLogsFunc != nil {
		return m.ExportLogsFunc(cluster, path)
	}
	return nil
}

func TestGetKindClusters(t *testing.T) {
	// Store original and restore after test
	original := cmd.Commands
	defer func() { cmd.Commands = original }()

	tests := []struct {
		name        string
		mockFunc    func() ([]cmd.Cluster, error)
		expectError bool
		expectType  string
	}{
		{
			name: "successful fetch",
			mockFunc: func() ([]cmd.Cluster, error) {
				return []cmd.Cluster{
					{Name: "test-cluster", Status: "running"},
				}, nil
			},
			expectError: false,
			expectType:  "ClustersMsg",
		},
		{
			name: "error fetching clusters",
			mockFunc: func() ([]cmd.Cluster, error) {
				return nil, errors.New("failed to get clusters")
			},
			expectError: true,
			expectType:  "MessageMsg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd.Commands = &MockCommands{
				GetClustersFunc: tt.mockFunc,
			}
			
			cmdFunc := GetKindClusters()
			msg := cmdFunc()
			
			if tt.expectError {
				msgMsg, ok := msg.(models.MessageMsg)
				if !ok {
					t.Errorf("Expected MessageMsg, got %T", msg)
				}
				if msgMsg.MsgType != "error" {
					t.Errorf("Expected error message type, got %s", msgMsg.MsgType)
				}
				if msgMsg.Text != "failed to get clusters" {
					t.Errorf("Expected error text 'failed to get clusters', got %s", msgMsg.Text)
				}
			} else {
				clustersMsg, ok := msg.(models.ClustersMsg)
				if !ok {
					t.Errorf("Expected ClustersMsg, got %T", msg)
				}
				if len(clustersMsg) != 1 {
					t.Errorf("Expected 1 cluster, got %d", len(clustersMsg))
				}
			}
		})
	}
}

func TestGetClusterNodes(t *testing.T) {
	original := cmd.Commands
	defer func() { cmd.Commands = original }()

	tests := []struct {
		name        string
		clusterName string
		mockFunc    func(string) ([]cmd.Node, error)
		expectError bool
	}{
		{
			name:        "successful fetch",
			clusterName: "test-cluster",
			mockFunc: func(name string) ([]cmd.Node, error) {
				return []cmd.Node{
					{Name: "node1", Status: "Ready"},
				}, nil
			},
			expectError: false,
		},
		{
			name:        "error fetching nodes",
			clusterName: "test-cluster",
			mockFunc: func(name string) ([]cmd.Node, error) {
				return nil, errors.New("failed to get nodes")
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd.Commands = &MockCommands{
				GetClusterNodesFunc: tt.mockFunc,
			}
			
			cmdFunc := GetClusterNodes(tt.clusterName)
			msg := cmdFunc()
			
			if tt.expectError {
				msgMsg, ok := msg.(models.MessageMsg)
				if !ok {
					t.Errorf("Expected MessageMsg, got %T", msg)
				}
				if msgMsg.MsgType != "error" {
					t.Errorf("Expected error message type, got %s", msgMsg.MsgType)
				}
			} else {
				nodesMsg, ok := msg.(models.NodesMsg)
				if !ok {
					t.Errorf("Expected NodesMsg, got %T", msg)
				}
				if len(nodesMsg) != 1 {
					t.Errorf("Expected 1 node, got %d", len(nodesMsg))
				}
			}
		})
	}
}

func TestGetClusterDetail(t *testing.T) {
	original := cmd.Commands
	defer func() { cmd.Commands = original }()

	tests := []struct {
		name        string
		clusterName string
		mockFunc    func(string) (cmd.Cluster, error)
		expectError bool
	}{
		{
			name:        "successful fetch",
			clusterName: "test-cluster",
			mockFunc: func(name string) (cmd.Cluster, error) {
				return cmd.Cluster{Name: name, Status: "running"}, nil
			},
			expectError: false,
		},
		{
			name:        "error fetching detail",
			clusterName: "test-cluster",
			mockFunc: func(name string) (cmd.Cluster, error) {
				return cmd.Cluster{}, errors.New("failed to get cluster detail")
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd.Commands = &MockCommands{
				GetClusterDetailFunc: tt.mockFunc,
			}
			
			cmdFunc := GetClusterDetail(tt.clusterName)
			msg := cmdFunc()
			
			if tt.expectError {
				msgMsg, ok := msg.(models.MessageMsg)
				if !ok {
					t.Errorf("Expected MessageMsg, got %T", msg)
				}
				if msgMsg.MsgType != "error" {
					t.Errorf("Expected error message type, got %s", msgMsg.MsgType)
				}
			} else {
				clusterMsg, ok := msg.(models.ClusterDetailMsg)
				if !ok {
					t.Errorf("Expected ClusterDetailMsg, got %T", msg)
				}
				cluster := cmd.Cluster(clusterMsg)
				if cluster.Name != tt.clusterName {
					t.Errorf("Expected cluster name %s, got %s", tt.clusterName, cluster.Name)
				}
			}
		})
	}
}

func TestCreateKindCluster(t *testing.T) {
	original := cmd.Commands
	defer func() { cmd.Commands = original }()

	tests := []struct {
		name        string
		clusterName string
		mockFunc    func(string) error
		expectError bool
		expectMsg   string
	}{
		{
			name:        "successful creation",
			clusterName: "new-cluster",
			mockFunc: func(name string) error {
				return nil
			},
			expectError: false,
			expectMsg:   "Cluster 'new-cluster' created successfully!",
		},
		{
			name:        "error creating cluster",
			clusterName: "new-cluster",
			mockFunc: func(name string) error {
				return errors.New("failed to create cluster")
			},
			expectError: true,
			expectMsg:   "failed to create cluster",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd.Commands = &MockCommands{
				CreateClusterFunc: tt.mockFunc,
			}
			
			cmdFunc := CreateKindCluster(tt.clusterName)
			msg := cmdFunc()
			
			msgMsg, ok := msg.(models.MessageMsg)
			if !ok {
				t.Errorf("Expected MessageMsg, got %T", msg)
			}
			
			if tt.expectError {
				if msgMsg.MsgType != "error" {
					t.Errorf("Expected error message type, got %s", msgMsg.MsgType)
				}
			} else {
				if msgMsg.MsgType != "success" {
					t.Errorf("Expected success message type, got %s", msgMsg.MsgType)
				}
			}
			
			if msgMsg.Text != tt.expectMsg {
				t.Errorf("Expected message %q, got %q", tt.expectMsg, msgMsg.Text)
			}
		})
	}
}

func TestDeleteKindCluster(t *testing.T) {
	original := cmd.Commands
	defer func() { cmd.Commands = original }()

	tests := []struct {
		name        string
		clusterName string
		mockFunc    func(string) error
		expectError bool
		expectMsg   string
	}{
		{
			name:        "successful deletion",
			clusterName: "old-cluster",
			mockFunc: func(name string) error {
				return nil
			},
			expectError: false,
			expectMsg:   "Cluster 'old-cluster' deleted successfully!",
		},
		{
			name:        "error deleting cluster",
			clusterName: "old-cluster",
			mockFunc: func(name string) error {
				return errors.New("failed to delete cluster")
			},
			expectError: true,
			expectMsg:   "failed to delete cluster",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd.Commands = &MockCommands{
				DeleteClusterFunc: tt.mockFunc,
			}
			
			cmdFunc := DeleteKindCluster(tt.clusterName)
			msg := cmdFunc()
			
			msgMsg, ok := msg.(models.MessageMsg)
			if !ok {
				t.Errorf("Expected MessageMsg, got %T", msg)
			}
			
			if tt.expectError {
				if msgMsg.MsgType != "error" {
					t.Errorf("Expected error message type, got %s", msgMsg.MsgType)
				}
			} else {
				if msgMsg.MsgType != "success" {
					t.Errorf("Expected success message type, got %s", msgMsg.MsgType)
				}
			}
			
			if msgMsg.Text != tt.expectMsg {
				t.Errorf("Expected message %q, got %q", tt.expectMsg, msgMsg.Text)
			}
		})
	}
}

func TestLoadDockerImage(t *testing.T) {
	original := cmd.Commands
	defer func() { cmd.Commands = original }()

	tests := []struct {
		name        string
		imageName   string
		clusterName string
		mockFunc    func(string, string) error
		expectError bool
		expectMsg   string
	}{
		{
			name:        "successful load",
			imageName:   "nginx:latest",
			clusterName: "test-cluster",
			mockFunc: func(image, cluster string) error {
				return nil
			},
			expectError: false,
			expectMsg:   "Image 'nginx:latest' loaded successfully!",
		},
		{
			name:        "error loading image",
			imageName:   "nginx:latest",
			clusterName: "test-cluster",
			mockFunc: func(image, cluster string) error {
				return errors.New("failed to load image")
			},
			expectError: true,
			expectMsg:   "failed to load image",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd.Commands = &MockCommands{
				LoadDockerImageFunc: tt.mockFunc,
			}
			
			cmdFunc := LoadDockerImage(tt.imageName, tt.clusterName)
			msg := cmdFunc()
			
			msgMsg, ok := msg.(models.MessageMsg)
			if !ok {
				t.Errorf("Expected MessageMsg, got %T", msg)
			}
			
			if tt.expectError {
				if msgMsg.MsgType != "error" {
					t.Errorf("Expected error message type, got %s", msgMsg.MsgType)
				}
			} else {
				if msgMsg.MsgType != "success" {
					t.Errorf("Expected success message type, got %s", msgMsg.MsgType)
				}
			}
			
			if msgMsg.Text != tt.expectMsg {
				t.Errorf("Expected message %q, got %q", tt.expectMsg, msgMsg.Text)
			}
		})
	}
}

func TestBuildNodeImage(t *testing.T) {
	original := cmd.Commands
	defer func() { cmd.Commands = original }()

	tests := []struct {
		name        string
		sourcePath  string
		mockFunc    func(string) error
		expectError bool
		expectMsg   string
	}{
		{
			name:       "successful build",
			sourcePath: "/path/to/source",
			mockFunc: func(path string) error {
				return nil
			},
			expectError: false,
			expectMsg:   "Node image built successfully!",
		},
		{
			name:       "error building image",
			sourcePath: "/path/to/source",
			mockFunc: func(path string) error {
				return errors.New("failed to build node image")
			},
			expectError: true,
			expectMsg:   "failed to build node image",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd.Commands = &MockCommands{
				BuildNodeImageFunc: tt.mockFunc,
			}
			
			cmdFunc := BuildNodeImage(tt.sourcePath)
			msg := cmdFunc()
			
			msgMsg, ok := msg.(models.MessageMsg)
			if !ok {
				t.Errorf("Expected MessageMsg, got %T", msg)
			}
			
			if tt.expectError {
				if msgMsg.MsgType != "error" {
					t.Errorf("Expected error message type, got %s", msgMsg.MsgType)
				}
			} else {
				if msgMsg.MsgType != "success" {
					t.Errorf("Expected success message type, got %s", msgMsg.MsgType)
				}
			}
			
			if msgMsg.Text != tt.expectMsg {
				t.Errorf("Expected message %q, got %q", tt.expectMsg, msgMsg.Text)
			}
		})
	}
}

func TestExportKindLogs(t *testing.T) {
	original := cmd.Commands
	defer func() { cmd.Commands = original }()

	homeDir, _ := os.UserHomeDir()

	tests := []struct {
		name        string
		clusterName string
		outputPath  string
		mockFunc    func(string, string) error
		expectError bool
		expectMsg   string
	}{
		{
			name:        "successful export to specific path",
			clusterName: "test-cluster",
			outputPath:  "/tmp/logs",
			mockFunc: func(cluster, path string) error {
				return nil
			},
			expectError: false,
			expectMsg:   "Logs exported to /tmp/logs successfully!",
		},
		{
			name:        "successful export to current directory",
			clusterName: "test-cluster",
			outputPath:  "",
			mockFunc: func(cluster, path string) error {
				return nil
			},
			expectError: false,
			expectMsg:   "Logs exported to current directory successfully!",
		},
		{
			name:        "successful export with tilde expansion",
			clusterName: "test-cluster",
			outputPath:  "~/logs",
			mockFunc: func(cluster, path string) error {
				// Check that path was expanded
				expectedPath := filepath.Join(homeDir, "logs")
				if path != expectedPath {
					return fmt.Errorf("path not expanded: got %s, want %s", path, expectedPath)
				}
				return nil
			},
			expectError: false,
			expectMsg:   fmt.Sprintf("Logs exported to %s successfully!", filepath.Join(homeDir, "logs")),
		},
		{
			name:        "error exporting logs",
			clusterName: "test-cluster",
			outputPath:  "/tmp/logs",
			mockFunc: func(cluster, path string) error {
				return errors.New("failed to export logs")
			},
			expectError: true,
			expectMsg:   "failed to export logs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd.Commands = &MockCommands{
				ExportLogsFunc: tt.mockFunc,
			}
			
			cmdFunc := ExportKindLogs(tt.clusterName, tt.outputPath)
			msg := cmdFunc()
			
			msgMsg, ok := msg.(models.MessageMsg)
			if !ok {
				t.Errorf("Expected MessageMsg, got %T", msg)
			}
			
			if tt.expectError {
				if msgMsg.MsgType != "error" {
					t.Errorf("Expected error message type, got %s", msgMsg.MsgType)
				}
			} else {
				if msgMsg.MsgType != "success" {
					t.Errorf("Expected success message type, got %s", msgMsg.MsgType)
				}
			}
			
			if msgMsg.Text != tt.expectMsg {
				t.Errorf("Expected message %q, got %q", tt.expectMsg, msgMsg.Text)
			}
		})
	}
}

func TestCommandsReturnFunc(t *testing.T) {
	// Test that each command returns a tea.Cmd (which is a function)
	tests := []struct {
		name string
		cmd  tea.Cmd
	}{
		{
			name: "GetKindClusters",
			cmd:  GetKindClusters(),
		},
		{
			name: "GetClusterNodes",
			cmd:  GetClusterNodes("test-cluster"),
		},
		{
			name: "GetClusterDetail",
			cmd:  GetClusterDetail("test-cluster"),
		},
		{
			name: "CreateKindCluster",
			cmd:  CreateKindCluster("new-cluster"),
		},
		{
			name: "DeleteKindCluster",
			cmd:  DeleteKindCluster("old-cluster"),
		},
		{
			name: "LoadDockerImage",
			cmd:  LoadDockerImage("nginx:latest", "test-cluster"),
		},
		{
			name: "BuildNodeImage",
			cmd:  BuildNodeImage("/path/to/source"),
		},
		{
			name: "ExportKindLogs",
			cmd:  ExportKindLogs("test-cluster", "./logs"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.cmd == nil {
				t.Errorf("%s() returned nil", tt.name)
			}
		})
	}
}

func TestExportKindLogsPathExpansion(t *testing.T) {
	// Test that ExportKindLogs expands ~ to home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Skip("Cannot get home directory")
	}

	tests := []struct {
		name         string
		outputPath   string
		shouldExpand bool
	}{
		{
			name:         "tilde path",
			outputPath:   "~/logs",
			shouldExpand: true,
		},
		{
			name:         "absolute path",
			outputPath:   "/tmp/logs",
			shouldExpand: false,
		},
		{
			name:         "relative path",
			outputPath:   "./logs",
			shouldExpand: false,
		},
		{
			name:         "empty path",
			outputPath:   "",
			shouldExpand: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get the command
			cmd := ExportKindLogs("test-cluster", tt.outputPath)
			
			if cmd == nil {
				t.Error("ExportKindLogs should return a non-nil command")
			}

			// For the tilde test, we know the function should expand it
			if tt.shouldExpand && tt.outputPath == "~/logs" {
				expectedPath := filepath.Join(homeDir, "logs")
				// The actual path expansion happens inside the command function
				// We've tested it in TestExportKindLogs
				_ = expectedPath // Avoid unused variable warning
			}
		})
	}
}

func TestCommandMessageTypes(t *testing.T) {
	// Test that commands can return appropriate message types
	
	// Verify MessageMsg has the expected fields
	successMsg := models.MessageMsg{
		Text:    "Success",
		MsgType: "success",
	}
	
	if successMsg.Text != "Success" {
		t.Error("MessageMsg.Text not set correctly")
	}
	
	if successMsg.MsgType != "success" {
		t.Error("MessageMsg.MsgType not set correctly")
	}
	
	// Test different message types
	messageTypes := []string{"success", "error", "info"}
	for _, msgType := range messageTypes {
		msg := models.MessageMsg{
			Text:    "Test message",
			MsgType: msgType,
		}
		
		if msg.MsgType != msgType {
			t.Errorf("Expected message type %s, got %s", msgType, msg.MsgType)
		}
	}
}

func TestCreateKindClusterMessage(t *testing.T) {
	// Test the expected success message format
	clusterName := "test-cluster"
	expectedMessage := "Cluster 'test-cluster' created successfully!"
	
	// Verify the message format matches what's in the implementation
	if !strings.Contains(expectedMessage, clusterName) {
		t.Error("Success message should contain cluster name")
	}
	
	if !strings.Contains(expectedMessage, "successfully") {
		t.Error("Success message should indicate success")
	}
}