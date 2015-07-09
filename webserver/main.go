package main

import (
	"flag"
	"log"

	"github.com/toomore/example_http/webserver/view"
)

var address = flag.String("a", ":59122", "Address")

func main() {
	flag.Parse()

	log.Println("Starting ...")
	view.Run(*address)
}
