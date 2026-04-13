package main

import (
	"fmt"

	"github.com/halas77/goscrape/engine"
)

func main() {
	input := `
			<div id="inventory-container">
				<article class="data-card">
					<div class="info">
						<span class="price-tag">$899.99</span>
						<hr/>
					</div>

					<span class="info" id="id-1" >
						<span class="price-tag">$530</span>
						<hr/>
					</span>
					<strong> look <b> Bolded </b> 
					Abebe
					</strong>
				</article>
				<b> Bolded </b>
			</div>
		`

	scrape, _ := engine.Scrape(input)
	// article := scrape.FindAll("article")

	// span := scrape.FindAll("div")

	// div := scrape.FindOne("span", map[string]any{"class": "info", "id": "id-1"})
	// fmt.Println("div", div.Root.Data)

	strong := scrape.FindOne("strong", map[string]any{"string": "abe"})
	fmt.Println("strong: ", strong)
	// span := div.FindAll("hr", nil)

	// for i := range span {
	// 	node := span[i]
	// 	fmt.Println("span", node.Root.Data)
	// }

	// hr := article.FindAll("hr")

	// text := "Hello world mother"

	// fmt.Println(RobustMatch(text, "Hello world mother")) // true (Exact)
	// fmt.Println(RobustMatch(text, "orld"))               // true (Substring)
	// fmt.Println(RobustMatch(text, "mother"))             // true (Substring)
	// fmt.Println(RobustMatch(text, "abc"))                // false
}
