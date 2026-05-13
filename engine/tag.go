package engine

import (
	"fmt"
	"io"
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

	document := Tag{root: node, Name: node.Data, Attrs: node.Attr}
	return document, nil
}

func LoadDocument(url string) (Tag, error) {
	resp, err := InitRequest(url, "GET", nil).Execute()

	if err != nil {
		return Tag{}, err
	}

	return InitDocument(resp)

}

func (t Tag) FindAll(name string, params ...map[string]any) Tags {
	var tags = Tags{}
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
	var tags Tags = t.FindAll(name, p)

	if len(tags) == 0 {
		return Tag{}
	}
	return tags[0]
}
func (t Tag) Find(selector string, params ...map[string]any) Tags {
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
				str := getNodeStrings(n)
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

func (t Tag) FindOne(selector string, params ...map[string]any) Tag {
	return t.Find(selector, params...).First()
}

func (t Tag) Text() string {
	return strings.TrimSpace(getNodeStrings(t.root))
}

func (ts Tags) Find(selector string, params ...map[string]any) Tags {
	var results Tags
	for _, t := range ts {
		results = append(results, t.Find(selector, params...)...)
	}
	return results
}

func (ts Tags) First() Tag {
	if len(ts) == 0 {
		return Tag{}
	}
	return ts[0]
}

func (ts Tags) FindOne(selector string, params ...map[string]any) Tag {
	return ts.Find(selector, params...).First()
}
