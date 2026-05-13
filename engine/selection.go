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

	var count uint8 = 0
	length := len(params)
	_, ok := params["string"]
	_, name := params["_name_"]

	if ok || name {
		length = 0
	}

	if length == int(count) {
		return true
	}

	for _, attr := range attributes {
		value, ok := params[attr.Key]

		if ok && value == attr.Val {
			count++
		}
	}

	if count > 0 && uint8(length) == count {
		return true
	} else if !isStrict && count > 0 {
		return true
	}

	return false
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
