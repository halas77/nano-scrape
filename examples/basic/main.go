package main

import (
	"fmt"
	"regexp"
)

func main() {
	// input := `
	// 		<div id="inventory-container">
	// 			<article class="data-card">
	// 				<div class="info">
	// 					<span class="price-tag">$899.99</span>
	// 					<hr/>
	// 				</div>

	// 				<span class="info" id="id-1" >
	// 					<span class="price-tag">$530</span>
	// 					<hr/>
	// 				</span>
	// 				<strong> look </strong>
	// 			</article>
	// 			<b> Bolded </b>
	// 		</div>
	// 	`

	// scrape, _ := engine.Scrape(input)
	// article := scrape.FindAll("article")

	// span := scrape.FindAll("div")

	// div := scrape.FindOne("span", map[string]any{"class": "info", "id": "id-1"})
	// fmt.Println("div", div.Root.Data)

	// span := div.FindAll("hr", nil)

	// for i := range span {
	// 	node := span[i]
	// 	fmt.Println("span", node.Root.Data)
	// }

	// hr := article.FindAll("hr")

	text := "Hello world mother"

	fmt.Println(RobustMatch(text, "Hello world mother")) // true (Exact)
	fmt.Println(RobustMatch(text, "orld"))               // true (Substring)
	fmt.Println(RobustMatch(text, "mother"))             // true (Substring)
	fmt.Println(RobustMatch(text, "abc"))                // false
}

// RobustMatch returns true if 'target' is found anywhere within 'main'.
// This includes exact matches and substrings.
func RobustMatch(main string, target string) bool {
	// 1. Escape special characters to treat the target as literal text
	pattern := regexp.QuoteMeta(target)

	// 2. Compile the regex
	re, err := regexp.Compile(pattern)
	if err != nil {
		// If the target is empty or invalid (unlikely with QuoteMeta), return false
		return false
	}

	// 3. MatchString returns true if the pattern matches any part of the string
	return re.MatchString(main)
}
