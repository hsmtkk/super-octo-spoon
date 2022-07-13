package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const s3PutEvent = "ObjectCreated:Put"

func main() {
	lambda.Start(handler)
}

func handler(sqsEvent events.SQSEvent) error {
	for _, msg := range sqsEvent.Records {
		ext, err := getExtFromMessage(msg)
		if err != nil {
			log.Printf("failed to get extension; %s", err)
		} else {
			log.Printf("extension of the file is %s", ext)
		}
	}
	return nil
}

func getExtFromMessage(e events.SQSMessage) (string, error) {
	log.Printf("SQS message: %s", e.Body)

	var snsEvent events.SNSEntity
	if err := json.Unmarshal([]byte(e.Body), &snsEvent); err != nil {
		return "", fmt.Errorf("failed to unmarshal; %s; %w", e.Body, err)
	}
	log.Printf("SNS message: %s", snsEvent.Message)

	if !strings.Contains(snsEvent.Message, s3PutEvent) {
		return "", nil
	}
	var s3event events.S3Event
	if err := json.Unmarshal([]byte(snsEvent.Message), &s3event); err != nil {
		return "", fmt.Errorf("failed to unmarshal; %s; %w", snsEvent.Message, err)
	}

	key, err := url.QueryUnescape(s3event.Records[0].S3.Object.Key)
	if err != nil {
		return "", fmt.Errorf("failed to unescape file name: %s", s3event.Records[0].S3.Object.Key)
	}

	return filepath.Ext(key), nil
}
