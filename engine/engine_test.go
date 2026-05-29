package engine

import (
	"fmt"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

/*
func TestGetNodeStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    string // HTML snippet
		expected string
	}{
		{
			name:     "Simple text",
			input:    "<div>Hello World</div>",
			expected: "Hello World",
		},
		{
			name:     "Ignore nested tags",
			input:    "<div>Hello <span>Inner</span> World</div>",
			expected: "Hello  World",
		},
		{
			name:     "Multiple text nodes",
			input:    "<div>Part 1Part 2</div>",
			expected: "Part 1Part 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Helper to turn string into a node
			doc, _ := html.Parse(strings.NewReader(tt.input))

			// html.Parse returns <html><body><your_input>... we want the div
			node := doc.FirstChild.LastChild.FirstChild

			got := getNodeStrings(node)
			if got != tt.expected {
				t.Errorf("getNodeStrings() = %v, want %v", got, tt.expected)
			}
		})
	}
} */

func TestHasIntersection(t *testing.T) {
	// Mock attributes
	attrs := []html.Attribute{
		{Key: "class", Val: "btn-primary"},
		{Key: "id", Val: "submit-button"},
	}

	tests := []struct {
		name     string
		params   map[string]any
		isStrict bool
		expected bool
	}{
		{
			name:     "Nil params should return true",
			params:   nil,
			expected: true,
		},
		{
			name:     "Single match (Non-Strict)",
			params:   map[string]any{"class": "btn-primary"},
			isStrict: false,
			expected: true,
		},
		{
			name:     "Full match (Strict)",
			params:   map[string]any{"class": "btn-primary", "id": "submit-button"},
			isStrict: true,
			expected: true,
		},
		{
			name:     "Partial match fails (Strict)",
			params:   map[string]any{"class": "btn-primary", "id": "wrong-id"},
			isStrict: true,
			expected: false,
		},
		{
			name:     "Ignore 'string' key in attribute count",
			params:   map[string]any{"class": "btn-primary", "string": "click me"},
			isStrict: true,
			expected: true, // Only 'class' is checked against attributes
		},
		{
			name:     "Ignore '_name_', and 'string' key in attribute count",
			params:   map[string]any{"class": "btn-primary", "string": "click me", "_name_": "div"},
			isStrict: true,
			expected: true, // Only 'class' is checked against attributes
		},
		{
			name:     "Ignore '_name_' key in attribute count",
			params:   map[string]any{"class": "btn-primary", "_name_": "div"},
			isStrict: true,
			expected: true, // Only 'class' is checked against attributes
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasIntersection(tt.params, attrs, tt.isStrict)
			if got != tt.expected {
				t.Errorf("hasIntersection() %s: got %v, want %v", tt.name, got, tt.expected)
			}
		})
	}
}

func TestFlexMatch(t *testing.T) {
	tests := []struct {
		name          string
		main          string
		target        string
		caseSensitive bool
		expected      bool
	}{
		{
			name:          "Case insensitive match",
			main:          "GoScrape is Cool",
			target:        "cool",
			caseSensitive: false,
			expected:      true,
		},
		{
			name:          "Case sensitive fail",
			main:          "GoScrape",
			target:        "goscrape",
			caseSensitive: true,
			expected:      false,
		},
		{
			name:          "Regex characters are escaped",
			main:          "Price is $100*",
			target:        "$100*",
			caseSensitive: false,
			expected:      true, // If QuoteMeta works, this passes
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := flexMatch(tt.main, tt.target, tt.caseSensitive)
			if got != tt.expected {
				t.Errorf("flexMatch() %s: got %v, want %v", tt.name, got, tt.expected)
			}
		})
	}
}

func TestFindMatchingAttributes(t *testing.T) {
	tests := []struct {
		name         string
		attr         []*Attribute
		elementAttrs []html.Attribute
		expected     bool
	}{
		{
			name: "returns false when all attributes value does not match",
			attr: []*Attribute{
				{Key: "id", Value: "submit-btn"},
				{Key: "class", Value: "btn-primary"},
			},
			elementAttrs: []html.Attribute{
				{Key: "class", Val: "btn-primary-2"},
				{Key: "id", Val: "submit-button"},
			},
			expected: false,
		},
		{
			name: "returns true when all attributes value match",
			attr: []*Attribute{
				{Key: "id", Value: "submit-btn"},
				{Key: "class", Value: "btn-primary"},
				{Key: "string", Value: "World"},
			},
			elementAttrs: []html.Attribute{
				{Key: "class", Val: "btn-primary"},
				{Key: "id", Val: "submit-btn"},
			},
			expected: true,
		},
	}

	input := `<div class="details">	Hello <span> Abebe Kebede </span> World </div>`

	reader := strings.NewReader(input)
	node, err := html.Parse(reader)

	if err != nil {
		fmt.Println("err ", err)
	}

	n := node.FirstChild.FirstChild.NextSibling.FirstChild

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tag := Tag{attrs: tc.attr}
			result := tag.FindMatchingAttributes(tc.elementAttrs, n)
			if result != tc.expected {
				t.Errorf("expected %v but found %v", tc.expected, result)
			}
		})
	}

	// 6. Test direct Tags ToMD
	tagsMD, err := items.ToMD()
	if err != nil {
		t.Errorf("items.ToMD failed: %v", err)
	}
	if !strings.Contains(tagsMD, "| name | text | class | id | attributes |") {
		t.Errorf("Tags ToMD header is incorrect: %s", tagsMD)
	}
	// Row 1 should contain specific fields but we're flexible with exact whitespace in text
	if !strings.Contains(tagsMD, "| div |") || !strings.Contains(tagsMD, "Go Programming") || !strings.Contains(tagsMD, "$39.99") || !strings.Contains(tagsMD, "| item |") || !strings.Contains(tagsMD, "data-category=books") {
		t.Errorf("Tags ToMD row 1 is incorrect: %s", tagsMD)
	}

	// 7. Test ExportMD
	mappedMD, err := ExportMD(mappedData)
	if err != nil {
		t.Errorf("ExportMD failed: %v", err)
	}
	if !strings.Contains(mappedMD, "| category | price | title |") {
		t.Errorf("ExportMD header is incorrect: %s", mappedMD)
	}
	if !strings.Contains(mappedMD, "| books | $39.99 | Go Programming |") {
		t.Errorf("ExportMD row 1 is incorrect: %s", mappedMD)
	}
}
