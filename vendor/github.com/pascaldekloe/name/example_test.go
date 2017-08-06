package name_test

import (
	"fmt"
	"github.com/pascaldekloe/name"
)

func ExampleCamelCase() {
	fmt.Println(name.CamelCase("pascal case", true))
	fmt.Println(name.CamelCase("snake_to_camel AND CamelToCamel?", false))

	// Output:
	// PascalCase
	// snakeToCamelANDCamelToCamel
}

func ExampleDelimit() {
	// Garbage to Lisp-case:
	fmt.Println(name.Delimit("* All Hype is aGoodThing (TM)", '-'))

	// Builds a Java property key:
	fmt.Println(name.DotSeparated("WebCrawler#socketTimeout"))

	// Output:
	// all-hype-is-a-good-thing-TM
	// web.crawler.socket.timeout
}
