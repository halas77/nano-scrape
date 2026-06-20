package engine

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

type TraverseCallback func(Tag) bool

type Attribute struct {
	Key   string
	Value string
}

func initTag(node *html.Node) *Tag {
	return &Tag{root: node, Name: node.Data, Attrs: node.Attr}
}

func InitDocument(input any) (*Tag, error) {
	var reader io.Reader

	switch v := input.(type) {
	case string:
		reader = strings.NewReader(v)
	case []byte:
		// Efficiently converts []byte to a string reader without extra allocations
		reader = strings.NewReader(string(v))
	case io.Reader:
		// If it's already a reader (like a network response body), use it directly
		reader = v
	default:
		return nil, fmt.Errorf("unsupported input type: %T", input)
	}

	// If the reader is also a closer (like an HTTP response body),
	// ensure it gets closed when this function finishes.
	if closer, ok := reader.(io.Closer); ok {
		defer closer.Close()
	}

	node, err := html.Parse(reader)
	if err != nil {
		return nil, err
	}

	document := initTag(node)
	return document, nil
}

func LoadDocument(url string) (*Tag, error) {
	resp, err := NewClient().Execute("GET", url, nil)

	if err != nil {
		return nil, err
	}

	return InitDocument(resp)

}

func (t *Tag) Next() *Tag {
	return initTag(t.root.NextSibling)
}

func (t *Tag) Previous() *Tag {
	return initTag(t.root.PrevSibling)
}

func (t *Tag) Parent() *Tag {
	return initTag(t.root.Parent)
}

func (t *Tag) FirstChild() *Tag {
	return initTag(t.root.FirstChild)
}

func (t *Tag) LastChild() *Tag {
	return initTag(t.root.LastChild)
}

func InitTag(node *html.Node) *Tag {
	return initTag(node)
}
