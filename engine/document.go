package engine

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// adding strings function to get strings of a tag not sub tags concatenate them and append them to an array

type Document struct {
	Root      *html.Node
	isStrict  bool
	params    map[string]any
	limit     uint8
	isInitial bool
}

func Scrape(input string) (Document, error) {
	reader := strings.NewReader(input)
	node, err := html.Parse(reader)
	doc := Document{Root: node, isInitial: true}
	fmt.Println("is initial: ", doc.isInitial)

	if err != nil {
		return doc, err
	}

	return doc, nil
}

type Traverser interface {
	// Traverse()
	FindAll(name string)
	FindOne()
}

func getNodeStrings(n *html.Node) string {
	nodeStrings := []string{}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			nodeStrings = append(nodeStrings, c.Data)
		}
	}

	return strings.Join(nodeStrings, "")
}

func hasIntersection(params map[string]any, attributes []html.Attribute, isStrict bool) bool {
	if params == nil {
		return true
	}

	var count uint8 = 0
	length := len(params)
	_, ok := params["string"]
	_, name := params["_name_"]

	if ok || name {
		length = 0
	}

	if length == int(count) {
		return true
	}

	for _, attr := range attributes {
		value, ok := params[attr.Key]

		if ok && value == attr.Val {
			count++
		}
	}

	if count > 0 && uint8(length) == count {
		return true
	} else if !isStrict && count > 0 {
		return true
	}

	return false
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

func (d Document) traverse(n *html.Node, nodes *[]Document) uint8 {
	name, _ := d.params["_name_"]
	limit := d.limit

	if n.Type == html.ElementNode && n.Data == name {
		canAddNode := false
		canAddNode = hasIntersection(d.params, n.Attr, false)
		target, ok := d.params["string"]

		if ok {
			str := getNodeStrings(n)
			fmt.Println("Been here:", str)
			if targetStr, ok := target.(string); ok {
				canAddNode = flexMatch(str, targetStr, false)
			}
		}

		if canAddNode {
			*nodes = append(*nodes, Document{Root: n})
			if limit == 1 {
				return 0
			}
		}

	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		status := d.traverse(c, nodes)
		if status == 0 {
			return 0
		}
	}

	return 1
}

func (d Document) FindAll(name string, params ...map[string]any) []Document {
	var docs = []Document{}
	if params == nil {
		var emptyMap map[string]any = make(map[string]any)
		d.params = emptyMap
	} else {
		d.params = params[0]
	}
	d.isStrict = false
	d.params["_name_"] = name

	d.traverse(d.Root, &docs)
	return docs
}

func (d Document) FindOne(name string, params ...map[string]any) Document {
	d.limit = 1
	var docs []Document = d.FindAll(name, params...)
	d.traverse(d.Root, &docs)

	if len(docs) == 0 {
		return Document{}
	}
	return docs[0]
}

var whitespaceRegex = regexp.MustCompile(`\s`)

func hasWhitespaceRegex(s string) bool {
	return whitespaceRegex.MatchString(s)
}

func (d Document) Print(root *html.Node, depth *uint16) string {
	n := d.Root
	if root != nil {
		n = root
	}

	depthCount := *depth
	var builder strings.Builder

	if n.Type == html.TextNode {
		input := strings.TrimSpace(n.Data)

		if input == "" {
			return ""
		}
		builder.WriteString(input)
		return builder.String()
	}

	if n.Type == html.ElementNode {
		*depth++
		builder.WriteString("\n")
		builder.WriteString("<")
		builder.WriteString(n.Data)
		builder.WriteString(">")
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		input := d.Print(c, depth)
		if input != "" {
			// builder.WriteString("\n")
			// fmt.Println("been ")
			builder.WriteString(input)
			// builder.WriteString("\n")
		}
	}

	if n.Type == html.ElementNode {
		fmt.Println("Depth", *depth, "Data", n.Data, "count", depthCount)
		*depth--
		// builder.WriteString("\n")
		builder.WriteString("</")
		builder.WriteString(n.Data)
		builder.WriteString(">")

	}
	return builder.String()
}

func FormatPseudoHTML(input string, indentWidth int) string {
	// 1. Clean the input
	input = strings.TrimSpace(input)

	var builder strings.Builder
	indentLevel := 0

	// Convert to runes for Unicode safety (Amharic, Mandarin, etc.)
	runes := []rune(input)
	n := len(runes)

	for i := range n {
		// Detect closing tag </
		if i+1 < n && runes[i] == '<' && runes[i+1] == '/' {
			indentLevel--

			// SAFETY: Prevent negative repeat count panic
			if indentLevel < 0 {
				indentLevel = 0
			}

			builder.WriteRune('\n')
			builder.WriteString(strings.Repeat(" ", indentLevel*indentWidth))
		}

		builder.WriteRune(runes[i])

		// Detect end of an opening tag >
		if runes[i] == '>' {
			// Look ahead to see if we should indent the next line
			if i+1 < n {
				// Only indent if the next tag is NOT a closing tag
				// and if we aren't looking at a self-closing tag like <br/>
				isClosingNext := (i+2 < n && runes[i+1] == '<' && runes[i+2] == '/')
				isSelfClosing := (i > 0 && runes[i-1] == '/')

				if !isClosingNext && !isSelfClosing {
					indentLevel++
					builder.WriteRune('\n')
					builder.WriteString(strings.Repeat(" ", indentLevel*indentWidth))
				}
			}
		}
	}

	return builder.String()
}
