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

func (ts Tags) First() Tag {

	if len(ts) == 0 {
		return Tag{}
	}
	return *ts[0]
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

// Find query

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

func (t *Tag) FindAll(name string, attr []*Attribute) Tags {

	var tags = Tags{}
	t.Find(name, attr, func(t *Tag) {
		tags = append(tags, t)
	})

	return tags
}

func (t Tag) FindFirst(name string, params ...map[string]any) Tag {
	t.limit = 1
	// return t.FindAll(name, params...).First()
	return Tag{}
}

// Select Functionality for css selectors
func (t *Tag) Select(selector string, params ...map[string]any) Tags {
	sel, err := cascadia.Parse(selector)
	if err != nil {
		return Tags{}
	}

	matches := cascadia.QueryAll(t.root, sel)
	var tags = Tags{}

	var p map[string]any
	if len(params) > 0 {
		p = params[0]
	}

	for _, n := range matches {
		if p != nil {
			if !hasIntersection(p, n.Attr, false) {
				continue
			}
			if target, ok := p["string"]; ok {
				str := getNodeStrings(*t)
				if targetStr, ok := target.(string); ok {
					if !flexMatch(str, targetStr, false) {
						continue
					}
				}
			}
		}
		tags = append(tags, &Tag{root: n, Name: n.Data, Attrs: n.Attr})
	}
	return tags
}

func (t Tag) SelectOne(selector string, params ...map[string]any) Tag {
	return t.Select(selector, params...).First()
}

func (t *Tag) Text() string {
	return strings.TrimSpace(t.getNodeStrings(t.root))
}

func (ts Tags) Select(selector string, params ...map[string]any) *Tags {
	var results Tags
	for _, t := range ts {
		results = append(results, t.Select(selector, params...)...)
	}
	return &results
}

func (ts Tags) SelectOne(selector string, params ...map[string]any) Tag {
	return ts.Select(selector, params...).First()
}
