package engine

import (
	"strings"

	"golang.org/x/net/html"
)

type Document struct {
	Root *html.Node
}

func IsType[T any](val any) bool {
	_, ok := val.(T)
	return ok
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

func hasIntersection(params map[string]any, attributes []html.Attribute, isStrict bool) bool {
	if params == nil {
		return true
	}

	var count uint8 = 0
	length := len(params)
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

func traverse(n *html.Node, name string, nodes *[]Document, limit int8, params map[string]any) uint8 {

	if n.Type == html.ElementNode && n.Data == name {
		if hasIntersection(params, n.Attr, false) {
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
