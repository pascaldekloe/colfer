[![GoDoc](https://godoc.org/github.com/pascaldekloe/name?status.svg)](https://godoc.org/github.com/pascaldekloe/name)

Naming convention library for the Go programming language (golang).

This is free and unencumbered software released into the
[public domain](http://creativecommons.org/publicdomain/zero/1.0).
package name_test


### Inspiration

* `name.CamelCase("pascal case", true)` returns *PascalCase*
* `name.CamelCase("snake_to_camel AND CamelToCamel?", false)` returns *snakeToCamelANDCamelToCamel*
* `name.Delimit("* All Hype is aGoodThing (TM)", '-')` returns *all-hype-is-a-good-thing-TM*
* `name.DotSeparated("WebCrawler#socketTimeout")` returns *web.crawler.socket.timeout*
