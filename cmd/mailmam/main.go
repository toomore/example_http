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

	"github.com/toomore/simpleaws/s3"
	"github.com/toomore/simpleaws/ses"
	"github.com/toomore/simpleaws/sqs"
	"github.com/toomore/simpleaws/utils"
)

var (
	AWSID          = flag.String("awsid", os.Getenv("AWSID"), "AWSID")
	AWSKEY         = flag.String("awskey", os.Getenv("AWSKEY"), "AWSKEY")
	S3Bucket       = flag.String("s3bucket", "", "AWS S3 bucket")
	S3Region       = flag.String("s3region", "", "AWS S3 region")
	SESRegion      = flag.String("sesregion", "", "AWS SES region")
	SQSReceiverMax = flag.Int64("sqsrecmax", 10, "AWS SQS receiver max")
	SQSRegion      = flag.String("sqsregion", "", "AWS SQS region")
	SQSURL         = flag.String("sqsurl", "", "AWS SQS queue URL")
	ncpu           = flag.Int("ncpu", runtime.NumCPU(), "指定 CPU 數量，預設為實際 CPU 數量")
	retry          = flag.Int64("retry", 5, "Get queue in zero to retry times")
	s3Object       *s3.S3
	sesObject      *ses.SES
	sqsObject      *sqs.SQS
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
		for i, m := range msg.Messages {
			// Decode base64, ParseQuery
			if body, err := utils.Base64Decode([]byte(*m.Body)); err == nil {
				if bodymap, err := url.ParseQuery(string(body)); err == nil {
					go func(i int, bodymap url.Values, rh *string) {
						defer wg.Done()
						runtime.Gosched()
						if s3ouput, err := s3Object.Get(bodymap.Get("tplpath")); err == nil {
							s3ouputdata, _ := ioutil.ReadAll(s3ouput.Body)
							if tpl, err := template.New("tpl").Parse(string(s3ouputdata)); err == nil {
								var tplcontent bytes.Buffer
								tpl.Execute(&tplcontent, bodymap)
								log.Println(sesObject.Send(
									&mail.Address{
										Name:    bodymap.Get("sendername"),
										Address: bodymap.Get("senderemail"),
									},
									[]*mail.Address{
										&mail.Address{
											Name:    bodymap.Get("name"),
											Address: bodymap.Get("email")},
									},
									bodymap.Get("subject"), tplcontent.String()))
							}
						}
						log.Println(i, bodymap)
						sqsObject.Delete(rh)
					}(i, bodymap, m.ReceiptHandle)
				}
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
