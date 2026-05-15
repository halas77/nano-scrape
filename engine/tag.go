package engine

import (
	"strings"

	"github.com/andybalholm/cascadia"
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

type Tags []Tag
type TagCallback func(Tag)

func (t Tag) Find(name string, params map[string]any, cb TagCallback) {
	var p map[string]any = make(map[string]any)

	if params != nil {
		p = params
	}
	p["_name_"] = name

	selectionParams := SelectionParams{params: p}
	traverse(t, t.limit, true, func(t Tag) bool {
		isMatch := nameSelector(t, selectionParams.params)
		if isMatch {
			cb(t)
		}
		return isMatch
	})
}

func (t Tag) FindAll(name string, params ...map[string]any) Tags {
	var tags = Tags{}
	t.Find(name, params[0], func(t Tag) {
		tags = append(tags, t)
	})

	return tags
}

func (t Tag) FindFirst(name string, params ...map[string]any) Tag {

	t.limit = 1
	var tags Tags = t.FindAll(name, params...)

	if len(tags) == 0 {
		return Tag{}
	}
	return tags[0]
}

func (t Tag) Select(selector string, params ...map[string]any) Tags {
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
				str := getNodeStrings(t)
				if targetStr, ok := target.(string); ok {
					if !flexMatch(str, targetStr, false) {
						continue
					}
				}
			}
		}
		tags = append(tags, Tag{root: n, Name: n.Data, Attrs: n.Attr})
	}
	return tags
}

func (t Tag) SelectOne(selector string, params ...map[string]any) Tag {
	return t.Select(selector, params...).First()
}

func (t Tag) Text() string {
	return strings.TrimSpace(getNodeStrings(t))
}

func (ts Tags) Select(selector string, params ...map[string]any) Tags {
	var results Tags
	for _, t := range ts {
		results = append(results, t.Select(selector, params...)...)
	}
	return results
}

func (ts Tags) First() Tag {
	if len(ts) == 0 {
		return Tag{}
	}
	return ts[0]
}

func (ts Tags) SelectOne(selector string, params ...map[string]any) Tag {
	return ts.Select(selector, params...).First()
}
