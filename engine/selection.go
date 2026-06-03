package engine

import (
	"fmt"

	"golang.org/x/net/html"
)

type SelectionParams struct {
	name    string
	params  map[string]any
	attrs   []*Attribute
	limit   uint8
	recurse bool
}

func (t *Tag) traverse(n *html.Node, f func(*html.Node) bool) uint8 {

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		// t.nameSelector(c)
		exit := f(c)
		if t.limit == 1 && exit {
			return 0
		}

		if t.recurse {
			status := t.traverse(c, f)
			if status == 0 {
				return 0
			}
		}
	}

	return 1
}

func (t *Tag) FindMatchingAttributes(elementAttrs []html.Attribute, n *html.Node) bool {
	attrs := t.attrs

	if attrs == nil {
		return true
	}

	if t.maps == nil {
		lookup := make(map[string]string)

		for _, attr := range attrs {
			// normalizedKey := strings.ToLower(attr.Key)
			normalizedKey := attr.Key
			lookup[normalizedKey] = attr.Value
		}

		t.maps = &lookup
	}

	var attrsLength uint8 = uint8(len(attrs))
	var counter uint8 = 0

	for _, attr := range elementAttrs {
		// normalizedKey := strings.ToLower(attr.Key)
		normalizedKey := attr.Key
		if valA, found := (*t.maps)[normalizedKey]; found {
			if valA == attr.Val {
				counter++
			}
		}
	}

	if stringVal, found := (*t.maps)["string"]; found {
		str := t.getNodeStrings(n)
		fmt.Println("str ", str)
		if flexMatch(str, stringVal, false) {
			counter++
		}
	}

	return attrsLength == counter
}

func (t *Tag) nameSelector(n *html.Node) bool {
	if n.Type == html.ElementNode && n.Data == t.name {
		return t.FindMatchingAttributes(n.Attr, n)
	}
	return false
}

// Add Traverse for selection Params

func hasIntersection(params map[string]any, attributes []html.Attribute, isStrict bool) bool {
	if params == nil {
		return true
	}

	// Calculate how many attributes we actually need to match (ignoring internal keys)
	expectedAttrCount := 0
	for key := range params {
		if key != "string" && key != "_name_" {
			expectedAttrCount++
		}
	}

	// If no attributes are provided to match, it's an automatic pass for attributes
	if expectedAttrCount == 0 {
		return true
	}

	matchCount := 0
	for _, attr := range attributes {
		if val, ok := params[attr.Key]; ok && val == attr.Val {
			matchCount++
		}
	}

	if isStrict {
		return matchCount == expectedAttrCount
	}
	return matchCount > 0

}

func nameSelector(t Tag, params map[string]any) bool {
	name, _ := params["_name_"]
	canAddNode := false
	n := t.root

	if n.Type == html.ElementNode && n.Data == name {

		if params == nil {
			canAddNode = true
		} else {
			canAddNode = hasIntersection(params, n.Attr, false)
		}
		target, ok := params["string"]

		if ok {
			str := getNodeStrings(t)
			if targetStr, ok := target.(string); ok {
				canAddNode = flexMatch(str, targetStr, false)
			}
		}
	}

	return canAddNode
}

// use a callback to make the filtering logic only use traverse function
func traverse(tag Tag, limit uint8, recurse bool, cb TraverseCallback) uint8 {
	n := tag.root

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		t := initTag(c)
		exit := cb(*t)
		if limit == 1 && exit {
			return 0
		}

		if recurse {
			status := traverse(*t, limit, recurse, cb)
			if status == 0 {
				return 0
			}
		}
	}

	return 1
}
