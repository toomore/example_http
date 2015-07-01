package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/toomore/example_http/webserver/session"
)

const loginpwd = "f9007add8286e2cb912d44cff34ac179"

var sessionkey = []byte("toomore")

type outputdata struct {
	User string
}

func index(w http.ResponseWriter, resp *http.Request) {
	if resp.URL.Path != "/" {
		http.NotFound(w, resp)
		return
	}
	var Session = session.New(sessionkey, w, resp)
	var result outputdata
	if Session.Get("user") != "" {
		result.User = Session.Get("user")
	}
	tpl["/"].Execute(w, result)
	log.Println(resp.Header["User-Agent"])
	log.Println(resp.FormValue("q"))
	log.Println(resp.Cookie("session"))
}

func board(w http.ResponseWriter, resp *http.Request) {
	var Session = session.New(sessionkey, w, resp)
	var result outputdata
	if Session.Get("user") != "" {
		result.User = Session.Get("user")
	}
	tpl["/board"].Execute(w, result)
}

func sendmail(w http.ResponseWriter, resp *http.Request) {
	switch resp.Method {
	case "GET":
		var Session = session.New(sessionkey, w, resp)
		var result outputdata
		if Session.Get("user") != "" {
			result.User = Session.Get("user")
		}
		tpl["/sendmail"].Execute(w, result)
	case "POST":
		resp.ParseForm()

		if file, h, err := resp.FormFile("template"); err == nil {
			defer file.Close()
			if h.Header.Get("Content-Type") == "text/html" {
				log.Println(h.Filename, h.Header.Get("Content-Type"), file)
			}
		}
		if file, h, err := resp.FormFile("csv"); err == nil {
			defer file.Close()
			if h.Header.Get("Content-Type") == "text/csv" {
				log.Println(h.Filename, h.Header.Get("Content-Type"), file)
			}
		}
	}
}

func login(w http.ResponseWriter, resp *http.Request) {
	if resp.Method == "POST" {
		resp.ParseForm()
		if resp.FormValue("email") != "" && resp.FormValue("pwd") != "" {
			//log.Println(resp.PostForm)
			hashpwd := md5.Sum([]byte(resp.FormValue("pwd")))
			if loginpwd == fmt.Sprintf("%x", hashpwd) {
				var Session = session.New(sessionkey, w, resp)
				Session.Set("user", resp.FormValue("email"))
				Session.Save()
				log.Printf("%+v", Session.Hashvalues)
				log.Println("Password Right!")
			}
			http.Redirect(w, resp, "/board", http.StatusSeeOther)
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

func needLogin(httpFunc func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, resp *http.Request) {
		log.Println("In wrapper", resp.UserAgent())
		var Session = session.New(sessionkey, w, resp)
		log.Printf(">>>[%s] %+v", Session.Get("user"), Session.Hashvalues)
		if Session.Get("user") == "" {
			http.Redirect(w, resp, "/", http.StatusTemporaryRedirect)
		}
		httpFunc(w, resp)
	}
}

func main() {
	log.Println("Hello Toomore")
	tpl = make(map[string]*template.Template)

	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/board", needLogin(board))
	http.HandleFunc("/sendmail", needLogin(sendmail))

	// template.ParseFiles need func.
	if tpl["/"], err = template.ParseFiles("./template/base.html", "./template/index.html"); err != nil {
		log.Fatal("No template", err)
	}

	if tpl["/board"], err = template.ParseFiles("./template/base.html", "./template/board.html"); err != nil {
		log.Fatal("No template", err)
	}
	if tpl["/sendmail"], err = template.ParseFiles("./template/base.html", "./template/sendmail.html"); err != nil {
		log.Fatal("No template", err)
	}

	if err := http.ListenAndServe(":59122", nil); err != nil {
		log.Fatal(err)
	}
}
