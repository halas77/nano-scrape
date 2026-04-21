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

	// raw := `
	// <div>
	// 	<p>
	// 		ሰላም (Hello) ✋
	// 	</p>
	// 	<span>
	// 		你好 (Nǐ hǎo) 🚀
	// 	</span>
	// 	<div>
	// 		你好 (Nǐ hǎo) 🚀
	// 	</div>
	// </div>
	//  `
	scrape, _ := engine.Scrape(input)
	article := scrape.FindOne("article")

	fmt.Println(article.Print())

	// divs := scrape.FindAll("div")
	// fmt.Println("divs:", divs)

	// formatted := engine.FormatPseudoHTML(input, 4)
	// fmt.Println(formatted)

	// div := scrape.FindOne("span", map[string]any{"class": "info", "id": "id-1"})
	// fmt.Println("div", div.Root.Data)

	// strong := scrape.FindOne("strong", map[string]any{"string": "abe"})
	// fmt.Println("strong: ", strong.Root.Data)

	// for i := range span {
	// 	node := span[i]
	// 	fmt.Println("span", node.Root.Data)
	// }

	// hr := article.FindAll("hr")

	// text := "Hello world mother"

}
