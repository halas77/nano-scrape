package engine

/*

import (
	"fmt"
	"strings"

	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

type Document struct {
	Root        *html.Node
	isStrict    bool
	params      map[string]any
	limit       uint8
	isInitial   bool
	currentNode *html.Node
}

type Selection struct {
	Nodes []*html.Node
	Doc   *Document
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

func (d *Document) Find(selector string) *Selection {
	sel, err := cascadia.Parse(selector)
	if err != nil {
		return &Selection{Doc: d}
	}
	matches := cascadia.QueryAll(d.Root, sel)
	return &Selection{Nodes: matches, Doc: d}
}

func (s *Selection) Find(selector string) *Selection {
	sel, err := cascadia.Parse(selector)
	if err != nil || len(s.Nodes) == 0 {
		return &Selection{Doc: s.Doc}
	}

	var newNodes []*html.Node
	for _, n := range s.Nodes {
		matches := cascadia.QueryAll(n, sel)
		newNodes = append(newNodes, matches...)
	}
	return &Selection{Nodes: newNodes, Doc: s.Doc}
}

func matchAttributes(n *html.Node, target map[string]string) bool {
	if len(target) == 0 {
		return true
	}

	actualAttrs := make(map[string]string)
	for _, attr := range n.Attr {
		actualAttrs[attr.Key] = attr.Val
	}

	for key, expectedVal := range target {
		actualVal, exists := actualAttrs[key]
		if !exists || actualVal != expectedVal {
			return false
		}
	}

	return true
}

func searchNodes(n *html.Node, name string, attrs map[string]string, matches *[]*html.Node) {
	if n.Type == html.ElementNode && n.Data == name {
		if matchAttributes(n, attrs) {
			*matches = append(*matches, n)
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		searchNodes(c, name, attrs, matches)
	}
}

func (s *Selection) First() *Selection {
	if len(s.Nodes) == 0 {
		return s
	}
	return &Selection{Nodes: s.Nodes[:1], Doc: s.Doc}
}

func (s *Selection) Eq(i int) *Selection {
	if i < 0 || i >= len(s.Nodes) {
		return &Selection{Doc: s.Doc}
	}
	return &Selection{Nodes: []*html.Node{s.Nodes[i]}, Doc: s.Doc}
}

func (s *Selection) Length() int {
	return len(s.Nodes)
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

func (d Document) Print(indentWidth ...int) string {

	n := d.Root
	if d.currentNode != nil {
		n = d.currentNode
	}

	var builder strings.Builder
	factor := 0
	width := 0
	if indentWidth != nil {
		width = int(indentWidth[0])
	}

	if n.Type != html.DocumentNode {
		factor = 3
	}

	if n.Type == html.TextNode {
		builder.WriteString(strings.Repeat(" ", int(width)))
		input := strings.TrimSpace(n.Data)

		if input == "" {
			return ""
		}
		builder.WriteString(input)
		return builder.String()
	}

	if n.Type == html.ElementNode {
		builder.WriteString(strings.Repeat(" ", int(width)))
		builder.WriteString("<")
		builder.WriteString(n.Data)
		builder.WriteString(">")
	}

	canIndent := false

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		canIndent = true

		d.currentNode = c
		input := d.Print(width + factor)
		if c == n.FirstChild {
			builder.WriteString("\n")
		}
		if input != "" {
			builder.WriteString(input)
			builder.WriteString("\n")
		}
	}

	if n.Type == html.ElementNode {
		if canIndent {
			builder.WriteString(strings.Repeat(" ", int(width)))
		}

		builder.WriteString("</")
		builder.WriteString(n.Data)
		builder.WriteString(">")

	}
	return builder.String()
}
*/
