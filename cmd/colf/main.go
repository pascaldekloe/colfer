package main

import (
	"log"

	colfer "github.com/pascaldekloe/colfer/go"
)

func main() {
	objects, err := colfer.ReadDefs()
	if err != nil {
		log.Fatal(err)
	}

	for _, o := range objects {
		if err := o.Generate(); err != nil {
			log.Fatal(err)
		}
	}
}
