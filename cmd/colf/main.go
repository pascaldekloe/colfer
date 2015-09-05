package main

import (
	"log"
	"reflect"

	colfer "github.com/pascaldekloe/colfer/go"
)

func main() {
	// Hardcoded test definition
	d := &colfer.Object{
		Name: "tstobj",
		Fields: []*colfer.Field {
			{0, "b", reflect.Bool},
			{1, "i", reflect.Int},
			{13, "s", reflect.String},
		},
	}

	if err := d.Generate(); err != nil {
		log.Fatal(err)
	}
}
