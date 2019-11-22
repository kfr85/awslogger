package main

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

var (
	Region      string
	Profile     string
	ProfilePath string
	LogGroup    string
	LogStream   string
	LogMessage  string
)

func GetCredential(file string, profile string, region string) *session.Session {
	creds := credentials.NewSharedCredentials(file, profile)
	_, err := creds.Get()
	if err != nil {
		panic(err)
	}

	sess, err := session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region)},
	)
	return sess
}

func Loop() {}

func PutLog(group string, stream string, message string, sess *session.Session) {
	svc := cloudwatchlogs.New(sess)
	token := ConfirmSequenceToken(group, stream, sess)

	configInputLogEvent := new(cloudwatchlogs.InputLogEvent).
		SetMessage(message).
		SetTimestamp(aws.TimeUnixMilli(time.Now()))

	configInputLogEventList := []*cloudwatchlogs.InputLogEvent{configInputLogEvent}
	configPutLogEventsInput := new(cloudwatchlogs.PutLogEventsInput)

	if token == "" {
		configPutLogEventsInput = configPutLogEventsInput.
			SetLogEvents(configInputLogEventList).
			SetLogGroupName(group).
			SetLogStreamName(stream)
	} else {
		configPutLogEventsInput = configPutLogEventsInput.
			SetLogEvents(configInputLogEventList).
			SetLogGroupName(group).
			SetLogStreamName(stream).
			SetSequenceToken(token)
	}

	_, err := svc.PutLogEvents(configPutLogEventsInput)

	if err != nil {
		panic(err)
	}
}

func ConfirmSequenceToken(group string, stream string, sess *session.Session) string {
	svc := cloudwatchlogs.New(sess)

	configDescribeLogStreamsInput := new(cloudwatchlogs.DescribeLogStreamsInput).
		SetLogGroupName(group).
		SetLogStreamNamePrefix(stream)

	res, err := svc.DescribeLogStreams(configDescribeLogStreamsInput)
	if err != nil {
		panic(err)
	}

	if res.LogStreams[0].UploadSequenceToken == nil {
		return ""
	} else {
		return *res.LogStreams[0].UploadSequenceToken
	}
}

func CreateLogStream() {
}

func main() {
	Profile = "trfm"
	Region = "ap-northeast-1"
	LogGroup = "test2"
	LogStream = "go"
	LogMessage = "hello go"
	fmt.Println(LogMessage)
	sess := GetCredential("", Profile, Region)
	PutLog(LogGroup, LogStream, LogMessage, sess)
}
