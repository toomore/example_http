package view

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

func init() {
	if tpl["/campaign/create"], err = template.ParseFiles("./template/base.html", "./template/campaign_create.html"); err != nil {
		log.Fatal("No template", err)
	}
}

func Campaign(w http.ResponseWriter, resp *http.Request) {
	var action string
	var actionList = strings.Split(resp.URL.Path, "/")

	if len(actionList) > 2 {
		action = actionList[2]
	}

	switch resp.Method {
	case "GET":
		switch action {
		case "create":
			tpl["/campaign/create"].Execute(w, nil)

		default: //GET index.
			w.Write([]byte("Not implement."))
		}
	case "POST":
		switch action {
		case "create":
			resp.ParseForm()
			uploadTemplate(resp.FormFile("template"))
		default:
			w.Write([]byte("Not implement."))
		}
	default:
		w.Write([]byte("Not implement."))
	}
}
