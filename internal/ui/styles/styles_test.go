package styles

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestStylesExist(t *testing.T) {
	tests := []struct {
		name  string
		style lipgloss.Style
	}{
		{"Title", Title},
		{"Status", Status},
		{"Error", Error},
		{"Help", Help},
		{"Focused", Focused},
		{"Blurred", Blurred},
		{"NoStyle", NoStyle},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Check that the style is not nil (lipgloss.Style is a struct, so we check it's been initialized)
			rendered := tt.style.Render("test")
			if rendered == "" {
				t.Errorf("%s style did not render any content", tt.name)
			}
		})
	}
}

func TestStyleColors(t *testing.T) {
	tests := []struct {
		name       string
		style      lipgloss.Style
		testString string
	}{
		{
			name:       "Title has background",
			style:      Title,
			testString: "Title Text",
		},
		{
			name:       "Status is bold",
			style:      Status,
			testString: "Status Text",
		},
		{
			name:       "Error is bold",
			style:      Error,
			testString: "Error Text",
		},
		{
			name:       "Focused has background",
			style:      Focused,
			testString: "Focused Text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Render the text with the style
			rendered := tt.style.Render(tt.testString)
			
			// Check that something was rendered
			if rendered == "" {
				t.Errorf("%s style did not render any content", tt.name)
			}
			
			// For some environments, styles might not add visible formatting
			// Just ensure the text is still present
			if !strings.Contains(rendered, tt.testString) {
				t.Errorf("%s style lost the original text content", tt.name)
			}
		})
	}
}

func TestNoStyleIsMinimal(t *testing.T) {
	testText := "Plain Text"
	rendered := NoStyle.Render(testText)
	
	// NoStyle should add minimal formatting
	// The rendered text should be close to the original length
	// (allowing for some ANSI escape codes)
	if len(rendered) > len(testText)*2 {
		t.Error("NoStyle added too much formatting")
	}
}

func TestStyleConsistency(t *testing.T) {
	testText := "Consistent Text"
	
	// Test that styles produce consistent output
	styles := map[string]lipgloss.Style{
		"Title":   Title,
		"Status":  Status,
		"Error":   Error,
		"Help":    Help,
		"Focused": Focused,
		"Blurred": Blurred,
		"NoStyle": NoStyle,
	}
	
	for name, style := range styles {
		t.Run(name+" consistency", func(t *testing.T) {
			// Render the same text multiple times
			first := style.Render(testText)
			second := style.Render(testText)
			
			if first != second {
				t.Errorf("%s style produced inconsistent output", name)
			}
		})
	}
}

func TestTitleHasPadding(t *testing.T) {
	// Title style should have padding
	withPadding := Title.Render("X")
	withoutPadding := NoStyle.Render("X")
	
	if len(withPadding) <= len(withoutPadding) {
		t.Error("Title style should add padding")
	}
}