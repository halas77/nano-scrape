package engine

import (
	"fmt"
	"strings"
	"testing"

	"github.com/halas77/nano-scrape/engine"
	"golang.org/x/net/html"
)

func generateMockHTML(count int) string {
	html := "<html><body>"
	html += `<span class="spam-category-item" id="header-id">Header is hear</span>`
	for i := range count {
		html += fmt.Sprintf(`<div class="category-item">Category %d</div>`, i)
	}
	html += `<span class="spam-category-item" id="footer-id">Footer is hear</span>`
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
			if normalizedKey == "string" {
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

func traverse(n *html.Node, name string, attrs []*engine.Attribute, f func(*html.Node)) {

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if nameSelector(c, name, attrs) {
			f(c)
		}

		traverse(c, name, attrs, f)
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
				{Key: "class", Val: "btn-primary-2"},
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
				{Key: "class", Val: "btn-primary"},
				{Key: "id", Val: "submit-btn"},
			},
			expected: true,
		},
	}

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

	b.StartTimer()

	attr := []*engine.Attribute{
		{Key: "id", Value: "footer-id"},
	}

	for b.Loop() {
		traverse(node, "span", attr, func(n *html.Node) {

		})
	}
}

func BenchmarkFind(b *testing.B) {
	scrape, err := engine.InitDocument(generateMockHTML(50))

	if err != nil {
		fmt.Println(err)
		return
	}

	name := "div"
	params := []*engine.Attribute{
		{
			Key:   "class",
			Value: "main-panel",
		},
	}


	b.StartTimer() 
	for b.Loop() {
		scrape.Find(name, params, func(foundTag *engine.Tag) {
		})
	}
}

func BenchmarkEmpty(b *testing.B) {

	b.StartTimer() 
	for b.Loop() {
		for range 1 {
		}
	}
}

func BenchmarkFindAll(b *testing.B) {
	scrape, err := engine.InitDocument(generateMockHTML(50))
	if err != nil {
		return
	}

	name := "span"
	params := []*engine.Attribute{
		{
			Key:   "id",
			Value: "footer-id",
		},
	}

	for b.Loop() {
		_ = scrape.FindAll(name, params)
	}
}

func BenchmarkFindFirst(b *testing.B) {
	scrape, err := engine.InitDocument(generateMockHTML(50))
	if err != nil {
		return
	}

	name := "div"
	params := []*engine.Attribute{
		// {
		// 	Key:   "id",
		// 	Value: "footer-id",
		// },
		{
			Key:   "class",
			Value: "category-item",
		},
	}

	b.StartTimer()
	for b.Loop() {
		scrape.FindFirst(name, params)
	}
}

// Test select functionality
func BenchmarkSelect(b *testing.B) {
	scrape, err := engine.InitDocument(generateMockHTML(50))
	if err != nil {
		return
	}

	selector := "#header-id"
	// params := map[string]any{"id": "footer-id"}

	b.StartTimer()
	for b.Loop() {
		scrape.SelectFirst(selector)
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
		scrape.SelectAll(selector)
	}
}
