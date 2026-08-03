package engine

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
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
			main:          "nano-scrape is Cool",
			target:        "cool",
			caseSensitive: false,
			expected:      true,
		},
		{
			name:          "Case sensitive fail",
			main:          "nano-scrape",
			target:        "Nano-Scrape",
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

	items := doc.SelectAll(".item")
	if len(*items) != 2 {
		t.Fatalf("Expected 2 items, got %d", len(*items))
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

func TestClient_ProxyRotator_NilClient(t *testing.T) {
	var c *Client
	err := c.ProxyRotator("socks5://127.0.0.1:1080")
	if err == nil {
		t.Fatal("expected error when client is nil, got nil")
	}
}

func TestClient_ProxyRotator_Rotation(t *testing.T) {
	// 1. Setup target HTTP server
	targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))
	defer targetServer.Close()

	// 2. Start two distinct local SOCKS5 proxies
	var proxy1Hits, proxy2Hits int32
	p1Addr, cleanup1 := startMockSocks5Server(t, &proxy1Hits)
	defer cleanup1()

	p2Addr, cleanup2 := startMockSocks5Server(t, &proxy2Hits)
	defer cleanup2()

	// 3. Initialize your client and configure ProxyRotator
	c := &Client{
		client: &http.Client{},
	}

	if err := c.ProxyRotator(p1Addr, p2Addr); err != nil {
		t.Fatalf("failed to configure ProxyRotator: %v", err)
	}

	// 4. Send multiple requests and check hit counts
	totalRequests := 4
	for i := 0; i < totalRequests; i++ {
		resp, err := c.client.Get(targetServer.URL)
		if err != nil {
			t.Fatalf("request %d failed: %v", i, err)
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}

	// 5. Assert traffic was routed across both proxies
	p1Total := atomic.LoadInt32(&proxy1Hits)
	p2Total := atomic.LoadInt32(&proxy2Hits)

	if p1Total == 0 || p2Total == 0 {
		t.Errorf("expected traffic on both proxies, got Proxy 1: %d, Proxy 2: %d", p1Total, p2Total)
	}

	if p1Total+p2Total < int32(totalRequests) {
		t.Errorf("expected at least %d total proxy connections, got %d", totalRequests, p1Total+p2Total)
	}
}
