package main

import (
	"log"
	"net/http"

	"github.com/toomore/example_http/webserver/config"
	"github.com/toomore/example_http/webserver/session"
	"github.com/toomore/example_http/webserver/view"
)

func needLogin(httpFunc func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, resp *http.Request) {
		log.Println("In wrapper", resp.UserAgent())
		var Session = session.New(config.SESSIONKEY, w, resp)
		log.Printf(">>>[%s] %+v", Session.Get("user"), Session.Hashvalues)
		if Session.Get("user") == "" {
			http.Redirect(w, resp, "/", http.StatusTemporaryRedirect)
		}
		httpFunc(w, resp)
	}
}

func main() {
	log.Println("Hello Toomore")

	http.HandleFunc("/", view.Index)
	http.HandleFunc("/login", view.Login)
	http.HandleFunc("/board", needLogin(view.Board))
	http.HandleFunc("/sendmail", needLogin(view.Sendmail))

	if err := http.ListenAndServe(":59122", nil); err != nil {
		log.Fatal(err)
	}
}
