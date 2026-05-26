package main

import (
	"fmt"

	"github.com/halas77/goscrape/engine"
)

func main() {
	requestTest()
	// stringTest()
}

func requestTest() {

	input := `
	<article class="data-card" data-category="electronics">
		<div class="info">
			<span class="price-tag">$899.99</span>
			<span class="stock">In Stock</span>
		</div>

		<div class="details">
			<span class="brand">Apple</span>
			<span class="model">MacBook Air</span>
		</div>
	</article>
	`
	// scrape, err := engine.LoadDocument("http://127.0.0.1:5501/examples/basic/index.html")
	scrape, err := engine.InitDocument(input)

	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return
	}

	// main := scrape.FindAll("span", map[string]any{"class": "item-count"})
	// scrape.Find("span", map[string]any{"class": "item-count"}, func(t engine.Tag) {
	// 	fmt.Println(t.Print(), ", ")
	// })
	params := []*engine.Attribute{
		{
			Key:   "class",
			Value: "details",
		},
	}

	main := scrape.FindAll("div", params)
	fmt.Println(main.Print())

}

func stringTest() {
	input := `
			<div id="inventory-container">
				<article class="data-card" data-category="electronics">
					<div class="info">
						<span class="price-tag">$899.99</span>
						<span class="stock">In Stock</span>
					</div>

					<div class="details">
						<span class="brand">Apple</span>
						<span class="model">MacBook Air</span>
					</div>
				</article>

				<article class="data-card" data-category="appliances">
					<div class="info">
						<span class="price-tag">$450.00</span>
						<span class="stock">Out of Stock</span>
					</div>
					<div class="details">
						<span class="brand">Samsung</span>
						<span class="model">Fridge</span>
					</div>
				</article>
			</div>
		`

	root, err := engine.InitDocument(input)

	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return
	}

	fmt.Println("=== 1. CSS Selector Search (.price-tag) ===")
	prices := root.Select(".price-tag")
	for _, p := range prices {
		fmt.Println("Price found:", p.Print())
	}

	fmt.Println("\n=== 2. CSS + Attribute Filtering (category: electronics) ===")
	electronics := root.Select("article.data-card", map[string]any{"data-category": "electronics"})
	for _, e := range electronics {
		brand := e.SelectOne("span", map[string]any{"class": "brand"})
		fmt.Printf("Electronic Item Brand: %s\n", brand.Print())
	}

	fmt.Println("\n=== 3. CSS + Deep Text Filtering ('Out of Stock') ===")
	outOfStockItems := root.Select("article", map[string]any{"string": "Out of Stock"})
	for _, item := range outOfStockItems {
		model := item.SelectOne(".model")
		fmt.Printf("Item Out of Stock: %s\n", model.Print())
	}

	fmt.Println("\n=== 4. Chained CSS Search ===")
	brands := root.Select("article[data-category='electronics']").Select(".brand")
	for _, b := range brands {
		fmt.Println("Brand in Electronics:", b.Print())
	}
}
