package main

import (
	"fmt"
	"strings"

	"github.com/halas77/goscrape/engine"
)

func main() {
	requestTest()
	stringTest()
	exportDemo()
}

func requestTest() {

	input := `
	<article class="data-card" data-category="electronics">
		<div class="info">
			<span class="price-tag">$899.99</span>
			<span class="stock">In Stock</span>
		</div>

		<div class="details">
			<div class="price-tag"> ETB 909.99 </div>
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
	params := []*engine.Attribute{
		{
			Key:   "class",
			Value: "details",
		},
	}

	fmt.Println(scrape.FindFirst("div", params).Print())

	// main := scrape.FindAll("span", map[string]any{"class": "item-count"})
	// scrape.Find("span", map[string]any{"class": "item-count"}, func(t engine.Tag) {
	// 	fmt.Println(t.Print(), ", ")
	// })

	// main := scrape.FindAll("div", params)
	// fmt.Println(main.Print())

	// scrape.Find("div", params, func(t engine.Tag) {
	// 	fmt.Println(t.Print())
	// })

	// scrape.Query(".price-tag", func(t *engine.Tag) {
	// 	fmt.Println(t.Print())
	// })

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
	prices := root.SelectAll(".price-tag")
	for i := range *prices {
		p := (*prices)[i]
		fmt.Println("Price found:", p.Print())
	}

	fmt.Println("\n=== 2. CSS + Attribute Filtering (category: electronics) ===")
	electronics := root.SelectAll("article.data-card[data-category='electronics']")
	for i := range *electronics {
		e := (*electronics)[i]
		// Use a specific CSS selector instead of the extra params map
		brand := e.Select("span.brand")
		fmt.Printf("Electronic Item Brand: %s\n", brand.Print())
	}

	fmt.Println("\n=== 3. CSS + Deep Text Filtering ('Out of Stock') ===")
	// Note: CSS selectors don't natively support deep text matching in all variants,
	// but for now we simplify the example to use a selector.
	outOfStockItems := root.SelectAll("article")
	for _, item := range *outOfStockItems {
		if strings.Contains(item.Text(), "Out of Stock") {
			model := item.Select(".model")
			fmt.Printf("Item Out of Stock: %s\n", model.Print())
		}
	}

	fmt.Println("\n=== 4. Chained CSS Search ===")
	brands := root.SelectAll("article[data-category='electronics']").SelectAll(".brand")
	for _, b := range *brands {
		fmt.Println("Brand in Electronics:", b.Print())
	}
}

func exportDemo() {
	input := `
			<div id="store-front">
				<div class="product-item" data-sku="SKU-NEURAL">
					<h3 class="name">Neural Link V1</h3>
					<span class="price">$1200.00</span>
					<a href="/products/neural" class="link">View Details</a>
				</div>
				<div class="product-item" data-sku="SKU-GLOVE">
					<h3 class="name">Haptic Glove Pro</h3>
					<span class="price">$650.00</span>
					<a href="/products/haptic" class="link">View Details</a>
				</div>
			</div>
		`

	root, err := engine.InitDocument(input)
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return
	}

	products := root.SelectAll(".product-item")

	fmt.Println("\n=== 5. Built-in Structured Export mapping (Tag mapping to structured JSON/CSV) ===")
	// Map extracted data with both text nodes and attributes (with @ prefix)
	mapping := map[string]string{
		"sku":   "@data-sku",
		"name":  "h3.name",
		"price": "span.price",
		"url":   "a.link @href",
	}

	mappedData := products.Map(mapping)

	// 1. Export structured data to JSON
	jsonBytes, err := engine.ExportJSON(mappedData)
	if err != nil {
		fmt.Println("Error exporting JSON:", err)
		return
	}
	fmt.Println("Structured Mapped JSON:\n", string(jsonBytes))

	// 2. Export structured data to CSV
	csvStr, err := engine.ExportCSV(mappedData)
	if err != nil {
		fmt.Println("Error exporting CSV:", err)
		return
	}
	fmt.Println("Structured Mapped CSV:\n", csvStr)

	// 3. Export structured data to Markdown
	mdStr, err := engine.ExportMD(mappedData)
	if err != nil {
		fmt.Println("Error exporting MD:", err)
		return
	}
	fmt.Println("Structured Mapped Markdown Table:\n", mdStr)

	// 4. Write them directly to files
	err = engine.WriteMappedJSON("scraped_products.json", mappedData)
	if err != nil {
		fmt.Println("Error writing JSON file:", err)
	} else {
		fmt.Println("Successfully wrote structured data to: scraped_products.json")
	}

	err = engine.WriteMappedCSV("scraped_products.csv", mappedData)
	if err != nil {
		fmt.Println("Error writing CSV file:", err)
	} else {
		fmt.Println("Successfully wrote structured data to: scraped_products.csv")
	}

	err = engine.WriteMappedMD("scraped_products.md", mappedData)
	if err != nil {
		fmt.Println("Error writing MD file:", err)
	} else {
		fmt.Println("Successfully wrote structured data to: scraped_products.md")
	}

	fmt.Println("\n=== 6. Direct Tags Collection Export (JSON/CSV/MD) ===")
	// Direct Tags to JSON
	directJSON, err := products.ToJSON()
	if err == nil {
		fmt.Println("Direct Tags JSON preview (first 250 chars):\n", string(directJSON)[:250]+"...")
	}

	// Direct Tags to CSV
	directCSV, err := products.ToCSV()
	if err == nil {
		fmt.Println("Direct Tags CSV:\n", directCSV)
	}

	// Direct Tags to Markdown
	directMD, err := products.ToMD()
	if err == nil {
		fmt.Println("Direct Tags Markdown:\n", directMD)
	}
}

