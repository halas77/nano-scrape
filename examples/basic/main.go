package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/halas77/nano-scrape/nano"
	"github.com/things-go/go-socks5"
)

func main() {
	// checkIP()
	// requestTest()
	stringTest()
	// exportDemo()
	// stringTest()

	// proxyTester()

	// testFromPost()
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
	scrape, err := nano.LoadDocument(input)
	fmt.Println(scrape.FindFirst("article").Print())

	// fmt.Println(scrape)

	// scrape, err := nano.InitDocument(input)

	// fmt.Println(scrape.FindFirst("main").Print())

	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return
	}

	// params := []*nano.Attribute{
	// 	{
	// 		Key:   "class",
	// 		Value: "details",
	// 	},
	// }

	// main := scrape.FindAll("div", params)
	// fmt.Println(main.Print())

	// scrape.Find("div", params, func(t nano.Tag) {
	// 	fmt.Println(t.Print())
	// })

	// scrape.Query(".price-tag", func(t *nano.Tag) {
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

	root, err := nano.InitDocument(input)

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
		brand := e.SelectFirst("span.brand")
		fmt.Printf("Electronic Item Brand: %s\n", brand.Print())
	}

	fmt.Println("\n=== 3. CSS + Deep Text Filtering ('Out of Stock') ===")
	// Note: CSS selectors don't natively support deep text matching in all variants,
	// but for now we simplify the example to use a selector.
	outOfStockItems := root.SelectAll("article")
	for _, item := range *outOfStockItems {
		if strings.Contains(item.Text(), "Out of Stock") {
			model := item.SelectFirst(".model")
			fmt.Printf("Item Out of Stock: %s\n", model.Print())
		}
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

	root, err := nano.InitDocument(input)
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
	jsonBytes, err := nano.ExportJSON(mappedData)
	if err != nil {
		fmt.Println("Error exporting JSON:", err)
		return
	}
	fmt.Println("Structured Mapped JSON:\n", string(jsonBytes))

	// 2. Export structured data to CSV
	csvStr, err := nano.ExportCSV(mappedData)
	if err != nil {
		fmt.Println("Error exporting CSV:", err)
		return
	}
	fmt.Println("Structured Mapped CSV:\n", csvStr)

	// 3. Export structured data to Markdown
	mdStr, err := nano.ExportMD(mappedData)
	if err != nil {
		fmt.Println("Error exporting MD:", err)
		return
	}
	fmt.Println("Structured Mapped Markdown Table:\n", mdStr)

	// 4. Write them directly to files
	err = nano.WriteMappedJSON("scraped_products.json", mappedData)
	if err != nil {
		fmt.Println("Error writing JSON file:", err)
	} else {
		fmt.Println("Successfully wrote structured data to: scraped_products.json")
	}

	err = nano.WriteMappedCSV("scraped_products.csv", mappedData)
	if err != nil {
		fmt.Println("Error writing CSV file:", err)
	} else {
		fmt.Println("Successfully wrote structured data to: scraped_products.csv")
	}

	err = nano.WriteMappedMD("scraped_products.md", mappedData)
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

func proxyTester() {
	publicProxies := []string{
		"http://2.26.3.66:8080",
		// "http://2.26.17.187:8888",
		// "http://23.247.136.254:80",
	}

	input := "http://127.0.0.1:8000/home"
	request := nano.NewClient()

	request.ProxyRotator(publicProxies...)
	requestsCount := 4
	for i := 1; i <= requestsCount; i++ {
		_, err := request.Get(input)

		if err != nil {
			fmt.Printf("[Req %d] ❌ Failed: %v (The public proxy might be dead)\n", i, err)
			return
		}

		fmt.Printf("[Req %d] ✅ Success! Server Response:\n%s\n", i, "<response body not captured>")
	}

	fmt.Println("🏁 Test finished.")
}

func testFromPost() {
	url := "http://127.0.0.1:8000/login"
	req := nano.NewClient()

	// Get login page with the same client used for POST so cookies/session are shared.
	body, err := req.Get(url)
	if err != nil {
		fmt.Println("Error loading login page:", err)
		return
	}

	scrape, err := nano.InitDocument(body)
	if err != nil {
		fmt.Println("Error parsing login page:", err)
		return
	}

	params := []*nano.Attribute{{Key: "name", Value: "_token"}}
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

	_, err = req.SendForm("POST", url, payload)

	if err != nil {
		fmt.Println("Error making POST request:", err)
		return
	}

	fmt.Println("Login request sent successfully")
	fmt.Println("Cookies:", req.CookiesFor(url))

	body, err2 := req.Get("http://127.0.0.1:8000/companies")
	if err2 != nil {
		fmt.Println("Error loading login page:", err2)
		return
	}

	scrape, err3 := nano.InitDocument(body)
	if err3 != nil {
		fmt.Println("Error parsing companies page:", err3)
		return
	}

	fmt.Println(scrape.FindAll("tr").Print())
}

type HttpBinResponse struct {
	Origin string `json:"origin"` // This holds the IP address seen by the server
}

// Helper to launch a local go-socks5 proxy server returning its full URL and a teardown function
func startMockSocks5Proxy() (string, func()) {
	server := socks5.NewServer()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(fmt.Sprintf("failed to start socks5 listener: %v", err))
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go server.ServeConn(conn)
		}
	}()

	proxyURL := fmt.Sprintf("socks5://%s", listener.Addr().String())
	return proxyURL, func() { listener.Close() }
}

type LaravelResponse struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
	Full string `json:"full"`
}

func checkIP() {
	input := "http://127.0.0.1:8000/home-ip"

	proxy1URL, cleanup1 := startMockSocks5Proxy()
	defer cleanup1()

	proxy2URL, cleanup2 := startMockSocks5Proxy()
	defer cleanup2()

	// 2. Build the list of active proxies
	publicProxies := []string{
		proxy1URL,
		proxy2URL,
	}

	fmt.Println(proxy1URL, " ", proxy2URL)

	// Assuming 'nano' is your custom internal package
	request := nano.NewClient()
	request.ProxyRotator(publicProxies...)

	requestsCount := 4
	for i := 0; i <= requestsCount; i++ {
		body, err := request.Get(input)

		if err != nil {
			// Changed 'return' to 'continue' so one dead proxy doesn't kill your entire loop
			fmt.Printf("[Req %d] ❌ Failed: %v (The public proxy might be dead)\n", i, err)
			continue
		}

		// 2. Read the response stream once, then unmarshal into multiple structs
		bodyBytes, err := io.ReadAll(body)
		if err != nil {
			fmt.Printf("[Req %d] ❌ Error reading response body: %v\n", i, err)
			continue
		}

		// 3. Decode into target struct
		var target HttpBinResponse
		if err := json.Unmarshal(bodyBytes, &target); err != nil {
			fmt.Printf("[Req %d] ❌ Error decoding JSON: %v\n", i, err)
			continue
		}

		// 4. Parse JSON returned by Laravel
		var laravelData LaravelResponse
		if err := json.Unmarshal(bodyBytes, &laravelData); err != nil {
			fmt.Printf("[Req %d] ❌ Failed parsing JSON: %v\nBody was: %s\n", i, err, string(bodyBytes))
			return
		}

		fmt.Printf("[Req %d] ✅ Laravel saw client IP: %s | Port: %d | Matches proxy list? %v\n",
			i, laravelData.IP, laravelData.Port)

		// 4. Extract the client IP from the struct and display it
		fmt.Printf("[Req %d] ✅ Success! Server detected your proxy IP as: %s\n", i, target.Origin)
	}

	fmt.Println("🏁 Test finished.")
}

func ScrapeTargetWithProxies() {
	proxies := []string{
		"http://proxy-us.example.com:3128",
		"http://proxy-eu.example.com:3128",
	}
	target := "https://example.com"

	// 1. Initialize your custom HTTP client wrapper
	client := nano.NewClient()

	// 2. Attach your proxy list using the ProxyRotator method.
	// This automatically configures the internal Go http.Transport to cycle proxies safely.
	client.ProxyRotator(proxies...)

	fmt.Println("🚀 Request pipeline initialized with proxy rotation...")

	// 3. Execute a network request to pull down the webpage data stream.
	// This returns an io.Reader (specifically a *bytes.Buffer), safely closing the network socket internally.
	responseStream, err := client.Get(target)
	if err != nil {
		log.Fatalf("❌ Network request failed: %v", err)
	}

	// 4. Pass the streaming data directly into the HTML Parser.
	// InitDocument automatically reads from the io.Reader stream.
	rootTag, err := nano.InitDocument(responseStream)
	if err != nil {
		log.Fatalf("❌ Failed to parse HTML content: %v", err)
	}

	fmt.Printf("✅ Successfully fetched and parsed: %s (Root Tag: %s)\n", target, rootTag.Name)

	// 5. Use the Tag query functions to extract data from the document tree.
	// Find the very first <h1> element on the page
	firstHeading := rootTag.FindFirst("h1")
	if firstHeading != nil {
		fmt.Printf("🎯 First Heading found on page: %s\n", firstHeading.Name)
	}

	// Stream and print all links found on the page using a callback
	fmt.Println("🔗 Listing all links found on the page:")
	rootTag.Find("a", nil, func(linkTag *nano.Tag) {
		// You can access link attributes via linkTag.Attrs
		fmt.Printf("   Found link element: %s\n", linkTag.Name)
	})
}
