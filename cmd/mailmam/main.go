package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"runtime"
	"sync"

	"github.com/toomore/simpleaws/sqs"
	"github.com/toomore/simpleaws/utils"
)

var AWSID = flag.String("awsid", os.Getenv("AWSID"), "AWSID")
var AWSKEY = flag.String("awskey", os.Getenv("AWSKEY"), "AWSKEY")
var SQSRegion = flag.String("sqsregion", "", "AWS SQS region")
var SQSURL = flag.String("sqsurl", "", "AWS SQS queue URL")
var SQSReceiverMax = flag.Int64("sqsrecmax", 10, "AWS SQS receiver max")
var sqsObject *sqs.SQS

func getQmsg(rmax int64) {
	if msg, err := sqsObject.Receive(rmax); err == nil {
		var wg sync.WaitGroup
		wg.Add(len(msg.Messages))
		for i, m := range msg.Messages {
			// Decode base64, ParseQuery
			if body, err := utils.Base64Decode([]byte(*m.Body)); err == nil {
				if bodymap, err := url.ParseQuery(string(body)); err == nil {
					go func(i int, bodymap url.Values, rh *string) {
						defer wg.Done()
						runtime.Gosched()
						log.Println(i, bodymap)
						sqsObject.Delete(rh)
					}(i, bodymap, m.ReceiptHandle)
				}
			}
		}
		wg.Wait()
	}
}

func main() {
	flag.Parse()
	if *AWSID == "" || *AWSID == "" {
		log.Fatal("Lost AWSID or AWSKEY")
	}
	if *SQSRegion != "" && *SQSURL != "" {
		sqsObject = sqs.New(*AWSID, *AWSKEY, *SQSRegion, *SQSURL)
	}
	getQmsg(*SQSReceiverMax)
}
