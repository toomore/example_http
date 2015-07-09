package main

import (
	"log"

	"github.com/toomore/example_http/webserver/view"
)

func main() {
	log.Println("Hello Toomore")
	view.Run(":59122")
}
