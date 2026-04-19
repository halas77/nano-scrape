package engine

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// adding strings function to get strings of a tag not sub tags concatenate them and append them to an array

type Document struct {
	Root      *html.Node
	isStrict  bool
	params    map[string]any
	limit     uint8
	isInitial bool
}

func Scrape(input string) (Document, error) {
	reader := strings.NewReader(input)
	node, err := html.Parse(reader)
	doc := Document{Root: node, isInitial: true}
	fmt.Println("is initial: ", doc.isInitial)

	if err != nil {
		return doc, err
	}

	return doc, nil
}

type Traverser interface {
	// Traverse()
	FindAll(name string)
	FindOne()
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

func flexMatch(main string, target string, caseSensitive bool) bool {
	pattern := regexp.QuoteMeta(target)

	if !caseSensitive {
		pattern = "(?i)" + pattern
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}

	return re.MatchString(main)
}

func (d Document) traverse(n *html.Node, nodes *[]Document) uint8 {
	name, _ := d.params["_name_"]
	limit := d.limit

	if n.Type == html.ElementNode && n.Data == name {
		canAddNode := false
		canAddNode = hasIntersection(d.params, n.Attr, false)
		target, ok := d.params["string"]

		if ok {
			str := getNodeStrings(n)
			fmt.Println("Been here:", str)
			if targetStr, ok := target.(string); ok {
				canAddNode = flexMatch(str, targetStr, false)
			}
		}

		if canAddNode {
			*nodes = append(*nodes, Document{Root: n})
			if limit == 1 {
				return 0
			}
		}

	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		status := d.traverse(c, nodes)
		if status == 0 {
			return 0
		}
	}

	return 1
}

func (d Document) FindAll(name string, params ...map[string]any) []Document {
	var docs = []Document{}
	if params == nil {
		var emptyMap map[string]any = make(map[string]any)
		d.params = emptyMap
	} else {
		d.params = params[0]
	}
	d.isStrict = false
	d.params["_name_"] = name

	d.traverse(d.Root, &docs)
	return docs
}

func (d Document) FindOne(name string, params ...map[string]any) Document {
	d.limit = 1
	var docs []Document = d.FindAll(name, params...)
	d.traverse(d.Root, &docs)

	if len(docs) == 0 {
		return Document{}
	}
	return docs[0]
}
