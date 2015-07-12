package config

import "os"

const (
	S3Bucket  = "toomore-aet"
	S3Prefix  = "tpl/"
	S3Region  = "us-east-1"
	SQSRegion = "ap-northeast-1"
	SQSURL    = "https://sqs.ap-northeast-1.amazonaws.com/271756324461/test_toomore"
)

var (
	AWSID      = os.Getenv("AWSID")
	AWSKEY     = os.Getenv("AWSKEY")
	LOGINPWD   = os.Getenv("LOGINPWD")
	SESSIONKEY = []byte(os.Getenv("SESSIONKEY"))
)
