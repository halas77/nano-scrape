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

					<span class="info">
						<span class="price-tag">$530</span>
						<hr/>
					</span>
					<strong> look </strong>
				</article>
			</div>
		`

	scrape, _ := engine.Scrape(input)
	// article := scrape.FindAll("article")

	// span := scrape.FindAll("div")

	div := scrape.FindOne("span")
	fmt.Println("div", div.Root.Data)

	// for i := range span {
	// 	node := span[i]
	// 	fmt.Println("span", node.Root.Data)
	// }

	// hr := article.FindAll("hr")
}
