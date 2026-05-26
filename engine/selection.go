package engine

import (
	"strings"

	"golang.org/x/net/html"
)

type SelectionParams struct {
	name    string
	params  map[string]any
	attrs   []*Attribute
	limit   uint8
	recurse bool
}

func (s SelectionParams) traverse(tag Tag, f TraverseCallback) uint8 {
	n := tag.root
	var limit uint8 = s.limit
	recurse := s.recurse

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		t := initTag(c)
		exit := f(t)
		if s.limit == 1 && exit {
			return 0
		}

		if recurse {
			status := traverse(t, limit, recurse, f)
			if status == 0 {
				return 0
			}
		}
	}

	return 1
}

func (s SelectionParams) FindMatchingAttributes(elementAttrs []html.Attribute) bool {
	lookup := make(map[string]string)
	attrs := s.attrs

	for _, attr := range attrs {
		normalizedKey := strings.ToLower(attr.Key)
		lookup[normalizedKey] = attr.Value
	}

	var attrsLength uint8 = uint8(len(attrs))
	var counter uint8 = 0

	for _, attr := range elementAttrs {
		normalizedKey := strings.ToLower(attr.Key)
		if valA, found := lookup[normalizedKey]; found {
			if normalizedKey == "string" {
				// use flex for equality and increment counter if it is true true and continue
			}

			if valA == attr.Val {
				counter++
			}
		}
	}

	return attrsLength == counter
}

func (s SelectionParams) nameSelector(tag Tag) bool {
	n := tag.root
	if n.Type == html.ElementNode && n.Data == s.name {
		return s.FindMatchingAttributes(n.Attr)
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
		exit := cb(t)
		if limit == 1 && exit {
			return 0
		}

		if recurse {
			status := traverse(t, limit, recurse, cb)
			if status == 0 {
				return 0
			}
		}
	}

	return 1
}
