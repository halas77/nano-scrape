package engine

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// adding strings function to get strings of a tag not sub tags concatenate them and append them to an array
type Document struct {
	Root *html.Node
}

func Scrape(input string) (Document, error) {
	reader := strings.NewReader(input)
	node, err := html.Parse(reader)
	doc := Document{Root: node}

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
	if ok {
		length--
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
	// 1. Escape special characters to treat the target as literal text
	pattern := regexp.QuoteMeta(target)

	// 2. If caseSensitive is false, prepend the "ignore case" flag (?i)
	if !caseSensitive {
		pattern = "(?i)" + pattern
	}

	// 3. Compile the regex
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}

	// 4. Match against the main string
	return re.MatchString(main)
}

func traverse(n *html.Node, name string, nodes *[]Document, limit int8, params map[string]any) uint8 {

	if n.Type == html.ElementNode && n.Data == name {
		canAddNode := false

		canAddNode = hasIntersection(params, n.Attr, false)

		target, ok := params["string"]

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
		status := traverse(c, name, nodes, limit, params)
		if status == 0 {
			return 0
		}
	}

	return 1
}

func (d Document) FindAll(name string, params map[string]any) []Document {
	var docs = []Document{}
	traverse(d.Root, name, &docs, -1, params)
	return docs
}

func (d Document) FindOne(name string, params map[string]any) Document {
	var docs = []Document{}
	traverse(d.Root, name, &docs, 1, params)
	if len(docs) == 0 {
		return Document{}
	}
	return docs[0]
}
