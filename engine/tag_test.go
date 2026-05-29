package engine

import (
	"fmt"
	"testing"

	"github.com/andybalholm/cascadia"
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

func BenchmarkMatchingAttribute(b *testing.B) {
	scrape, _ := InitDocument(generateMockHTML(50))
	attrs := []html.Attribute{
		{Key: "class", Val: "category-item"},
	}

	attr := []*Attribute{
		{Key: "class", Value: "category-item"},
	}

	scrape.attrs = attr

	b.StartTimer()
	for b.Loop() {
		scrape.FindMatchingAttributes(attrs)
	}
}

func BenchmarkCascadiaQuery(b *testing.B) {
	scrape, _ := InitDocument(generateMockHTML(50))
	sel, _ := cascadia.Parse(".category-item")

	b.StartTimer()

	for b.Loop() {
		cascadia.QueryAll(scrape.root, sel)
	}
}

func BenchmarkTraverseTest(b *testing.B) {
	scrape, err := InitDocument(generateMockHTML(50))

	if err != nil {
		fmt.Println(err)
		return
	}
	attr := []*Attribute{
		{Key: "class", Value: "category-item"},
	}

	scrape.recurse = true
	scrape.attrs = attr

	b.StartTimer()

	for b.Loop() {
		scrape.traverse(scrape.root, func(t *html.Node) bool {
			if scrape.FindMatchingAttributes(t.Attr) {
				// fmt.Println("attr ", t.Attr)
			}
			return false
		})
	}
}

// BenchmarkFind benchmarks the Find method
func BenchmarkFind(b *testing.B) {
	scrape, err := InitDocument(generateMockHTML(50))

	if err != nil {
		fmt.Println(err)
		return
	}

	// 1. Setup the data needed for the test
	name := "div"

	params := []*Attribute{
		// {
		// 	Key:   "id",
		// 	Value: "footer-id",
		// },
		{
			Key:   "class",
			Value: "category-item",
		},
	}

	// 2. Reset the timer to exclude the setup time above

	b.StartTimer() // locate the position we will start the timer for the bench test
	// 3. Run the actual loop
	for b.Loop() {
		scrape.Find(name, params, func(foundTag *Tag) {
			// Keep the callback minimal so we bench the method, not the callback logic
		})
	}
}

// BenchmarkFindAll benchmarks the FindAll method
func BenchmarkFindAll(b *testing.B) {
	scrape, err := InitDocument(generateMockHTML(50))
	if err != nil {
		return
	}

	name := "div"
	params := []*Attribute{
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
		_ = scrape.FindAll(name, params)
	}
}

func BenchmarkSelectAll(b *testing.B) {
	scrape, err := InitDocument(generateMockHTML(50))
	if err != nil {
		return
	}

	selector := ".category-item"
	// params := map[string]any{"id": "footer-id"}

	b.StartTimer()
	for b.Loop() {
		// scrape.Select(selector)
		scrape.QueryAll(selector)
	}
}

func BenchmarkFindFirst(b *testing.B) {
	scrape, err := InitDocument(generateMockHTML(50))
	if err != nil {
		return
	}

	name := "span"
	params := map[string]any{"id": "header-id"}

	b.StartTimer()
	for b.Loop() {
		scrape.FindFirst(name, params)
	}
}

// Test select functionality
func BenchmarkSelectOne(b *testing.B) {
	scrape, err := InitDocument(generateMockHTML(50))
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
