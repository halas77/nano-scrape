# Helper Functions

---

These are helper methods on the `*Tag` struct.

## `Print`

```go
func (tag *Tag) Print(depth ...uint16) string
```

Converts a `Tag` and all of its nested children elements back into a clean, formatted, and beautifully indented HTML string. It is primarily used for debugging, logging, or verifying exactly what HTML data your scraper has captured.

- **What it takes:** `depth`: _(Optional)_ An optional integer specifying the baseline indentation level for the output block. It defaults to `0`.
- **What it returns:** A formatted HTML string representation of the node tree. If the tag is empty or invalid, it returns `"Empty"`.

#### `💡 Usage Example`

When you use `FindFirst("div")` on the sample article, it grabs the very first `<div>` block (`class="info"`) and prints its entire structure with standardized alignment.

```go
package main

import (
	"fmt"
	"log"
)

func main() {
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

	scrape, err := engine.InitDocument(input)
	if err != nil {
		log.Fatal(err)
	}

	// 1. Locate the first <div> element (which is the one with class="info")
	firstDiv := scrape.FindFirst("div")

	// 2. Convert it to a formatted HTML string and display it
	fmt.Println(firstDiv.Print())
}
```

#### `📋 Printed Output`

```
<div class="info">
   <span class="price-tag">
      $899.99
   </span>
   <span class="stock">
      In Stock
   </span>
</div>
```

## `func (Tags) Print`

```go
func (ts Tags) Print(depth ...uint16) string
```

Converts an entire collection of elements (`Tags`) into a structured, indexed array-like layout. It iterates through every element in the collection, applies pretty-printing indentations to their internal structures, and surrounds each element with index indicators (`0: [ ... ]`). This is highly useful for checking multiple query results at once.

- **What it takes:** `depth`: _(Optional)_ An optional integer baseline that dictates the structural padding/indentation offset for the HTML elements inside the blocks.
- **What it returns:** A single formatted string rendering all elements in an indexed list format.

#### 💡 Usage Example

When you call `FindAll("div").Print()` on the sample article, it gathers all three `<div>` elements found anywhere in the tree structure and lists them out sequentially.

```go
package main

import (
	"fmt"
	"log"
)

func main() {
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

	scrape, err := engine.InitDocument(input)
	if err != nil {
		log.Fatal(err)
	}

	// 1. Locate ALL <div> elements on the page (creates a Tags collection)
	allDivs := scrape.FindAll("div")

	// 2. Format and print the entire collection together
	fmt.Println(allDivs.Print())
}
```

#### `📋 Printed Output`

```
0: [
<div class="info">
   <span class="price-tag">
      $899.99
   </span>
   <span class="stock">
      In Stock
   </span>
</div>
],
1: [
<div class="details">
   <div class="price-tag">
      ETB 909.99
   </div>
   <span class="brand">
      Apple
   </span>
   <span class="model">
      MacBook Air
   </span>
</div>
],
2: [
<div class="price-tag">
   ETB 909.99
</div>
]
```
