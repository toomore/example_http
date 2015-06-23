package main

import (
	"log"
	"net/http"
	"text/template"
)

func home(w http.ResponseWriter, resp *http.Request) {
	tpl["/"].Execute(w, nil)
	log.Println(resp.Header["User-Agent"])
	log.Println(resp.FormValue("q"))
}

var tpl map[string]*template.Template
var err error

func main() {
	log.Println("Hello Toomore")
	tpl = make(map[string]*template.Template)

	http.HandleFunc("/", home)
	if tpl["/"], err = template.ParseFiles("index.htm"); err != nil {
		log.Fatal("No template")
	}
	if err := http.ListenAndServe(":59122", nil); err != nil {
		log.Fatal(err)
	}
}
