package main

import (
	"log"
	"net/http"
	"text/template"
)

func home(w http.ResponseWriter, resp *http.Request) {
	if t, err := template.ParseFiles("index.htm"); err == nil {
		t.Execute(w, nil)
	}
	log.Println(resp.Header["User-Agent"])
}

func main() {
	log.Println("Hello Toomore")
	http.HandleFunc("/", home)
	if err := http.ListenAndServe(":59122", nil); err != nil {
		log.Fatal(err)
	}
}
