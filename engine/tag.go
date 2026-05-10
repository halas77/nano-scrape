package engine

import (
	"strings"

	"golang.org/x/net/html"
)

type Tag struct {
	root  *html.Node
	Name  string
	Attrs []html.Attribute
	Class string
	Id    string
}

type TraverseCallback func(*html.Node) bool

func InitDocument(input string) (Tag, error) {
	reader := strings.NewReader(input)
	node, err := html.Parse(reader)
	document := Tag{root: node, Name: node.Data, Attrs: node.Attr}
	if err != nil {
		return Tag{}, err
	}

	return document, nil
}

func (t Tag) FindAll(name string, params ...map[string]any) []Tag {
	var tags = []Tag{}
	var p map[string]any = make(map[string]any)

	if len(params) > 0 {
		p = params[0]
	}
	p["_name_"] = name

	selectionParams := SelectionParams{params: p}
	selectionParams.traverse(t.root, &tags)
	return tags
}

func (t Tag) FindFirst(name string, params ...map[string]any) Tag {
	var p map[string]any = make(map[string]any)

	if len(params) > 0 {
		p = params[0]
	}
	p["limit"] = 0

	var tags []Tag = t.FindAll(name, p)

	if len(tags) == 0 {
		return Tag{}
	}
	return tags[0]
}
