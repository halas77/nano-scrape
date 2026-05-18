package tests

import (
	"testing"

	"github.com/halas77/goscrape/engine"
)

// A sample HTML content to use in benchmarks.
const benchmarkHTML = `
<!DOCTYPE html>
<html>
<head>
	<title>Benchmark Page</title>
</head>
<body>
	<div class="container">
		<header>
			<h1>Welcome to GoScrape Benchmark</h1>
		</header>
		<main>
			<p class="intro">This is a paragraph used to benchmark the performance of the library.</p>
			<ul class="item-list">
				<li class="item" id="item-1">First item</li>
				<li class="item" id="item-2">Second item</li>
				<li class="item" id="item-3">Third item</li>
			</ul>
		</main>
		<footer>
			<p>Footer content here.</p>
		</footer>
	</div>
</body>
</html>
`

// BenchmarkInitDocument benchmarks how long it takes to parse an HTML document.
func BenchmarkInitDocument(b *testing.B) {
	// Enable memory allocation reporting to see how many bytes and allocations are made.
	b.ReportAllocs()

	// b.N is dynamically set by Go's testing tool. The loop will run b.N times
	// to get an accurate measurement of the function's execution time.
	for i := 0; i < b.N; i++ {
		_, err := engine.InitDocument(benchmarkHTML)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkTagSelect benchmarks how fast we can select elements using CSS selectors.
func BenchmarkTagSelect(b *testing.B) {
	// 1. Setup: Parse the document once before starting the benchmark.
	doc, err := engine.InitDocument(benchmarkHTML)
	if err != nil {
		b.Fatal(err)
	}

	// 2. Reset the timer to exclude the setup time (parsing) from the benchmark results.
	b.ResetTimer()
	b.ReportAllocs()

	// 3. Run the selection operation b.N times.
	for i := 0; i < b.N; i++ {
		// Benchmark selecting all list items with class "item"
		results := doc.Select("li.item")
		if len(results) == 0 {
			b.Fatal("expected to find items")
		}
	}
}
