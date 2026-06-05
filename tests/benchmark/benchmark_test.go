package engine

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/halas77/goscrape/engine"
)

// COMMAND = go test -bench=Educational ./tests -benchmem -count=3 > tests/benchmark/benchmark_results.txt

// loadBenchmarkTag initializes a Tag from our static HTML file.
func loadBenchmarkTag(b *testing.B) *engine.Tag {
	b.Helper()
	path := filepath.Join("../", "testdata", "benchmark_data.html")
	data, err := os.ReadFile(path)
	if err != nil {
		b.Fatalf("Failed to read: %v", err)
	}
	tag, err := engine.InitDocument(string(data))
	if err != nil {
		b.Fatalf("Failed to init: %v", err)
	}
	return tag
}

// BenchmarkFindEducational measures the performance of the Find method.
func BenchmarkFindEducational(b *testing.B) {
	tag := loadBenchmarkTag(b)
	name := "span"
	attrs := []*engine.Attribute{{Key: "id", Value: "footer-id"}}
	b.ResetTimer()
	for b.Loop() {
		tag.Find(name, attrs, func(t *engine.Tag) {})
	}
}

// BenchmarkSelectOneEducational measures the performance of SelectOne on a single Tag.
func BenchmarkSelectOneEducational(b *testing.B) {
	tag := loadBenchmarkTag(b)
	selector := ".category-item"
	b.ResetTimer()
	for b.Loop() {
		tag.SelectOne(selector)
	}
}

// BenchmarkFindAllEducational measures the performance of retrieving many tags at once.
func BenchmarkFindAllEducational(b *testing.B) {
	tag := loadBenchmarkTag(b)
	name := "div"
	attrs := []*engine.Attribute{{Key: "class", Value: "category-item"}}
	b.ResetTimer()
	for b.Loop() {
		tag.FindAll(name, attrs)
	}
}

// BenchmarkSelectEducational measures CSS selector performance.
func BenchmarkSelectEducational(b *testing.B) {
	tag := loadBenchmarkTag(b)
	selector := ".category-item"
	b.ResetTimer()
	for b.Loop() {
		tag.Select(selector)
	}
}
