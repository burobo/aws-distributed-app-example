package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	fmt.Printf("sqsEvent: %v", sqsEvent)
	fmt.Printf("sqsEvent.Records: %v", sqsEvent.Records)

	for _, message := range sqsEvent.Records {
		fmt.Printf("message.Body :%s \n", message.Body)

		var messageBody map[string]interface{}
		err := json.Unmarshal([]byte(message.Body), &messageBody)

		if err != nil {
			return err
		}

		var inquiry map[string]interface{}
		err = json.Unmarshal([]byte(messageBody["Message"].(string)), &inquiry)

		if err != nil {
			return err
		}

		err = send(os.Getenv("SES_EMAIL_FROM"), inquiry["InquirerEmailAddress"].(string), "Your inquiry was accepted.", "Your inquiry was accepted.", os.Getenv("SES_REGION"))

		if err != nil {
			return err
		}
	}

	return nil
}

func send(from string, to string, title string, body string, region string) error {
	svc := session.New(&aws.Config{
		Region: aws.String(region),
	})
	mailClient := ses.New(svc)
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(to),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(body),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(title),
			},
		},
		Source: aws.String(from),
	}
	output, err := mailClient.SendEmail(input)
	fmt.Printf("send mail output: %v", output)
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}

func main() {
	lambda.Start(handler)
}
