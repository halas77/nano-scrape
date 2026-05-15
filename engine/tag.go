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

func (t Tag) Find(name string, params ...map[string]any) Tags {
	var tags Tags
	p := getParams(params)
	p["_name_"] = name

	selectionParams := SelectionParams{params: p}
	traverse(t, t.limit, true, func(match Tag) bool {
		if nameSelector(match, selectionParams.params) {
			tags = append(tags, match)
		}
		return false
	})
	return tags
}

func (t Tag) FindOne(name string, params ...map[string]any) Tag {
	t.limit = 1
	return t.Find(name, params...).First()
}

func (t Tag) Matches(params map[string]any) bool {
	if params == nil {
		return true
	}
	if !hasIntersection(params, t.Attrs, false) {
		return false
	}
	if target, ok := params["string"]; ok {
		if targetStr, ok := target.(string); ok {
			return flexMatch(t.Text(), targetStr, false)
		}
	}
	return true
}

func (t Tag) query(selector string, params map[string]any, cb func(Tag) bool) {
	sel, err := cascadia.Parse(selector)
	if err != nil {
		return
	}

	for _, n := range cascadia.QueryAll(t.root, sel) {
		candidate := initTag(n)
		if candidate.Matches(params) {
			if cb(candidate) {
				return
			}
		}
	}
}

func (t Tag) Select(selector string, params ...map[string]any) (results Tags) {
	t.query(selector, getParams(params), func(match Tag) bool {
		results = append(results, match)
		return false
	})
	return
}

func (t Tag) SelectOne(selector string, params ...map[string]any) (result Tag) {
	t.query(selector, getParams(params), func(match Tag) bool {
		result = match
		return true
	})
	return
}

func (t Tag) Text() string {
	return strings.TrimSpace(getNodeStrings(t))
}

func (ts Tags) Select(selector string, params ...map[string]any) (results Tags) {
	p := getParams(params)
	for _, t := range ts {
		t.query(selector, p, func(match Tag) bool {
			results = append(results, match)
			return false
		})
	}
	return
}

func (ts Tags) SelectOne(selector string, params ...map[string]any) Tag {
	return ts.Select(selector, params...).First()
}

func (ts Tags) Each(cb func(Tag)) {
	for _, t := range ts {
		cb(t)
	}
}

func (ts Tags) First() Tag {
	if len(ts) == 0 {
		return Tag{}
	}
	return ts[0]
}

func getParams(params []map[string]any) map[string]any {
	if len(params) > 0 {
		return params[0]
	}
	return nil
}
