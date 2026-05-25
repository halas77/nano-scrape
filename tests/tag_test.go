package engine

import (
	"fmt"
	"strings"
	"testing"

	"github.com/halas77/goscrape/engine"
	"golang.org/x/net/html"
)

func generateMockHTML(count int) string {
	html := "<html><body>"
	for i := range count {
		html += fmt.Sprintf(`<div class="category-item">Category %d</div>`, i)
	}
	html += "</body></html>"
	return html
}

func FindMatchingAttributes(attrs []*engine.Attribute, elementAttrs []html.Attribute) bool {
	lookup := make(map[string]string)
	for _, attr := range attrs {
		normalizedKey := strings.ToLower(attr.Key)
		lookup[normalizedKey] = attr.Value
	}

	var attrsLength uint8 = uint8(len(attrs))
	var counter uint8 = 0

	for _, attr := range elementAttrs {
		normalizedKey := strings.ToLower(attr.Key)
		if valA, found := lookup[normalizedKey]; found {
			if normalizedKey == "class" && engine.FlexSearch(valA, attr.Val, false) {
				// if some one flexi Match to check the value and increment counter if it is true and continue
				fmt.Println("---------- Hello ----------- ")
				counter++
				continue
			}

			if normalizedKey == "string" {
				// use flex for equality and increment counter if it is true true and continue
			}

			if valA == attr.Val {
				counter++
			}
		}
	}

	return attrsLength == counter
}

func nameSelector(n *html.Node, name string, attrs []*engine.Attribute) bool {
	if n.Type == html.ElementNode && n.Data == name {
		return FindMatchingAttributes(attrs, n.Attr)
	}

	return false
}

func traverse(n *html.Node, name string, attrs []*engine.Attribute) {

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if nameSelector(c, name, attrs) {
			// fmt.Println("Found Match ", c.Data)
		}

		traverse(c, name, attrs)
	}
}

func TestFindMatchingAttributes(t *testing.T) {
	tests := []struct {
		name         string
		attr         []*engine.Attribute
		elementAttrs []html.Attribute
		expected     bool
	}{
		{
			name: "returns false when all attributes value does not match",
			attr: []*engine.Attribute{
				{Key: "id", Value: "submit-btn"},
				{Key: "class", Value: "btn-primary"},
			},
			elementAttrs: []html.Attribute{
				{Key: "class", Val: "btn-primary"},
				{Key: "id", Val: "submit-button"},
			},
			expected: false,
		},
		{
			name: "returns true when all attributes value match",
			attr: []*engine.Attribute{
				{Key: "id", Value: "submit-btn"},
				{Key: "class", Value: "btn-primary"},
			},
			elementAttrs: []html.Attribute{
				{Key: "class", Val: "btn btn-primary"},
				{Key: "id", Val: "submit-btn"},
			},
			expected: true,
		},
	}

	t.Run("Is Flexi work ", func(t *testing.T) {
		result := engine.FlexSearch("btn btn-primary", "btn-primary", false)

		if result != true {
			t.Errorf("expected %v but found %v", true, result)
		}
	})

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := FindMatchingAttributes(tc.attr, tc.elementAttrs)
			if result != tc.expected {
				t.Errorf("expected %v but found %v", tc.expected, result)
			}
		})
	}
}

func BenchmarkMiniTest(b *testing.B) {
	reader := strings.NewReader(generateMockHTML(50))
	node, err := html.Parse(reader)

	if err != nil {
		fmt.Println(err)
		return
	}

	b.StartTimer() // locate the position we will start the timer for the bench test
	// 3. Run the actual loop

	attr := []*engine.Attribute{
		{Key: "id", Value: "footer-id"},
	}

	for b.Loop() {
		traverse(node, "span", attr)
	}
}

// BenchmarkFind benchmarks the Find method
func BenchmarkFind(b *testing.B) {
	scrape, err := engine.InitDocument(generateMockHTML(50))

	if err != nil {
		fmt.Println(err)
		return
	}

	// 1. Setup the data needed for the test
	name := "main"
	params := map[string]any{"class": "main-panel"}
	// 2. Reset the timer to exclude the setup time above

	b.StartTimer() // locate the position we will start the timer for the bench test
	// 3. Run the actual loop

	attr := []*engine.Attribute{
		{Key: "id", Value: "footer-id"},
	}

	for b.Loop() {
		scrape.Find(name, params, func(foundTag engine.Tag) {
			// Keep the callback minimal so we bench the method, not the callback logic
		})
	}
}

func BenchmarkEmpty(b *testing.B) {

	b.StartTimer() // locate the position we will start the timer for the bench test
	// 3. Run the actual loop
	for b.Loop() {
		for range 1 {
			// fmt.Println("i:", i)
		}
	}
}

// BenchmarkFindAll benchmarks the FindAll method
/*func BenchmarkFindAll(b *testing.B) {
	t := Tag{limit: 10}
	name := "test-tag"
	params := map[string]any{"status": "active"}

	for b.Loop() {
		_ = t.FindAll(name, params)
	}
}

// BenchmarkFind_NilParams benchmarks how the code handles nil map inputs
func BenchmarkFindFirst(b *testing.B) {
	scrape, err := engine.InitDocument(generateMockHTML(50))
	if err != nil {
		return
	}

	name := "span"
	params := map[string]any{"id": "header-id"}

	b.StartTimer()
	for b.Loop() {
		t.Find(name, nil, func(foundTag Tag) {})
	}
}

// Test select functionality
func BenchmarkSelectOne(b *testing.B) {
	scrape, err := engine.InitDocument(generateMockHTML(50))
	if err != nil {
		return
	}

	selector := "#header-id"
	// params := map[string]any{"id": "footer-id"}

	b.StartTimer()
	for b.Loop() {
		scrape.SelectOne(selector)
	}
}

func BenchmarkSelectAll(b *testing.B) {
	scrape, err := engine.InitDocument(generateMockHTML(50))
	if err != nil {
		return
	}

	selector := "#footer-id"
	// params := map[string]any{"id": "footer-id"}

	b.StartTimer()
	for b.Loop() {
		scrape.Select(selector)
	}
}
