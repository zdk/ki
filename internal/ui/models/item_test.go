package models

import "testing"

func TestNewItem(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		description string
		action      string
	}{
		{
			name:        "create menu item",
			title:       "Create Cluster",
			description: "Create a new KIND cluster",
			action:      "create",
		},
		{
			name:        "list clusters item",
			title:       "List Clusters",
			description: "View and delete KIND clusters",
			action:      "clusters",
		},
		{
			name:        "empty values",
			title:       "",
			description: "",
			action:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := NewItem(tt.title, tt.description, tt.action)
			
			if item.Title() != tt.title {
				t.Errorf("NewItem().Title() = %v, want %v", item.Title(), tt.title)
			}
			if item.Description() != tt.description {
				t.Errorf("NewItem().Description() = %v, want %v", item.Description(), tt.description)
			}
			if item.Action != tt.action {
				t.Errorf("NewItem().Action = %v, want %v", item.Action, tt.action)
			}
			if item.FilterValue() != tt.title {
				t.Errorf("NewItem().FilterValue() = %v, want %v", item.FilterValue(), tt.title)
			}
		})
	}
}

func TestItemMethods(t *testing.T) {
	item := Item{
		title:       "Test Title",
		description: "Test Description",
		Action:      "test-action",
	}

	t.Run("Title method", func(t *testing.T) {
		expected := "Test Title"
		if got := item.Title(); got != expected {
			t.Errorf("Item.Title() = %v, want %v", got, expected)
		}
	})

	t.Run("Description method", func(t *testing.T) {
		expected := "Test Description"
		if got := item.Description(); got != expected {
			t.Errorf("Item.Description() = %v, want %v", got, expected)
		}
	})

	t.Run("FilterValue method", func(t *testing.T) {
		expected := "Test Title"
		if got := item.FilterValue(); got != expected {
			t.Errorf("Item.FilterValue() = %v, want %v", got, expected)
		}
	})

	t.Run("Action field", func(t *testing.T) {
		expected := "test-action"
		if got := item.Action; got != expected {
			t.Errorf("Item.Action = %v, want %v", got, expected)
		}
	})
}