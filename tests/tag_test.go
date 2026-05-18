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

// BenchmarkFind benchmarks the Find method
func BenchmarkFind(b *testing.B) {
	// scrape, err := engine.InitDocument(generateMockHTML(50))

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	reader := strings.NewReader(generateMockHTML(50))
	node, _ := html.Parse(reader)

	var scrape = engine.InitTag(node)

	// // 1. Setup the data needed for the test
	// // name := "main"
	// // params := map[string]any{"class": "main-panel"}

	// callback := func(foundTag engine.Tag) {
	// 	// Keep the callback minimal so we bench the method, not the callback logic
	// }

	// 2. Reset the timer to exclude the setup time above

	// b.StartTimer() // locate the position we will start the timer for the bench test
	// 3. Run the actual loop
	for b.Loop() {
		// scrape.Find(name, nil, callback)
		scrape.Fi(func(t engine.Tag) {
			//
		})
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
func BenchmarkFind_NilParams(b *testing.B) {
	t := Tag{limit: 10}
	name := "test-tag"

	for b.Loop() {
		t.Find(name, nil, func(foundTag Tag) {})
	}
}
*/
