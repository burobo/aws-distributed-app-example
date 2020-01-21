package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/sns"
)

func main() {
	lambda.Start(handler)
}

func handler(e events.DynamoDBEvent) error {
	region := os.Getenv("AWS_REGION")

	sess := session.New(&aws.Config{
		Region: aws.String(region),
	})

	svc := sns.New(sess)

	log.Printf("publishing event: %v", e)
	for _, record := range e.Records {
		log.Printf("record : %v", record)
		if record.EventName == "INSERT" {
			var topicArn string
			var inquired Inquired

			if err := unmarshalStreamImage(record.Change.NewImage, &inquired); err != nil {
				return err
			}

			topicArn = os.Getenv("INQUIRED_TOPIC_ARN")
			eventByte, err := json.Marshal(inquired)

			if err != nil {
				return err
			}

			output, err := svc.Publish(&sns.PublishInput{
				TopicArn: aws.String(topicArn),
				Message:  aws.String(string(eventByte)),
			})

			if err != nil {
				return err
			}

			fmt.Printf("publish-output: %v", output)
		}
	}
	return nil
}

func unmarshalStreamImage(attribute map[string]events.DynamoDBAttributeValue, out interface{}) error {

	dbAttrMap := make(map[string]*dynamodb.AttributeValue)

	for k, v := range attribute {

		var dbAttr dynamodb.AttributeValue

		bytes, marshalErr := v.MarshalJSON()
		if marshalErr != nil {
			return marshalErr
		}

		json.Unmarshal(bytes, &dbAttr)
		dbAttrMap[k] = &dbAttr
	}

	return dynamodbattribute.UnmarshalMap(dbAttrMap, out)

}

type Inquired struct {
	EventID              string // HashKey
	EventType            string `json:"EventType"`
	EventVersion         int    `json:"EventVersion"`
	InquiryID            string // RangeKey
	InquirerName         string `json:"InquirerName"`
	InquirerEmailAddress string `json:"InquirerEmailAddress"`
	InquiryDetails       string `json:"InquiryDetails"`
}
