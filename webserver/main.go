package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/toomore/example_http/webserver/session"
)

var Session = session.New([]byte("toomore"))

func home(w http.ResponseWriter, resp *http.Request) {
	tpl["/"].Execute(w, nil)
	log.Println(resp.Header["User-Agent"])
	log.Println(resp.FormValue("q"))
	log.Println(resp.Cookie("session"))
}

func makeSession(value string, resp *http.Request) *http.Cookie {
	return &http.Cookie{
		Name:     "session",
		Value:    value,
		Path:     "/",
		Domain:   strings.Split(resp.Host, ":")[0],
		HttpOnly: true,
		//Expires: time.Now().Add(time.Duration(expires) * time.Second),
	}
}

func login(w http.ResponseWriter, resp *http.Request) {
	if resp.Method == "POST" {
		resp.ParseForm()
		Session.Set("age", "30")
		log.Printf("%+v", Session.Hashvalues)
		if resp.FormValue("email") != "" && resp.FormValue("pwd") != "" {
			//log.Println(resp.PostForm)
			hashpwd := md5.Sum([]byte(resp.FormValue("pwd")))
			if hashkey == fmt.Sprintf("%x", hashpwd) {
				http.SetCookie(w, makeSession("name=toomore", resp))
				Session.SetCookie(w, resp)
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

const hashkey = "f9007add8286e2cb912d44cff34ac179"

var (
	err error
	tpl map[string]*template.Template
)

func wrapper(httpFunc func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, resp *http.Request) {
		log.Println("In wrapper", resp.UserAgent())
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
