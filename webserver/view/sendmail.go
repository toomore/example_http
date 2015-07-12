package view

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"

	"github.com/toomore/example_http/webserver/config"
	"github.com/toomore/example_http/webserver/session"
	"github.com/toomore/example_http/webserver/utils"
	"github.com/toomore/simpleaws/sqs"
)

func init() {
	if tpl["/sendmail"], err = template.ParseFiles("./template/base.html", "./template/sendmail.html"); err != nil {
		log.Fatal("No template", err)
	}
}

func Sendmail(w http.ResponseWriter, resp *http.Request) {
	switch resp.Method {
	case "GET":
		var Session = session.New(config.SESSIONKEY, w, resp)
		var result outputdata
		if Session.Get("user") != "" {
			result.User = Session.Get("user")
		}
		tpl["/sendmail"].Execute(w, result)
	case "POST":
		resp.ParseForm()

		_, h, _ := uploadTemplate(resp.FormFile("template"))
		//var tpldata []byte
		//var filekey string
		//if err == nil {
		//	defer tplfile.Close()
		//	if h.Header.Get("Content-Type") == "text/html" {
		//		tpldata, _ = ioutil.ReadAll(tplfile)
		//		log.Println(h.Filename, h.Header.Get("Content-Type"), tplfile)
		//		s3Object := s3.New(config.AWSID, config.AWSKEY,
		//			config.S3Region, config.S3Bucket)
		//		filekey = fmt.Sprintf("%s%s", config.S3Prefix, h.Filename)
		//		log.Println(s3Object.Put(filekey, bytes.NewReader(tpldata)))
		//	}
		//}
		var filekey string
		filekey = fmt.Sprintf("%s%s", config.S3Prefix, h.Filename)

		csvfile, h, err := resp.FormFile("csv")
		var csvValues []url.Values
		if err == nil {
			defer csvfile.Close()
			var queue []string
			if h.Header.Get("Content-Type") == "text/csv" {
				var sqsObject = sqs.New(config.AWSID, config.AWSKEY,
					config.SQSRegion, config.SQSURL)
				csvValues = utils.Map2ValuesMust(utils.CSV2map(csvfile))
				queue = make([]string, len(csvValues))
				for i, v := range csvValues {
					v.Set("tplpath", filekey)
					v.Set("sendername", resp.FormValue("sendername"))
					v.Set("senderemail", resp.FormValue("senderemail"))
					v.Set("subject", resp.FormValue("subject"))
					queue[i] = v.Encode()
				}
				sqsObject.SendBatch(queue)
				log.Println(h.Filename, h.Header.Get("Content-Type"))
			}
			fmt.Fprintf(w, "template: %s, Nums: %d", filekey, len(queue))
		}
	}
}
