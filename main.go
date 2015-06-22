package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, resp *http.Request) {
	log.Println(resp)
	w.Write([]byte("In Home"))
}

func main() {
	log.Println("Hello Toomore")
	http.HandleFunc("/", home)
	if err := http.ListenAndServe(":59122", nil); err != nil {
		log.Fatal(err)
	}
}
