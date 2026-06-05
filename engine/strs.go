package engine

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func (t *Tag) getNodeStrings(n *html.Node) string {
	nodeStrings := []string{}

	t.traverse(n, func(c *html.Node) bool {
		if c.Type == html.TextNode {
			nodeStrings = append(nodeStrings, c.Data)
		}
		return false
	})

	return strings.Join(nodeStrings, "")
}

func getNodeStrings(t Tag) string {
	nodeStrings := []string{}

	traverse(t, 0, true, func(t Tag) bool {
		c := t.root
		if c.Type == html.TextNode {
			nodeStrings = append(nodeStrings, c.Data)
		}
		return false
	})

	return strings.Join(nodeStrings, "")
}

func printAttributes(attrs []html.Attribute) string {
	var builder strings.Builder

	var counter uint8 = 0
	for _, attr := range attrs {
		builder.WriteString(" ")

		builder.WriteString(attr.Key)
		builder.WriteString("=")
		builder.WriteString(`"`)
		builder.WriteString(attr.Val)
		builder.WriteString(`"`)
		counter++
	}

	return builder.String()
}

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
		builder.WriteString(printAttributes(node.Attr))
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

func (tag *Tag) Print(depth ...uint16) string {
	var d uint16 = 0
	if len(depth) > 0 {
		d = (depth[0])
	}

	if tag == nil || tag.root == nil {
		return "Empty"
	}

	return print(tag.root, d)
}

func (ts Tags) Print(depth ...uint16) string {
	length := len(ts)
	var builder strings.Builder

	for i, tag := range ts {
		if i > 0 {
			builder.WriteString("\n")
		}
		fmt.Fprint(&builder, i)
		builder.WriteString(": [\n")
		builder.WriteString(tag.Print(depth...))
		builder.WriteString("\n]")
		if length > (i + 1) {
			builder.WriteString(",")
		}
	}

	return builder.String()
}
