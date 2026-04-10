package engine

import (
	"fmt"
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

var count uint8 = 0

func traverse(n *html.Node, name string, nodes *[]Document, limit int8) uint8 {
	count++
	fmt.Println("count", count, "Tag", n.Data)

	if n.Type == html.ElementNode && n.Data == name {
		*nodes = append(*nodes, Document{Root: n})

		if limit == 1 {
			return 0
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		status := traverse(c, name, nodes, limit)
		if status == 0 {
			return 0
		}
	}

	return 1
}

func (d Document) FindAll(name string) []Document {
	var docs = []Document{}
	traverse(d.Root, name, &docs, -1)
	return docs
}

func (d Document) FindOne(name string) Document {
	var docs = []Document{}
	traverse(d.Root, name, &docs, 1)
	return docs[0]
}
