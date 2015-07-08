package config

import "os"

const (
	LOGINPWD  = "f9007add8286e2cb912d44cff34ac179"
	S3Bucket  = "toomore-aet"
	S3Region  = "us-east-1"
	SQSRegion = "ap-northeast-1"
	SQSURL    = "https://sqs.ap-northeast-1.amazonaws.com/271756324461/test_toomore"
)

var (
	AWSID      = os.Getenv("AWSID")
	AWSKEY     = os.Getenv("AWSKEY")
	SESSIONKEY = []byte("toomore")
)
