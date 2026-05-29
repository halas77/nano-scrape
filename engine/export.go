package engine

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

// TagExport represents a cleanly structured Tag for direct serialization.
type TagExport struct {
	Name       string            `json:"name"`
	Text       string            `json:"text"`
	Class      string            `json:"class,omitempty"`
	ID         string            `json:"id,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

// Export converts a single Tag into a clean, serializable TagExport structure.
func (t Tag) Export() TagExport {
	attrs := make(map[string]string)
	var class, id string
	for _, attr := range t.Attrs {
		attrs[attr.Key] = attr.Val
		if attr.Key == "class" {
			class = attr.Val
		} else if attr.Key == "id" {
			id = attr.Val
		}
	}
	return TagExport{
		Name:       t.Name,
		Text:       t.Text(),
		Class:      class,
		ID:         id,
		Attributes: attrs,
	}
}

// Export converts a Tags collection into a slice of serializable TagExport structures.
func (ts Tags) Export() []TagExport {
	var result []TagExport
	for _, t := range ts {
		result = append(result, t.Export())
	}
	return result
}

// extractValue finds the element matched by the selector and extracts either
// its text or its attribute value if the "@attrName" suffix is specified.
//
// Examples:
// - "h2.product-title" -> extracts inner text of the matching h2
// - "a.link @href"     -> extracts "href" attribute of the matching anchor tag
// - "@data-id"          -> extracts "data-id" attribute of the current tag itself
func (t Tag) extractValue(selector string) string {
	parts := strings.SplitN(selector, "@", 2)
	sel := strings.TrimSpace(parts[0])

	var match Tag
	if sel == "" {
		match = t
	} else {
		match = t.SelectOne(sel)
	}

	if match.root == nil {
		return ""
	}

	if len(parts) == 2 {
		attrName := strings.TrimSpace(parts[1])
		for _, attr := range match.Attrs {
			if attr.Key == attrName {
				return attr.Val
			}
		}
		return ""
	}

	return match.Text()
}

// Map extracts structured data from each Tag in the collection using a key-to-selector mapping.
// Key-to-selector mappings can specify tag and class names as well as attribute extractions with '@'.
//
// Example:
//
//	mapping := map[string]string{
//	    "title": "h2.product-title",
//	    "price": "span.price-tag",
//	    "url":   "a.product-link @href",
//	}
//	data := cards.Map(mapping)
func (ts Tags) Map(mapping map[string]string) []map[string]string {
	var result []map[string]string
	for _, t := range ts {
		row := make(map[string]string)
		for key, selector := range mapping {
			row[key] = t.extractValue(selector)
		}
		result = append(result, row)
	}
	return result
}

// ToJSON converts the Tags collection directly into a pretty-printed JSON byte slice.
func (ts Tags) ToJSON() ([]byte, error) {
	return json.MarshalIndent(ts.Export(), "", "  ")
}

// ToCSV converts the Tags collection directly into a CSV representation (as a string).
func (ts Tags) ToCSV() (string, error) {
	exports := ts.Export()
	if len(exports) == 0 {
		return "", nil
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	headers := []string{"name", "text", "class", "id", "attributes"}
	if err := writer.Write(headers); err != nil {
		return "", err
	}

	// Write records
	for _, exp := range exports {
		var attrPairs []string
		var keys []string
		for k := range exp.Attributes {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			attrPairs = append(attrPairs, fmt.Sprintf("%s=%s", k, exp.Attributes[k]))
		}
		attrsStr := strings.Join(attrPairs, "; ")

		record := []string{
			exp.Name,
			exp.Text,
			exp.Class,
			exp.ID,
			attrsStr,
		}
		if err := writer.Write(record); err != nil {
			return "", err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// WriteJSON writes the Tags collection directly into a file formatted as JSON.
func (ts Tags) WriteJSON(filename string) error {
	data, err := ts.ToJSON()
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// WriteCSV writes the Tags collection directly into a file formatted as CSV.
func (ts Tags) WriteCSV(filename string) error {
	csvStr, err := ts.ToCSV()
	if err != nil {
		return err
	}
	return os.WriteFile(filename, []byte(csvStr), 0644)
}

// ToMD converts the Tags collection directly into a Markdown table representation (as a string).
func (ts Tags) ToMD() (string, error) {
	exports := ts.Export()
	if len(exports) == 0 {
		return "", nil
	}

	var buf strings.Builder
	headers := []string{"name", "text", "class", "id", "attributes"}

	// Write header
	buf.WriteString("| " + strings.Join(headers, " | ") + " |\n")
	buf.WriteString("| " + strings.Repeat("--- | ", len(headers)-1) + "--- |\n")

	// Write records
	for _, exp := range exports {
		var attrPairs []string
		var keys []string
		for k := range exp.Attributes {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			val := strings.ReplaceAll(exp.Attributes[k], "|", "\\|")
			attrPairs = append(attrPairs, fmt.Sprintf("%s=%s", k, val))
		}
		attrsStr := strings.Join(attrPairs, "; ")

		row := []string{
			strings.ReplaceAll(exp.Name, "|", "\\|"),
			strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(exp.Text, "|", "\\|"), "\n", " "), "\r", ""),
			strings.ReplaceAll(exp.Class, "|", "\\|"),
			strings.ReplaceAll(exp.ID, "|", "\\|"),
			strings.ReplaceAll(attrsStr, "|", "\\|"),
		}
		buf.WriteString("| " + strings.Join(row, " | ") + " |\n")
	}

	return buf.String(), nil
}

// WriteMD writes the Tags collection directly into a file formatted as Markdown.
func (ts Tags) WriteMD(filename string) error {
	mdStr, err := ts.ToMD()
	if err != nil {
		return err
	}
	return os.WriteFile(filename, []byte(mdStr), 0644)
}

// ExportJSON converts a mapped slice of string-to-string maps into a pretty-printed JSON byte slice.
func ExportJSON(data []map[string]string) ([]byte, error) {
	return json.MarshalIndent(data, "", "  ")
}

// ExportCSV converts a mapped slice of string-to-string maps into a CSV representation (as a string).
func ExportCSV(data []map[string]string) (string, error) {
	if len(data) == 0 {
		return "", nil
	}

	// Gather unique headers and sort them for deterministic order
	headerMap := make(map[string]bool)
	for _, row := range data {
		for k := range row {
			headerMap[k] = true
		}
	}

	var headers []string
	for k := range headerMap {
		headers = append(headers, k)
	}
	sort.Strings(headers)

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	if err := writer.Write(headers); err != nil {
		return "", err
	}

	// Write rows
	for _, row := range data {
		record := make([]string, len(headers))
		for i, h := range headers {
			record[i] = row[h]
		}
		if err := writer.Write(record); err != nil {
			return "", err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// ExportMD converts a mapped slice of string-to-string maps into a Markdown table representation (as a string).
func ExportMD(data []map[string]string) (string, error) {
	if len(data) == 0 {
		return "", nil
	}

	// Gather unique headers and sort them for deterministic order
	headerMap := make(map[string]bool)
	for _, row := range data {
		for k := range row {
			headerMap[k] = true
		}
	}

	var headers []string
	for k := range headerMap {
		headers = append(headers, k)
	}
	sort.Strings(headers)

	var buf strings.Builder

	// Write header
	buf.WriteString("| " + strings.Join(headers, " | ") + " |\n")
	buf.WriteString("| " + strings.Repeat("--- | ", len(headers)-1) + "--- |\n")

	// Write rows
	for _, row := range data {
		record := make([]string, len(headers))
		for i, h := range headers {
			val := strings.ReplaceAll(row[h], "|", "\\|")
			// Also replace newlines to avoid breaking the table
			val = strings.ReplaceAll(val, "\n", " ")
			val = strings.ReplaceAll(val, "\r", "")
			record[i] = val
		}
		buf.WriteString("| " + strings.Join(record, " | ") + " |\n")
	}

	return buf.String(), nil
}

// WriteMappedJSON saves mapped structured data directly to a JSON file.
func WriteMappedJSON(filename string, data []map[string]string) error {
	jsonBytes, err := ExportJSON(data)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, jsonBytes, 0644)
}

// WriteMappedCSV saves mapped structured data directly to a CSV file.
func WriteMappedCSV(filename string, data []map[string]string) error {
	csvStr, err := ExportCSV(data)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, []byte(csvStr), 0644)
}

// WriteMappedMD saves mapped structured data directly to a Markdown file.
func WriteMappedMD(filename string, data []map[string]string) error {
	mdStr, err := ExportMD(data)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, []byte(mdStr), 0644)
}
