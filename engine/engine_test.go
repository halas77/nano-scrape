package engine

import (
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

func TestExport(t *testing.T) {
	htmlContent := `
		<div id="container">
			<div class="item" data-category="books">
				<span class="title">Go Programming</span>
				<span class="price">$39.99</span>
			</div>
			<div class="item" data-category="electronics">
				<span class="title">Mechanical Keyboard</span>
				<span class="price">$99.99</span>
			</div>
		</div>
	`

	doc, err := InitDocument(htmlContent)
	if err != nil {
		t.Fatalf("Failed to initialize document: %v", err)
	}

	items := doc.Select(".item")
	if len(items) != 2 {
		t.Fatalf("Expected 2 items, got %d", len(items))
	}

	// 1. Test Map mapping
	mapping := map[string]string{
		"category": "@data-category",
		"title":    "span.title",
		"price":    "span.price",
	}

	mappedData := items.Map(mapping)
	if len(mappedData) != 2 {
		t.Fatalf("Expected 2 mapped items, got %d", len(mappedData))
	}

	if mappedData[0]["category"] != "books" || mappedData[0]["title"] != "Go Programming" || mappedData[0]["price"] != "$39.99" {
		t.Errorf("Incorrect mapping for first item: %v", mappedData[0])
	}
	if mappedData[1]["category"] != "electronics" || mappedData[1]["title"] != "Mechanical Keyboard" || mappedData[1]["price"] != "$99.99" {
		t.Errorf("Incorrect mapping for second item: %v", mappedData[1])
	}

	// 2. Test ExportJSON
	jsonBytes, err := ExportJSON(mappedData)
	if err != nil {
		t.Errorf("ExportJSON failed: %v", err)
	}
	jsonStr := string(jsonBytes)
	if !strings.Contains(jsonStr, `"category": "books"`) || !strings.Contains(jsonStr, `"title": "Mechanical Keyboard"`) {
		t.Errorf("ExportJSON output does not contain expected fields: %s", jsonStr)
	}

	// 3. Test ExportCSV
	csvStr, err := ExportCSV(mappedData)
	if err != nil {
		t.Errorf("ExportCSV failed: %v", err)
	}
	if !strings.Contains(csvStr, "category,price,title") {
		t.Errorf("ExportCSV header is incorrect: %s", csvStr)
	}
	if !strings.Contains(csvStr, "books,$39.99,Go Programming") {
		t.Errorf("ExportCSV row 1 is incorrect: %s", csvStr)
	}

	// 4. Test direct Tags ToJSON
	tagsJSON, err := items.ToJSON()
	if err != nil {
		t.Errorf("items.ToJSON failed: %v", err)
	}
	tagsJSONStr := string(tagsJSON)
	if !strings.Contains(tagsJSONStr, `"class": "item"`) {
		t.Errorf("Tags ToJSON output does not contain expected class field: %s", tagsJSONStr)
	}

	// 5. Test direct Tags ToCSV
	tagsCSV, err := items.ToCSV()
	if err != nil {
		t.Errorf("items.ToCSV failed: %v", err)
	}
	if !strings.Contains(tagsCSV, "name,text,class,id,attributes") {
		t.Errorf("Tags ToCSV header is incorrect: %s", tagsCSV)
	}
}

