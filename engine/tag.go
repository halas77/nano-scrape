package engine

import (
	"strings"

	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

type Tag struct {
	root    *html.Node
	limit   uint8
	attrs   []*Attribute
	recurse bool
	name    string
	maps    *map[string]string

	Name  string
	Attrs []html.Attribute
	Class string
	Id    string
}

type Tags []*Tag
type TagCallback func(*Tag)

func (ts *Tags) First() *Tag {

	if ts == nil || len(*ts) == 0 {
		return nil
	}
	return (*ts)[0]
}

func (t *Tag) Query(selector string, f func(*Tag)) {
	sel, err := cascadia.Parse(selector)
	if err != nil {
		return
	}

	t.recurse = true

	t.traverse(t.root, func(n *html.Node) bool {
		hasMatch := sel.Match(n)
		if hasMatch {
			f(&Tag{root: n, Name: n.Data, Attrs: n.Attr})
		}

		return hasMatch
	})
}

func (t *Tag) QueryAll(selector string) *Tags {

	var tags = &Tags{}
	t.Query(selector, func(t *Tag) {
		*tags = append(*tags, t)
	})

	return tags
}

func (t *Tag) QueryOne(selector string) *Tag {
	t.limit = 1
	return t.QueryAll(selector).First()
}

func (t *Tag) Find(name string, attrs []*Attribute, cb TagCallback) {

	t.recurse = true
	t.name = name
	t.attrs = attrs

	t.traverse(t.root, func(node *html.Node) bool {
		isMatch := t.nameSelector(node)
		if isMatch {
			cb(&Tag{root: node, Name: node.Data, Attrs: node.Attr})
		}
		return isMatch
	})
}

func (t *Tag) FindAll(name string, attribute ...[]*Attribute) *Tags {

	var tags = &Tags{}
	var attr []*Attribute
	if len(attribute) > 0 {
		attr = attribute[0]
	}

	t.Find(name, attr, func(t *Tag) {
		*tags = append(*tags, t)
	})

	return tags
}

func (t *Tag) FindFirst(name string, attr ...[]*Attribute) *Tag {
	t.limit = 1
	return t.FindAll(name, attr...).First()
}

func (t *Tag) SelectAll(selector string) *Tags {
	return t.QueryAll(selector)
}

func (t *Tag) Select(selector string) *Tag {
	return t.QueryOne(selector)
}

func (t *Tag) Text() string {
	return strings.TrimSpace(t.getNodeStrings(t.root))
}
