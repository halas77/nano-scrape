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

func initTag(node *html.Node) Tag {
	return Tag{root: node, Name: node.Data, Attrs: node.Attr}
}

func InitDocument(input any) (Tag, error) {
	var reader io.Reader

	switch v := any(input).(type) {
	case string:
		reader = strings.NewReader(v)
	case []byte:
		reader = strings.NewReader(string(v))
	default:
		return Tag{}, fmt.Errorf("unsupported input type: %T", input)
	}

	node, err := html.Parse(reader)
	if err != nil {
		return Tag{}, err
	}

	document := initTag(node)
	return document, nil
}

func LoadDocument(url string) (Tag, error) {
	resp, err := InitRequest(url, "GET", nil).Execute()

	if err != nil {
		return Tag{}, err
	}

	return InitDocument(resp)

}

func (t Tag) Next() Tag {
	return initTag(t.root.NextSibling)
}

func (t Tag) Previous() Tag {
	return initTag(t.root.PrevSibling)
}

func (t Tag) Parent() Tag {
	return initTag(t.root.Parent)
}

func (t Tag) FirstChild() Tag {
	return initTag(t.root.FirstChild)
}

func (t Tag) LastChild() Tag {
	return initTag(t.root.LastChild)
}
