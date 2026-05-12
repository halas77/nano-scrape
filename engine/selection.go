package engine

import (
	"golang.org/x/net/html"
)

type SelectionParams struct {
	name      string
	className string
	id        string
	params    map[string]any
	limit     uint8
}

type TraverseCallback func(*html.Node) bool

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

func nameSelector(n *html.Node, params map[string]any) bool {
	name, _ := params["_name_"]
	canAddNode := false

	if n.Type == html.ElementNode && n.Data == name {
		canAddNode = hasIntersection(params, n.Attr, false)
		target, ok := params["string"]

		if ok {
			str := getNodeStrings(n)
			if targetStr, ok := target.(string); ok {
				canAddNode = flexMatch(str, targetStr, false)
			}
		}
	}

	return canAddNode
}

// use a callback to make the filtering logic only use traverse function
func traverse(n *html.Node, limit uint8, recurse bool, cb TraverseCallback) uint8 {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		exit := cb(c)
		if limit == 1 && exit {
			return 0
		}

		if recurse {
			status := traverse(c, limit, recurse, cb)
			if status == 0 {
				return 0
			}
		}
	}

	return 1
}
