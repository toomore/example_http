package view

import (
	"html/template"
	"log"
	"net/http"

	"github.com/toomore/example_http/webserver/config"
	"github.com/toomore/example_http/webserver/session"
)

var tpl map[string]*template.Template

type outputdata struct {
	User string
}

func init() {
	tpl = make(map[string]*template.Template)
}

func NeedLogin(httpFunc func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
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

func Run(address string) {
	http.HandleFunc("/", Index)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/board", NeedLogin(Board))
	http.HandleFunc("/sendmail", NeedLogin(Sendmail))
	http.HandleFunc("/campaign", NeedLogin(Campaign))
	http.HandleFunc("/campaign/create", NeedLogin(Campaign))

	log.Printf("Address: %s\n", address)

	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatal(err)
	}
}
