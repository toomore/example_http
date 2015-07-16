// Mailman - To get queue from SQS and send by SES.
//
/*
Install:

	go install github.com/toomore/example_http/cmd/mailman

Usage:

	mailman [flags]

The flags are:

	-awsid
			AWSID (Default env: AWSID)
	-awskey
			AWSKEY (Default env: AWSKEY)
	-s3bucket
			S3 bucket (Default env: S3BUCKET)
	-s3region
			S3 bucket region (Default env: S3REGION)
	-sesregion
			SES region (Default env: SESREGION)
	-sqsrecmax
			SQS receive max messages limit (Default: 10)
	-sqsregion
			SQS region (Default env: SQSREGION)
	-sqsurl
			SQS URL (Default env: SQSURL)
	-ncpu
			CPU nums (Default: CPU nums)
	-retry
			Get SQS retry times (Default: 5)

*/
package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"net/mail"
	"net/url"
	"os"
	"runtime"
	"sync"
	"text/template"
	"time"

	webutils "github.com/toomore/example_http/webserver/utils"
	"github.com/toomore/simpleaws/s3"
	"github.com/toomore/simpleaws/ses"
	"github.com/toomore/simpleaws/sqs"
	"github.com/toomore/simpleaws/utils"
)

var (
	AWSID          = flag.String("awsid", os.Getenv("AWSID"), "AWSID")
	AWSKEY         = flag.String("awskey", os.Getenv("AWSKEY"), "AWSKEY")
	S3Bucket       = flag.String("s3bucket", os.Getenv("S3BUCKET"), "AWS S3 bucket")
	S3Region       = flag.String("s3region", os.Getenv("S3REGION"), "AWS S3 region")
	SESRegion      = flag.String("sesregion", os.Getenv("SESREGION"), "AWS SES region")
	SQSReceiverMax = flag.Int64("sqsrecmax", 10, "AWS SQS receiver max")
	SQSRegion      = flag.String("sqsregion", os.Getenv("SQSREGION"), "AWS SQS region")
	SQSURL         = flag.String("sqsurl", os.Getenv("SQSURL"), "AWS SQS queue URL")
	ncpu           = flag.Int("ncpu", runtime.NumCPU(), "指定 CPU 數量，預設為實際 CPU 數量")
	retry          = flag.Int64("retry", 5, "Get queue in zero to retry times")
	s3Object       *s3.S3
	sesObject      *ses.SES
	sqsObject      *sqs.SQS
	tplcache       map[string]string
)

func getQmsg(rmax int64) {
	var (
		delta     int64
		wg        sync.WaitGroup
		zerotimes int64
	)

Send:
	if msg, err := sqsObject.Receive(rmax); err == nil {
		wg.Add(len(msg.Messages))
		var mt = &sync.Mutex{}
		for i, m := range msg.Messages {
			// Decode base64, ParseQuery
			if body, err := utils.Base64Decode([]byte(*m.Body)); err == nil {
				if bodymap, err := url.ParseQuery(string(body)); err == nil {
					go func(i int, bodymap url.Values, rh *string) {
						defer wg.Done()
						runtime.Gosched()

						var ok bool
						var s3ouputbody string
						var tplpath = bodymap.Get("tplpath")

						if s3ouputbody, ok = tplcache[tplpath]; !ok {
							log.Println("No cache")
							mt.Lock()
							if s3ouputbody, ok = tplcache[tplpath]; !ok {
								if s3ouput, err := s3Object.Get(tplpath); err == nil {
									if s3ouputbyte, err := ioutil.ReadAll(s3ouput.Body); err == nil {
										tplcache[tplpath] = string(s3ouputbyte)
										s3ouputbody = tplcache[tplpath]
										log.Printf("save cache: %+v", s3ouput.Body)
										mt.Unlock()
									} else {
										mt.Unlock()
										return
									}
								} else {
									mt.Unlock()
									return
								}
							} else {
								log.Println("Pass by cache.")
								mt.Unlock()
							}
						}
						if tpl, err := template.New("tpl").Parse(s3ouputbody); err == nil {
							var tplcontent bytes.Buffer
							tpl.Execute(&tplcontent, webutils.Values2Map(bodymap))
							if sendresult, err := sesObject.Send(
								&mail.Address{
									Name:    bodymap.Get("sendername"),
									Address: bodymap.Get("senderemail"),
								},
								[]*mail.Address{
									&mail.Address{
										Name:    bodymap.Get("name"),
										Address: bodymap.Get("email")},
								},
								bodymap.Get("subject"), tplcontent.String()); err == nil {
								log.Println("[OK]", i, bodymap, sendresult)
								sqsObject.Delete(rh)
							} else {
								log.Println("[Error]", i, bodymap, err)
							}
						}
					}(i, bodymap, m.ReceiptHandle)
				} else {
					log.Println(err)
					wg.Done()
				}
			} else {
				log.Println(err)
				wg.Done()
			}
		}
		wg.Wait()
		if zerotimes < *retry {
			if len(msg.Messages) == 0 {
				zerotimes++
				log.Printf("In retry check. [%d]", zerotimes)
				delta = 1 << uint64(zerotimes)
				log.Printf("Retry in %d seconds.", delta)
				time.Sleep(time.Duration(delta) * time.Second)
			} else {
				zerotimes = 0
			}
			goto Send
		}
	}
	log.Println("Done")
}

func init() {
	tplcache = make(map[string]string)
}

func main() {
	flag.Parse()

	runtime.GOMAXPROCS(*ncpu)

	if *AWSID == "" || *AWSID == "" {
		log.Fatal("Lost AWSID or AWSKEY")
	}
	if *SQSRegion == "" || *SQSURL == "" {
		log.Fatal("Lost SQSRegion or SQSURL")
	}
	if *SESRegion == "" {
		log.Fatal("Lost SESRegion")
	}
	if *S3Region == "" || *S3Bucket == "" {
		log.Fatal("Lost S3Region or S3Bucket")
	}
	sqsObject = sqs.New(*AWSID, *AWSKEY, *SQSRegion, *SQSURL)
	sesObject = ses.New(*AWSID, *AWSKEY, *SESRegion)
	s3Object = s3.New(*AWSID, *AWSKEY, *S3Region, *S3Bucket)
	getQmsg(*SQSReceiverMax)
}
