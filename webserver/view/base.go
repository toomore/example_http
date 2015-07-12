package view

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/toomore/example_http/webserver/config"
	"github.com/toomore/example_http/webserver/session"
	"github.com/toomore/simpleaws/s3"
)

var tpl map[string]*template.Template
var s3Object = s3.New(config.AWSID, config.AWSKEY, config.S3Region, config.S3Bucket)

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

func uploadTemplate(tplfile multipart.File, h *multipart.FileHeader, err error) (multipart.File, *multipart.FileHeader, error) {
	var filekey string
	var tpldata []byte

	if err == nil {
		defer tplfile.Close()
		if h.Header.Get("Content-Type") == "text/html" {
			tpldata, _ = ioutil.ReadAll(tplfile)
			log.Println(h.Filename, h.Header.Get("Content-Type"), tplfile)
			filekey = fmt.Sprintf("%s%s", config.S3Prefix, h.Filename)
			log.Println(s3Object.Put(filekey, bytes.NewReader(tpldata)))
		}
	}
	return tplfile, h, err
}
