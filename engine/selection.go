package engine

import (
	"strings"

	"golang.org/x/net/html"
)

type SelectionParams struct {
	name      string
	className string
	id        string
	params    map[string]any
	limit     uint8
}

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

func getNodeStrings(n *html.Node) string {
	nodeStrings := []string{}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			nodeStrings = append(nodeStrings, c.Data)
		}
	}

	return strings.Join(nodeStrings, "")
}

// use a callback to make the filtering logic only use traverse function
func (s SelectionParams) traverse(n *html.Node, nodes *[]Tag) uint8 {
	name, _ := s.params["_name_"]
	limit := s.params["limit"]

	if n.Type == html.ElementNode && n.Data == name {
		canAddNode := false
		canAddNode = hasIntersection(s.params, n.Attr, false)
		target, ok := s.params["string"]

		// fmt.Println(target, ok)
		if ok {
			str := getNodeStrings(n)
			if targetStr, ok := target.(string); ok {
				canAddNode = flexMatch(str, targetStr, false)
			}
		}

		if canAddNode {
			*nodes = append(*nodes, Tag{root: n, Name: n.Data, Attrs: n.Attr, Class: "", Id: ""})
			if limit == 1 {
				return 0
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		status := s.traverse(c, nodes)
		if status == 0 {
			return 0
		}
	}

	return 1
}
