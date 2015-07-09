package view

import (
	"html/template"
	"log"
	"net/http"

	"github.com/toomore/example_http/webserver/config"
	"github.com/toomore/example_http/webserver/session"
)

func init() {
	if tpl["/board"], err = template.ParseFiles("./template/base.html", "./template/board.html"); err != nil {
		log.Fatal("No template", err)
	}
}

func Board(w http.ResponseWriter, resp *http.Request) {
	var Session = session.New(config.SESSIONKEY, w, resp)
	var result outputdata
	if Session.Get("user") != "" {
		result.User = Session.Get("user")
	}
	tpl["/board"].Execute(w, result)
}
