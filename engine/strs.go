package engine

import (
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func print(node *html.Node, indentWidth uint16) string {

	var builder strings.Builder
	var factor uint16 = 0
	width := indentWidth

	if node.Type != html.DocumentNode {
		factor = 3
	}

	if node.Type == html.TextNode {
		builder.WriteString(strings.Repeat(" ", int(width)))
		input := strings.TrimSpace(node.Data)

		if input == "" {
			return ""
		}
		builder.WriteString(input)
		return builder.String()
	}

	if node.Type == html.ElementNode {
		builder.WriteString(strings.Repeat(" ", int(width)))
		builder.WriteString("<")
		builder.WriteString(node.Data)
		builder.WriteString(">")
	}

	canIndent := false

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		canIndent = true

		input := print(c, width+factor)
		if c == node.FirstChild {
			builder.WriteString("\n")
		}
		if input != "" {
			builder.WriteString(input)
			builder.WriteString("\n")
		}
	}

	if node.Type == html.ElementNode {
		if canIndent {
			builder.WriteString(strings.Repeat(" ", int(width)))
		}

		builder.WriteString("</")
		builder.WriteString(node.Data)
		builder.WriteString(">")
	}

	return builder.String()
}

func (tag Tag) Print(depth ...uint16) string {
	var d uint16 = 0
	if len(depth) > 0 {
		d = (depth[0])
	}

	if tag.root == nil {
		return "Empty"
	}

	return print(tag.root, d)
}

func flexMatch(main string, target string, caseSensitive bool) bool {
	pattern := regexp.QuoteMeta(target)

	if !caseSensitive {
		pattern = "(?i)" + pattern
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}

	return re.MatchString(main)
}
