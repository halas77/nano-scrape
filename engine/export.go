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

type TagExport struct {
	Name       string            `json:"name"`
	Text       string            `json:"text"`
	Class      string            `json:"class,omitempty"`
	ID         string            `json:"id,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

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

func (ts *Tags) Export() []TagExport {
	var result []TagExport
	for _, t := range *ts {
		result = append(result, t.Export())
	}
	return result
}

func (t Tag) extractValue(selector string) string {
	parts := strings.SplitN(selector, "@", 2)
	sel := strings.TrimSpace(parts[0])

	var match *Tag
	if sel == "" {
		match = &t
	} else {
		match = t.SelectFirst(sel)
	}

	if match == nil || match.root == nil {
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

func (ts *Tags) Map(mapping map[string]string) []map[string]string {
	var result []map[string]string
	for _, t := range *ts {
		row := make(map[string]string)
		for key, selector := range mapping {
			row[key] = t.extractValue(selector)
		}
		result = append(result, row)
	}
	return result
}

func (ts *Tags) ToJSON() ([]byte, error) {
	return json.MarshalIndent(ts.Export(), "", "  ")
}

func (ts *Tags) ToCSV() (string, error) {
	exports := ts.Export()
	if len(exports) == 0 {
		return "", nil
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	headers := []string{"name", "text", "class", "id", "attributes"}
	if err := writer.Write(headers); err != nil {
		return "", err
	}

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

func (ts *Tags) WriteJSON(filename string) error {
	data, err := ts.ToJSON()
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func (ts *Tags) WriteCSV(filename string) error {
	csvStr, err := ts.ToCSV()
	if err != nil {
		return err
	}
	return os.WriteFile(filename, []byte(csvStr), 0644)
}

func (ts *Tags) ToMD() (string, error) {
	exports := ts.Export()
	if len(exports) == 0 {
		return "", nil
	}

	var buf strings.Builder
	headers := []string{"name", "text", "class", "id", "attributes"}

	buf.WriteString("| ")
	buf.WriteString(strings.Join(headers, " | "))
	buf.WriteString(" |\n")
	buf.WriteString("| ")
	buf.WriteString(strings.Repeat("--- | ", len(headers)-1))
	buf.WriteString("--- |\n")

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
		buf.WriteString("| ")
		buf.WriteString(strings.Join(row, " | "))
		buf.WriteString(" |\n")
	}

	return buf.String(), nil
}

func (ts *Tags) WriteMD(filename string) error {
	mdStr, err := ts.ToMD()
	if err != nil {
		return err
	}
	return os.WriteFile(filename, []byte(mdStr), 0644)
}

func ExportJSON(data []map[string]string) ([]byte, error) {
	return json.MarshalIndent(data, "", "  ")
}

func ExportCSV(data []map[string]string) (string, error) {
	if len(data) == 0 {
		return "", nil
	}

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

	if err := writer.Write(headers); err != nil {
		return "", err
	}

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

func ExportMD(data []map[string]string) (string, error) {
	if len(data) == 0 {
		return "", nil
	}

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

	buf.WriteString("| ")
	buf.WriteString(strings.Join(headers, " | "))
	buf.WriteString(" |\n")
	buf.WriteString("| ")
	buf.WriteString(strings.Repeat("--- | ", len(headers)-1))
	buf.WriteString("--- |\n")

	for _, row := range data {
		record := make([]string, len(headers))
		for i, h := range headers {
			val := strings.ReplaceAll(row[h], "|", "\\|")
			val = strings.ReplaceAll(val, "\n", " ")
			val = strings.ReplaceAll(val, "\r", "")
			record[i] = val
		}
		buf.WriteString("| ")
		buf.WriteString(strings.Join(record, " | "))
		buf.WriteString(" |\n")
	}

	return buf.String(), nil
}

func WriteMappedJSON(filename string, data []map[string]string) error {
	jsonBytes, err := ExportJSON(data)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, jsonBytes, 0644)
}

func WriteMappedCSV(filename string, data []map[string]string) error {
	csvStr, err := ExportCSV(data)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, []byte(csvStr), 0644)
}

func WriteMappedMD(filename string, data []map[string]string) error {
	mdStr, err := ExportMD(data)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, []byte(mdStr), 0644)
}
