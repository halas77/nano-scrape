package engine

import (
	"strings"

	"golang.org/x/net/html"
)

type Tag struct {
	root  *html.Node
	limit uint8

	Name  string
	Attrs []html.Attribute
	Class string
	Id    string
}

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
	traverse(t.root, t.limit, true, func(n *html.Node) bool {
		isMatch := nameSelector(n, selectionParams.params)
		if isMatch {
			tags = append(tags, Tag{root: n, Name: n.Data, Attrs: n.Attr, Class: "", Id: ""})
		}
		return isMatch
	})
	return tags
}

func (t Tag) FindFirst(name string, params ...map[string]any) Tag {
	var p map[string]any = make(map[string]any)

	if len(params) > 0 {
		p = params[0]
	}

	t.limit = 1
	var tags []Tag = t.FindAll(name, p)

	if len(tags) == 0 {
		return Tag{}
	}
	return tags[0]
}
