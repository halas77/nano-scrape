package main

import (
	"fmt"

	"github.com/halas77/goscrape/engine"
)

func main() {
	// requestTest()
	stringTest()
	exportDemo()
	// stringTest()

	// proxyTester()

	testFromPost()
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

	input = "http://127.0.0.1:5501/examples/basic/index.html"
	scrape, err := engine.LoadDocument(input)
	// scrape, err := engine.InitDocument(input)

	fmt.Println(scrape.FindFirst("main").Print())

	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return
	}

	// params := []*engine.Attribute{
	// 	{
	// 		Key:   "class",
	// 		Value: "details",
	// 	},
	// }

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

	products := root.Select(".product-item")

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

	// 3. Write them directly to files
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

	fmt.Println("\n=== 6. Direct Tags Collection Export (JSON/CSV) ===")
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
}

func proxyTester() {
	// paste your freshly gathered public proxies here (include http:// prefix)
	publicProxies := []string{
		"http://2.26.3.66:8080",
		// "http://2.26.17.187:8888",
		// "http://23.247.136.254:80",
	}

	// rotator := engine.NewProxyRotator(publicProxies)
	input := "http://127.0.0.1:8000/home"
	request := engine.InitRequest()

	request.ProxyRotator(publicProxies...)
	requestsCount := 4
	for i := 1; i <= requestsCount; i++ {
		body, err := request.Execute(input, "GET")

		if err != nil {
			fmt.Printf("[Req %d] ❌ Failed: %v (The public proxy might be dead)\n", i, err)
			return
		}

		// ipPattern := regexp.MustCompile(`\b(?:[0-9]{1,3}\.){3}[0-9]{1,3}\b`)

		// Find the first match inside the HTML string
		s, err := engine.InitDocument(body)
		if err != nil {
			fmt.Printf("[Req %d] ❌ Failed to parse HTML: %v\n", i, err)
			return
		}

		// callerIP := ipPattern.FindString(string(body))
		fmt.Printf("[Req %d] ✅ Success! Server Response:\n%s\n", i, s.Select("#ip-test").Print())
	}

	fmt.Println("🏁 Test finished.")
}

func testFromPost() {
	url := "http://127.0.0.1:8000/login"
	req := engine.InitRequest()

	// Get login page with the same client used for POST so cookies/session are shared.
	body, err := req.Execute(url, "GET")
	if err != nil {
		fmt.Println("Error loading login page:", err)
		return
	}

	scrape, err := engine.InitDocument(body)
	if err != nil {
		fmt.Println("Error parsing login page:", err)
		return
	}

	params := []*engine.Attribute{{Key: "name", Value: "_token"}}
	tokenInput := scrape.FindFirst("input", params)
	if tokenInput.Name == "" {
		fmt.Println("CSRF token input not found")
		return
	}

	token := ""
	for _, attr := range tokenInput.Attrs {
		if attr.Key == "value" {
			token = attr.Val
			break
		}
	}

	if token == "" {
		fmt.Println("CSRF token value missing")
		return
	}

	payload := map[string]string{
		"_token":   token,
		"email":    "superadmin@test.com",
		"password": "password123",
	}

	_, err = req.MakeFormPostRequest(url, "POST", payload)

	if err != nil {
		fmt.Println("Error making POST request:", err)
		return
	}

	fmt.Println("Login request sent successfully")
	fmt.Println("Cookies:", req.CookiesFor(url))

	body, err2 := req.Execute("http://127.0.0.1:8000/companies", "GET")
	if err2 != nil {
		fmt.Println("Error loading login page:", err2)
		return
	}

	scrape, err3 := engine.InitDocument(body)
	if err3 != nil {
		fmt.Println("Error parsing companies page:", err3)
		return
	}

	fmt.Println(scrape.FindAll("tr").Print())
}
