package view

import (
	"html/template"
	"log"
	"net/http"

	"github.com/toomore/example_http/webserver/config"
	"github.com/toomore/example_http/webserver/session"
)

var err error

func init() {
	if tpl["/"], err = template.ParseFiles("./template/base.html", "./template/index.html"); err != nil {
		log.Fatal("No template", err)
	}
}

func Index(w http.ResponseWriter, resp *http.Request) {
	if resp.URL.Path != "/" {
		log.Println(resp.URL.Path)
		http.NotFound(w, resp)
		return
	}
	var Session = session.New(config.SESSIONKEY, w, resp)
	var result outputdata
	if Session.Get("user") != "" {
		result.User = Session.Get("user")
	}
	tpl["/"].Execute(w, result)
	log.Println(resp.Header["User-Agent"])
	log.Println(resp.FormValue("q"))
	log.Println(resp.Cookie("session"))
}
