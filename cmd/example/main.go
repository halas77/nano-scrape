package main

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	rawHTML := "<html><head><title>Example</title></head><body><h1>Hello Go</h1></body></html>"
	doc, _ := html.Parse(strings.NewReader(rawHTML))

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			fmt.Println("Found tag:", n.Data)
		}
		for c := n.FirstChild; c != nil; c = n.NextSibling {
			f(c)
		}
	}
	f(doc)
}
