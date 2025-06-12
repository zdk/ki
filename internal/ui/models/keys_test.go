package models

import (
	"testing"

	"github.com/charmbracelet/bubbles/key"
)

func TestKeyBindings(t *testing.T) {
	tests := []struct {
		name    string
		binding key.Binding
		keys    []string
		help    string
	}{
		{
			name:    "up key",
			binding: Keys.Up,
			keys:    []string{"up", "k"},
			help:    "move up",
		},
		{
			name:    "down key",
			binding: Keys.Down,
			keys:    []string{"down", "j"},
			help:    "move down",
		},
		{
			name:    "quit key",
			binding: Keys.Quit,
			keys:    []string{"q", "ctrl+c"},
			help:    "quit",
		},
		{
			name:    "enter key",
			binding: Keys.Enter,
			keys:    []string{"enter"},
			help:    "select",
		},
		{
			name:    "create key",
			binding: Keys.Create,
			keys:    []string{"c"},
			help:    "create cluster",
		},
		{
			name:    "delete key",
			binding: Keys.Delete,
			keys:    []string{"d"},
			help:    "delete cluster",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Check if keys are correctly bound
			actualKeys := tt.binding.Keys()
			if len(actualKeys) != len(tt.keys) {
				t.Errorf("Expected %d keys, got %d", len(tt.keys), len(actualKeys))
			}

			// Check help text contains expected string
			help := tt.binding.Help()
			if help.Desc != tt.help {
				t.Errorf("Expected help desc %q, got %q", tt.help, help.Desc)
			}
		})
	}
}

func TestKeyMapHelpers(t *testing.T) {
	t.Run("ShortHelp", func(t *testing.T) {
		shortHelp := Keys.ShortHelp()
		if len(shortHelp) != 2 {
			t.Errorf("Expected 2 items in short help, got %d", len(shortHelp))
		}

		// Should contain Help and Quit
		// Check by comparing help text since we can't compare structs directly
		hasHelp := false
		hasQuit := false
		for _, binding := range shortHelp {
			help := binding.Help()
			if help.Key == "?" {
				hasHelp = true
			}
			if help.Key == "q" {
				hasQuit = true
			}
		}

		if !hasHelp {
			t.Error("Short help should contain Help key")
		}
		if !hasQuit {
			t.Error("Short help should contain Quit key")
		}
	})

	t.Run("FullHelp", func(t *testing.T) {
		fullHelp := Keys.FullHelp()
		if len(fullHelp) != 4 {
			t.Errorf("Expected 4 groups in full help, got %d", len(fullHelp))
		}

		// First group should have navigation keys
		if len(fullHelp[0]) != 5 {
			t.Errorf("Expected 5 keys in navigation group, got %d", len(fullHelp[0]))
		}

		// Second group should have action keys
		if len(fullHelp[1]) != 3 {
			t.Errorf("Expected 3 keys in action group, got %d", len(fullHelp[1]))
		}
	})
}

func TestKeyConflicts(t *testing.T) {
	// Check for conflicting key bindings
	keyMap := make(map[string]string)
	
	bindings := []struct {
		name    string
		binding key.Binding
	}{
		{"Up", Keys.Up},
		{"Down", Keys.Down},
		{"Left", Keys.Left},
		{"Right", Keys.Right},
		{"Enter", Keys.Enter},
		{"Back", Keys.Back},
		{"Quit", Keys.Quit},
		{"Help", Keys.Help},
		{"Create", Keys.Create},
		{"Delete", Keys.Delete},
		{"Refresh", Keys.Refresh},
		{"Load", Keys.Load},
		{"Build", Keys.Build},
		{"Logs", Keys.Logs},
		{"Nodes", Keys.Nodes},
		{"Detail", Keys.Detail},
		{"Yes", Keys.Yes},
		{"No", Keys.No},
		{"Tab", Keys.Tab},
	}

	for _, b := range bindings {
		for _, k := range b.binding.Keys() {
			if existing, ok := keyMap[k]; ok {
				// Special cases for context-specific keys:
				// 'n' is used for both "nodes" and "no" in different contexts
				// 'l' is used for both "right/select" and "load" in different contexts
				if (k == "n" && ((existing == "Nodes" && b.name == "No") || (existing == "No" && b.name == "Nodes"))) ||
				   (k == "l" && ((existing == "Right" && b.name == "Load") || (existing == "Load" && b.name == "Right"))) {
					continue
				}
				t.Errorf("Key conflict: %q is used by both %s and %s", k, existing, b.name)
			}
			keyMap[k] = b.name
		}
	}
}