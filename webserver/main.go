package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/toomore/example_http/webserver/session"
)

const hashkey = "f9007add8286e2cb912d44cff34ac179"

var sessionkey = []byte("toomore")

func home(w http.ResponseWriter, resp *http.Request) {
	tpl["/"].Execute(w, nil)
	log.Println(resp.Header["User-Agent"])
	log.Println(resp.FormValue("q"))
	log.Println(resp.Cookie("session"))
}

func login(w http.ResponseWriter, resp *http.Request) {
	if resp.Method == "POST" {
		resp.ParseForm()
		if resp.FormValue("email") != "" && resp.FormValue("pwd") != "" {
			//log.Println(resp.PostForm)
			hashpwd := md5.Sum([]byte(resp.FormValue("pwd")))
			if hashkey == fmt.Sprintf("%x", hashpwd) {
				var Session = session.New(sessionkey, w, resp)
				Session.Set("user", resp.FormValue("email"))
				Session.Save()
				log.Printf("%+v", Session.Hashvalues)
				log.Println("Password Right!")
			}
			w.Write([]byte("In POST"))
		} else {
			http.Redirect(w, resp, "/", http.StatusSeeOther)
		}
	} else {
		http.Redirect(w, resp, "/", http.StatusSeeOther)
	}
}

var (
	err error
	tpl map[string]*template.Template
)

func wrapper(httpFunc func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, resp *http.Request) {
		log.Println("In wrapper", resp.UserAgent())
		var Session = session.New(sessionkey, w, resp)
		log.Printf(">>> %+v", Session.Hashvalues)
		httpFunc(w, resp)
	}
}

func main() {
	log.Println("Hello Toomore")
	tpl = make(map[string]*template.Template)

	http.HandleFunc("/", wrapper(home))
	http.HandleFunc("/login", wrapper(login))

	// template.ParseFiles need func.
	if tpl["/"], err = template.ParseFiles("./template/base.html", "./template/index.html"); err != nil {
		log.Fatal("No template", err)
	}

	if err := http.ListenAndServe(":59122", nil); err != nil {
		log.Fatal(err)
	}
}
