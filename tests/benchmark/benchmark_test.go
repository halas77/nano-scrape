package engine

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/halas77/goscrape/engine"
)

// COMMAND = go test -bench=Path ./tests/benchmark -benchmem -count=3 > tests/benchmark/benchmark_results.txt

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

// BenchmarkFindPath measures the performance of the Find method.
func BenchmarkFindPath(b *testing.B) {
	tag := loadBenchmarkTag(b)
	name := "span"
	attrs := []*engine.Attribute{{Key: "id", Value: "category-item"}}
	b.ResetTimer()
	for b.Loop() {
		tag.Find(name, attrs, func(t *engine.Tag) {})
	}
}

// BenchmarkFindAllPath measures the performance of retrieving many tags at once.
func BenchmarkFindAllPath(b *testing.B) {
	tag := loadBenchmarkTag(b)
	name := "div"
	attrs := []*engine.Attribute{{Key: "class", Value: "category-item"}}
	b.ResetTimer()
	for b.Loop() {
		tag.FindAll(name, attrs)
	}
}

// BenchmarkFindFirstPath measures the performance of retrieving the first matching tag.
func BenchmarkFindFirstPath(b *testing.B) {
	tag := loadBenchmarkTag(b)
	name := "div"
	attrs := []*engine.Attribute{{Key: "class", Value: "category-item"}}
	b.ResetTimer()
	for b.Loop() {
		tag.FindFirst(name, attrs)
	}
}

// BenchmarkSelectPath measures the performance of Select on a single Tag.
func BenchmarkSelectPath(b *testing.B) {
	tag := loadBenchmarkTag(b)
	selector := ".category-item"
	b.ResetTimer()
	for b.Loop() {
		tag.Select(selector)
	}
}

// BenchmarkSelectAllPath measures CSS selector performance for retrieving many tags.
func BenchmarkSelectAllPath(b *testing.B) {
	tag := loadBenchmarkTag(b)
	selector := ".category-item"
	b.ResetTimer()
	for b.Loop() {
		tag.SelectAll(selector)
	}
}
