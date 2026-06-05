package engine

import (
	"fmt"
	"strings"
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

func TestFinds(t *testing.T) {
	scrape, _ := InitDocument(generateMockHTML(50))

	// attr1 := []html.Attribute{
	// 	{Key: "class", Val: "category-item"},
	// }

	attr2 := []*Attribute{
		{Value: "footer-id", Key: "id"},
	}

	t.Run("Does Find work", func(t *testing.T) {
		span := scrape.FindFirst("span", attr2)
		if !strings.Contains(span.Print(), "Footer is hear") {
			t.Errorf("FindFirst: Expected to fine Footer is hear but not exist on the string.")
		}
	})
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
		scrape.FindMatchingAttributes(attrs, nil)
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
			if scrape.FindMatchingAttributes(t.Attr, nil) {
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
	params := []*Attribute{
		{
			Key:   "id",
			Value: "header-id",
		},
		// {
		// 	Key:   "class",
		// 	Value: "category-item",
		// },
	}

	b.StartTimer()
	for b.Loop() {
		scrape.FindFirst(name, params)
	}
}

// Test select functionality
func BenchmarkSelect(b *testing.B) {
	scrape, err := InitDocument(generateMockHTML(50))
	if err != nil {
		return
	}

	selector := "#header-id"
	// params := map[string]any{"id": "footer-id"}

	b.StartTimer()
	for b.Loop() {
		scrape.QueryOne(selector)
	}
}
