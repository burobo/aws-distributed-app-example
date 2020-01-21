package main

import (
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/google/uuid"
	"github.com/guregu/dynamo"
)

type RequestBody struct {
	InquirerName         string `json:"inquirer_name,omitempty"`
	InquirerEmailAddress string `json:"inquirer_email_address,omitempty"`
	InquiryDetails       string `json:"inquiry_details,omitempty"`
}

func main() {
	lambda.Start(handler)
}

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	reqBody := RequestBody{}

	err := json.Unmarshal([]byte(req.Body), &reqBody)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
		}, err
	}

	inquiryID, err := uuid.NewRandom()

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}

	inquiry := Inquiry{
		inquiryID.String(),
		reqBody.InquirerName,
		reqBody.InquirerEmailAddress,
		reqBody.InquiryDetails,
	}

	eventID, err := uuid.NewRandom()

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}

	inquired := newInquired(
		eventID.String(),
		inquiryID.String(),
		reqBody.InquirerName,
		reqBody.InquirerEmailAddress,
		reqBody.InquiryDetails)

	db := dynamo.New(session.New(), &aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})

	inquiriesTable := db.Table("Inquiries")
	err = inquiriesTable.Put(inquiry).Run()

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}

	InquiryEventsTable := db.Table("InquiryEvents")
	err = InquiryEventsTable.Put(inquired).Run()

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
	}, nil
}

type Inquiry struct {
	InquiryID            string // HashKey
	InquirerName         string `dynamo:"inquirer_name"`
	InquirerEmailAddress string `dynamo:"inquirer_email_address"`
	InquiryDetails       string `dynamo:"inquiry_details"`
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

func newInquired(
	EventID string,
	InquiryID string,
	InquirerName string,
	InquirerEmailAddress string,
	InquiryDetails string) Inquired {
	return Inquired{
		EventID,
		"Inquired",
		0,
		InquiryID,
		InquirerName,
		InquirerEmailAddress,
		InquiryDetails}
}
