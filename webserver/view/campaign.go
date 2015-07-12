package view

import (
	"net/http"
	"strings"
)

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
			w.Write([]byte("In campaign create."))
			w.Write([]byte(resp.URL.Path))
			w.Write([]byte(action))

		default: //GET index.
			w.Write([]byte("Not implement."))
		}
	default:
		w.Write([]byte("Not implement."))
	}
}
