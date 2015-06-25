package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

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
		Domain:   resp.Host,
		HttpOnly: true,
		//Expires: time.Now().Add(time.Duration(expires) * time.Second),
	}
}

func login(w http.ResponseWriter, resp *http.Request) {
	if resp.Method == "POST" {
		resp.ParseForm()
		if resp.FormValue("email") != "" && resp.FormValue("pwd") != "" {
			//log.Println(resp.PostForm)
			hashpwd := md5.Sum([]byte(resp.FormValue("pwd")))
			if hashkey == fmt.Sprintf("%x", hashpwd) {
				http.SetCookie(w, makeSession("name=toomore", resp))
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

func main() {
	log.Println("Hello Toomore")
	tpl = make(map[string]*template.Template)

	http.HandleFunc("/", home)
	http.HandleFunc("/login", login)
	// template.ParseFiles need func.
	if tpl["/"], err = template.ParseFiles("./template/base.html", "./template/index.html"); err != nil {
		log.Fatal("No template", err)
	}
	if err := http.ListenAndServe(":59122", nil); err != nil {
		log.Fatal(err)
	}
}
