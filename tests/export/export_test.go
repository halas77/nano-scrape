package export_test

import (
	"encoding/csv"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/halas77/nano-scrape/nano"
)

// loadHTML initializes a Tag from the local test_data.html file.
func loadHTML(t *testing.T) *nano.Tag {
	t.Helper()
	path := "test_data.html"
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read test_data.html: %v", err)
	}
	tag, err := nano.InitDocument(string(data))
	if err != nil {
		t.Fatalf("Failed to init document: %v", err)
	}
	return tag
}

// TestDirectExport tests the direct serialization functions: ToJSON, ToCSV, ToMD, and their file writer counterparts.
func TestDirectExport(t *testing.T) {
	root := loadHTML(t)
	cards := root.SelectAll(".product-card")
	if len(*cards) != 3 {
		t.Fatalf("Expected 3 product cards, got %d", len(*cards))
	}

	// 1. ToJSON Validation
	jsonBytes, err := cards.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}
	var exports []nano.TagExport
	if err := json.Unmarshal(jsonBytes, &exports); err != nil {
		t.Fatalf("Failed to unmarshal ToJSON output: %v", err)
	}
	if len(exports) != 3 {
		t.Errorf("Expected 3 items in JSON, got %d", len(exports))
	}
	if exports[0].Name != "div" {
		t.Errorf("Expected first tag name 'div', got %s", exports[0].Name)
	}
	if exports[0].Class != "product-card" {
		t.Errorf("Expected class 'product-card', got %s", exports[0].Class)
	}
	if exports[0].Attributes["data-sku"] != "SKU-CHIPS-9" {
		t.Errorf("Expected data-sku 'SKU-CHIPS-9', got %s", exports[0].Attributes["data-sku"])
	}

	// 2. ToCSV Validation
	csvStr, err := cards.ToCSV()
	if err != nil {
		t.Fatalf("ToCSV failed: %v", err)
	}
	reader := csv.NewReader(strings.NewReader(csvStr))
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("Failed to parse CSV output: %v", err)
	}
	if len(records) != 4 { // header + 3 cards
		t.Errorf("Expected 4 CSV lines, got %d", len(records))
	}
	// CSV header fields
	if records[0][0] != "name" || records[0][1] != "text" || records[0][2] != "class" || records[0][3] != "id" || records[0][4] != "attributes" {
		t.Errorf("Unexpected CSV header: %v", records[0])
	}
	// First card verification
	if records[1][0] != "div" || records[1][2] != "product-card" || !strings.Contains(records[1][4], "data-sku=SKU-CHIPS-9") {
		t.Errorf("Unexpected CSV record 1: %v", records[1])
	}

	// 3. ToMD Validation
	mdStr, err := cards.ToMD()
	if err != nil {
		t.Fatalf("ToMD failed: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(mdStr), "\n")
	if len(lines) != 5 { // header, separator, 3 data rows
		t.Errorf("Expected 5 MD table rows, got %d", len(lines))
	}
	for i, line := range lines {
		if !strings.HasPrefix(line, "| ") || !strings.HasSuffix(line, " |") {
			t.Errorf("MD line %d is not properly formatted: %s", i, line)
		}
	}
	// Verify that the pipe symbol '|' inside the description text is escaped as '\|'
	if !strings.Contains(lines[2], "salted \\| freshly") {
		t.Errorf("Expected escaped '\\|' in MD output line 2, but got: %s", lines[2])
	}
	// Verify that newlines inside description are replaced by spaces in MD output
	if strings.Contains(lines[2], "\n") {
		t.Errorf("Expected newlines in description to be replaced with space in MD output line 2, but got: %s", lines[2])
	}

	// 4. File Writers Validation (WriteJSON, WriteCSV, WriteMD)
	tmpDir := t.TempDir()

	tmpJSON := filepath.Join(tmpDir, "direct_products.json")
	if err := cards.WriteJSON(tmpJSON); err != nil {
		t.Fatalf("WriteJSON failed: %v", err)
	}
	if fileBytes, err := os.ReadFile(tmpJSON); err != nil {
		t.Errorf("Failed to read created JSON file: %v", err)
	} else if len(fileBytes) == 0 {
		t.Errorf("Created JSON file is empty")
	}

	tmpCSV := filepath.Join(tmpDir, "direct_products.csv")
	if err := cards.WriteCSV(tmpCSV); err != nil {
		t.Fatalf("WriteCSV failed: %v", err)
	}
	if fileBytes, err := os.ReadFile(tmpCSV); err != nil {
		t.Errorf("Failed to read created CSV file: %v", err)
	} else if len(fileBytes) == 0 {
		t.Errorf("Created CSV file is empty")
	}

	tmpMD := filepath.Join(tmpDir, "direct_products.md")
	if err := cards.WriteMD(tmpMD); err != nil {
		t.Fatalf("WriteMD failed: %v", err)
	}
	if fileBytes, err := os.ReadFile(tmpMD); err != nil {
		t.Errorf("Failed to read created MD file: %v", err)
	} else if len(fileBytes) == 0 {
		t.Errorf("Created MD file is empty")
	}
}

// TestMappedExport tests structure mapping using Tags.Map and serialization via ExportJSON, ExportCSV, ExportMD and mapped file writers.
func TestMappedExport(t *testing.T) {
	root := loadHTML(t)
	cards := root.SelectAll(".product-card")

	mapping := map[string]string{
		"sku":      "@data-sku",
		"name":     "h3.name",
		"price":    "span.price",
		"desc":     "p.description",
		"tracking": "a.details-link @data-tracking",
	}

	mappedData := cards.Map(mapping)
	if len(mappedData) != 3 {
		t.Fatalf("Expected 3 mapped rows, got %d", len(mappedData))
	}

	// Verify exact mapping values
	expectedData := []map[string]string{
		{
			"sku":      "SKU-CHIPS-9",
			"name":     "Crispy Potato Chips",
			"price":    "$3.99",
			"desc":     "Lightly salted | freshly packaged product.\nBest consumed within 6 months.",
			"tracking": "btn-chips-1",
		},
		{
			"sku":      "SKU-ROBOT-88",
			"name":     "Smart Robot Vacuum",
			"price":    "$299.99",
			"desc":     "Automated floor cleaner.\nHas a 3-stage cleaning system | quiet mode.",
			"tracking": "btn-robot-2",
		},
		{
			"sku":      "SKU-COFFEE-M",
			"name":     "Organic Dark Roast",
			"price":    "$14.50",
			"desc":     "Rich aroma | 100% Arabica beans.\nGrown in high-altitude volcanic soils.",
			"tracking": "btn-coffee-3",
		},
	}

	for i, expected := range expectedData {
		for key, val := range expected {
			got := strings.TrimSpace(mappedData[i][key])
			want := strings.TrimSpace(val)
			if got != want {
				t.Errorf("Row %d, Key %s: got %q, want %q", i, key, got, want)
			}
		}
	}

	// 1. ExportJSON Validation
	jsonBytes, err := nano.ExportJSON(mappedData)
	if err != nil {
		t.Fatalf("ExportJSON failed: %v", err)
	}
	var decoded []map[string]string
	if err := json.Unmarshal(jsonBytes, &decoded); err != nil {
		t.Fatalf("Failed to decode ExportJSON: %v", err)
	}
	if len(decoded) != 3 {
		t.Errorf("Expected 3 items in decoded JSON, got %d", len(decoded))
	}

	// 2. ExportCSV Validation
	csvStr, err := nano.ExportCSV(mappedData)
	if err != nil {
		t.Fatalf("ExportCSV failed: %v", err)
	}
	r := csv.NewReader(strings.NewReader(csvStr))
	records, err := r.ReadAll()
	if err != nil {
		t.Fatalf("Failed to parse ExportCSV: %v", err)
	}
	if len(records) != 4 { // header + 3 rows
		t.Errorf("Expected 4 CSV lines, got %d", len(records))
	}
	// Headers sorted alphabetically by map key
	expectedHeaders := []string{"desc", "name", "price", "sku", "tracking"}
	for i, h := range expectedHeaders {
		if records[0][i] != h {
			t.Errorf("Expected header index %d to be %q, got %q", i, h, records[0][i])
		}
	}
	if records[1][1] != "Crispy Potato Chips" {
		t.Errorf("Expected index 1,1 to be 'Crispy Potato Chips', got %q", records[1][1])
	}

	// 3. ExportMD Validation
	mdStr, err := nano.ExportMD(mappedData)
	if err != nil {
		t.Fatalf("ExportMD failed: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(mdStr), "\n")
	if len(lines) != 5 {
		t.Errorf("Expected 5 MD lines, got %d", len(lines))
	}
	// Headers sorting checked
	if !strings.Contains(lines[0], "desc | name | price | sku | tracking") {
		t.Errorf("Unexpected MD header line: %q", lines[0])
	}
	// Check that pipe '|' inside mapped text is escaped and newline replaced in table cell
	if !strings.Contains(lines[2], "salted \\| freshly") {
		t.Errorf("Expected pipe to be escaped in MD output, line 2: %s", lines[2])
	}
	if strings.Contains(lines[2], "\n") {
		t.Errorf("Expected newline to be replaced with space in MD output, line 2: %s", lines[2])
	}

	// 4. Mapped File Writers Validation (WriteMappedJSON, WriteMappedCSV, WriteMappedMD)
	tmpDir := t.TempDir()

	tmpJSON := filepath.Join(tmpDir, "mapped_products.json")
	if err := nano.WriteMappedJSON(tmpJSON, mappedData); err != nil {
		t.Fatalf("WriteMappedJSON failed: %v", err)
	}
	if fileBytes, err := os.ReadFile(tmpJSON); err != nil {
		t.Errorf("Failed to read mapped JSON file: %v", err)
	} else if len(fileBytes) == 0 {
		t.Errorf("Created mapped JSON file is empty")
	}

	tmpCSV := filepath.Join(tmpDir, "mapped_products.csv")
	if err := nano.WriteMappedCSV(tmpCSV, mappedData); err != nil {
		t.Fatalf("WriteMappedCSV failed: %v", err)
	}
	if fileBytes, err := os.ReadFile(tmpCSV); err != nil {
		t.Errorf("Failed to read mapped CSV file: %v", err)
	} else if len(fileBytes) == 0 {
		t.Errorf("Created mapped CSV file is empty")
	}

	tmpMD := filepath.Join(tmpDir, "mapped_products.md")
	if err := nano.WriteMappedMD(tmpMD, mappedData); err != nil {
		t.Fatalf("WriteMappedMD failed: %v", err)
	}
	if fileBytes, err := os.ReadFile(tmpMD); err != nil {
		t.Errorf("Failed to read mapped MD file: %v", err)
	} else if len(fileBytes) == 0 {
		t.Errorf("Created mapped MD file is empty")
	}
}

// TestEmptyCollections tests how empty collections or empty mapped slices are exported.
func TestEmptyCollections(t *testing.T) {
	emptyTags := &nano.Tags{}

	// Direct export on empty tag collection
	exports := emptyTags.Export()
	if len(exports) != 0 {
		t.Errorf("Expected empty exports slice, got %v", exports)
	}

	jsonBytes, err := emptyTags.ToJSON()
	if err != nil {
		t.Fatalf("empty ToJSON failed: %v", err)
	}
	if string(jsonBytes) != "null" {
		t.Errorf("Expected empty JSON to serialize as 'null', got %q", string(jsonBytes))
	}

	csvStr, err := emptyTags.ToCSV()
	if err != nil {
		t.Fatalf("empty ToCSV failed: %v", err)
	}
	if csvStr != "" {
		t.Errorf("Expected empty CSV string, got %q", csvStr)
	}

	mdStr, err := emptyTags.ToMD()
	if err != nil {
		t.Fatalf("empty ToMD failed: %v", err)
	}
	if mdStr != "" {
		t.Errorf("Expected empty MD string, got %q", mdStr)
	}

	// Mapped export on empty mapped data
	var emptyMapped []map[string]string

	jsonBytes2, err := nano.ExportJSON(emptyMapped)
	if err != nil {
		t.Fatalf("empty ExportJSON failed: %v", err)
	}
	if string(jsonBytes2) != "null" {
		t.Errorf("Expected empty mapped JSON to serialize as 'null', got %q", string(jsonBytes2))
	}

	csvStr2, err := nano.ExportCSV(emptyMapped)
	if err != nil {
		t.Fatalf("empty ExportCSV failed: %v", err)
	}
	if csvStr2 != "" {
		t.Errorf("Expected empty mapped CSV, got %q", csvStr2)
	}

	mdStr2, err := nano.ExportMD(emptyMapped)
	if err != nil {
		t.Fatalf("empty ExportMD failed: %v", err)
	}
	if mdStr2 != "" {
		t.Errorf("Expected empty mapped MD, got %q", mdStr2)
	}
}
