package models

import (
	"testing"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"ki/internal/cmd"
)

func TestModelStruct(t *testing.T) {
	// Test that Model struct can be created and has all expected fields
	model := Model{
		CurrentView:         MainMenuView,
		MainMenu:            list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
		ClusterList:         list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
		NodeList:            list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
		TextInput:           textinput.New(),
		Help:                help.New(),
		Clusters:            []cmd.Cluster{},
		CurrentCluster:      nil,
		Message:             "Test message",
		MessageType:         "success",
		ShowHelp:            false,
		Quitting:            false,
		Loading:             false,
		Width:               800,
		Height:              600,
		InputPrompt:         "Enter value:",
		InputAction:         "create",
		SelectedCluster:     "test-cluster",
		ClusterToDelete:     "old-cluster",
		DeleteConfirmChoice: 0,
	}

	// Test field access
	if model.CurrentView != MainMenuView {
		t.Errorf("Expected CurrentView to be MainMenuView, got %v", model.CurrentView)
	}
	if model.Message != "Test message" {
		t.Errorf("Expected Message to be 'Test message', got %s", model.Message)
	}
	if model.MessageType != "success" {
		t.Errorf("Expected MessageType to be 'success', got %s", model.MessageType)
	}
	if model.Width != 800 {
		t.Errorf("Expected Width to be 800, got %d", model.Width)
	}
	if model.Height != 600 {
		t.Errorf("Expected Height to be 600, got %d", model.Height)
	}
	if model.InputPrompt != "Enter value:" {
		t.Errorf("Expected InputPrompt to be 'Enter value:', got %s", model.InputPrompt)
	}
	if model.InputAction != "create" {
		t.Errorf("Expected InputAction to be 'create', got %s", model.InputAction)
	}
	if model.SelectedCluster != "test-cluster" {
		t.Errorf("Expected SelectedCluster to be 'test-cluster', got %s", model.SelectedCluster)
	}
	if model.ClusterToDelete != "old-cluster" {
		t.Errorf("Expected ClusterToDelete to be 'old-cluster', got %s", model.ClusterToDelete)
	}
	if model.DeleteConfirmChoice != 0 {
		t.Errorf("Expected DeleteConfirmChoice to be 0, got %d", model.DeleteConfirmChoice)
	}
}

func TestModelWithClusters(t *testing.T) {
	clusters := []cmd.Cluster{
		{Name: "cluster1", Status: "running"},
		{Name: "cluster2", Status: "stopped"},
	}

	model := Model{
		Clusters: clusters,
	}

	if len(model.Clusters) != 2 {
		t.Errorf("Expected 2 clusters, got %d", len(model.Clusters))
	}
	if model.Clusters[0].Name != "cluster1" {
		t.Errorf("Expected first cluster name to be 'cluster1', got %s", model.Clusters[0].Name)
	}
	if model.Clusters[1].Status != "stopped" {
		t.Errorf("Expected second cluster status to be 'stopped', got %s", model.Clusters[1].Status)
	}
}

func TestModelWithCurrentCluster(t *testing.T) {
	cluster := &cmd.Cluster{
		Name:        "current-cluster",
		Status:      "running",
		KubeVersion: "v1.25.3",
		Nodes: []cmd.Node{
			{Name: "node1", Status: "Ready", Role: "control-plane"},
		},
	}

	model := Model{
		CurrentCluster: cluster,
	}

	if model.CurrentCluster == nil {
		t.Error("Expected CurrentCluster to be set")
	}
	if model.CurrentCluster.Name != "current-cluster" {
		t.Errorf("Expected CurrentCluster name to be 'current-cluster', got %s", model.CurrentCluster.Name)
	}
	if model.CurrentCluster.KubeVersion != "v1.25.3" {
		t.Errorf("Expected CurrentCluster KubeVersion to be 'v1.25.3', got %s", model.CurrentCluster.KubeVersion)
	}
	if len(model.CurrentCluster.Nodes) != 1 {
		t.Errorf("Expected 1 node in CurrentCluster, got %d", len(model.CurrentCluster.Nodes))
	}
}

func TestModelBooleanFields(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		setValue bool
	}{
		{name: "ShowHelp true", field: "ShowHelp", setValue: true},
		{name: "ShowHelp false", field: "ShowHelp", setValue: false},
		{name: "Quitting true", field: "Quitting", setValue: true},
		{name: "Quitting false", field: "Quitting", setValue: false},
		{name: "Loading true", field: "Loading", setValue: true},
		{name: "Loading false", field: "Loading", setValue: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := Model{}
			
			switch tt.field {
			case "ShowHelp":
				model.ShowHelp = tt.setValue
				if model.ShowHelp != tt.setValue {
					t.Errorf("Expected ShowHelp to be %v, got %v", tt.setValue, model.ShowHelp)
				}
			case "Quitting":
				model.Quitting = tt.setValue
				if model.Quitting != tt.setValue {
					t.Errorf("Expected Quitting to be %v, got %v", tt.setValue, model.Quitting)
				}
			case "Loading":
				model.Loading = tt.setValue
				if model.Loading != tt.setValue {
					t.Errorf("Expected Loading to be %v, got %v", tt.setValue, model.Loading)
				}
			}
		})
	}
}

func TestModelDeleteConfirmChoice(t *testing.T) {
	tests := []struct {
		name     string
		choice   int
		expected int
	}{
		{name: "Yes choice", choice: 0, expected: 0},
		{name: "No choice", choice: 1, expected: 1},
		{name: "Invalid choice", choice: 2, expected: 2},
		{name: "Negative choice", choice: -1, expected: -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := Model{
				DeleteConfirmChoice: tt.choice,
			}
			
			if model.DeleteConfirmChoice != tt.expected {
				t.Errorf("Expected DeleteConfirmChoice to be %d, got %d", tt.expected, model.DeleteConfirmChoice)
			}
		})
	}
}

func TestModelInputFields(t *testing.T) {
	tests := []struct {
		name            string
		inputPrompt     string
		inputAction     string
		selectedCluster string
		clusterToDelete string
	}{
		{
			name:            "create cluster",
			inputPrompt:     "Enter cluster name:",
			inputAction:     "create",
			selectedCluster: "",
			clusterToDelete: "",
		},
		{
			name:            "load image",
			inputPrompt:     "Enter image name:",
			inputAction:     "load-image",
			selectedCluster: "target-cluster",
			clusterToDelete: "",
		},
		{
			name:            "delete cluster",
			inputPrompt:     "",
			inputAction:     "",
			selectedCluster: "",
			clusterToDelete: "cluster-to-remove",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := Model{
				InputPrompt:     tt.inputPrompt,
				InputAction:     tt.inputAction,
				SelectedCluster: tt.selectedCluster,
				ClusterToDelete: tt.clusterToDelete,
			}
			
			if model.InputPrompt != tt.inputPrompt {
				t.Errorf("Expected InputPrompt to be %q, got %q", tt.inputPrompt, model.InputPrompt)
			}
			if model.InputAction != tt.inputAction {
				t.Errorf("Expected InputAction to be %q, got %q", tt.inputAction, model.InputAction)
			}
			if model.SelectedCluster != tt.selectedCluster {
				t.Errorf("Expected SelectedCluster to be %q, got %q", tt.selectedCluster, model.SelectedCluster)
			}
			if model.ClusterToDelete != tt.clusterToDelete {
				t.Errorf("Expected ClusterToDelete to be %q, got %q", tt.clusterToDelete, model.ClusterToDelete)
			}
		})
	}
}

func TestModelViewModes(t *testing.T) {
	viewModes := []ViewMode{
		MainMenuView,
		ClusterListView,
		ClusterDetailView,
		NodeListView,
		CreateClusterView,
		LoadImageView,
		BuildImageView,
		ExportLogsView,
		DeleteConfirmView,
	}

	for _, viewMode := range viewModes {
		t.Run("view_mode_"+string(rune(int(viewMode)+'0')), func(t *testing.T) {
			model := Model{
				CurrentView: viewMode,
			}
			
			if model.CurrentView != viewMode {
				t.Errorf("Expected CurrentView to be %v, got %v", viewMode, model.CurrentView)
			}
		})
	}
}

func TestModelTeaModelInterface(t *testing.T) {
	// Test that Model implements tea.Model interface
	model := Model{}
	
	// Test Init method
	cmd := model.Init()
	if cmd != nil {
		t.Error("Expected Init() to return nil")
	}
	
	// Test Update method
	updatedModel, updateCmd := model.Update(nil)
	// Can't compare models directly due to internal fields, just check it returns something
	if updatedModel == nil {
		t.Error("Expected Update() to return a model")
	}
	if updateCmd != nil {
		t.Error("Expected Update() to return nil command")
	}
	
	// Test View method
	view := model.View()
	if view != "" {
		t.Error("Expected View() to return empty string")
	}
}

func TestModelAsTeaModel(t *testing.T) {
	// Test that Model can be used as tea.Model
	var teaModel tea.Model = Model{}
	
	// Test that we can call tea.Model methods
	cmd := teaModel.Init()
	if cmd != nil {
		t.Error("Expected Init() to return nil")
	}
	
	newModel, newCmd := teaModel.Update(nil)
	if newModel == nil {
		t.Error("Expected Update() to return a model")
	}
	if newCmd != nil {
		t.Error("Expected Update() to return nil command")
	}
	
	view := teaModel.View()
	if view != "" {
		t.Error("Expected View() to return empty string")
	}
}

func TestModelMessageTypes(t *testing.T) {
	messageTypes := []string{"success", "error", "info", "warning", ""}
	
	for _, msgType := range messageTypes {
		t.Run("message_type_"+msgType, func(t *testing.T) {
			model := Model{
				Message:     "Test message for " + msgType,
				MessageType: msgType,
			}
			
			if model.MessageType != msgType {
				t.Errorf("Expected MessageType to be %q, got %q", msgType, model.MessageType)
			}
		})
	}
}

func TestModelDimensions(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
	}{
		{name: "small", width: 400, height: 300},
		{name: "medium", width: 800, height: 600},
		{name: "large", width: 1920, height: 1080},
		{name: "zero", width: 0, height: 0},
		{name: "negative", width: -100, height: -50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := Model{
				Width:  tt.width,
				Height: tt.height,
			}
			
			if model.Width != tt.width {
				t.Errorf("Expected Width to be %d, got %d", tt.width, model.Width)
			}
			if model.Height != tt.height {
				t.Errorf("Expected Height to be %d, got %d", tt.height, model.Height)
			}
		})
	}
}

func TestModelZeroValue(t *testing.T) {
	// Test zero value of Model
	var model Model
	
	// Check that zero values are as expected
	if model.CurrentView != MainMenuView {
		t.Errorf("Expected zero CurrentView to be MainMenuView (0), got %v", model.CurrentView)
	}
	if model.Message != "" {
		t.Errorf("Expected zero Message to be empty, got %q", model.Message)
	}
	if model.MessageType != "" {
		t.Errorf("Expected zero MessageType to be empty, got %q", model.MessageType)
	}
	if model.ShowHelp != false {
		t.Error("Expected zero ShowHelp to be false")
	}
	if model.Quitting != false {
		t.Error("Expected zero Quitting to be false")
	}
	if model.Loading != false {
		t.Error("Expected zero Loading to be false")
	}
	if model.Width != 0 {
		t.Errorf("Expected zero Width to be 0, got %d", model.Width)
	}
	if model.Height != 0 {
		t.Errorf("Expected zero Height to be 0, got %d", model.Height)
	}
	if model.DeleteConfirmChoice != 0 {
		t.Errorf("Expected zero DeleteConfirmChoice to be 0, got %d", model.DeleteConfirmChoice)
	}
	if model.CurrentCluster != nil {
		t.Error("Expected zero CurrentCluster to be nil")
	}
	if model.Clusters != nil {
		t.Error("Expected zero Clusters to be nil")
	}
}