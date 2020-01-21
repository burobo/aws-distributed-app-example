# aws-distributed-app-example
This project is a distributed application example works on AWS.
The system sends a notification Email to inquirer when an inquirer sends a request on the API Gateway endpoint.

Publisher: API Gateway -> Lambda -> DynamoDB -> DynamoDB Streams -> Lambda(DomainEventPublisher) -> SNS
Subscriber: SQS -> Lambda(DomainEventSubscriber) -> SES

If additional flows e.g. "Send an Email to customer support." required, just add subscribers that subscribes "Inquired" Event.

# Usage
### Create s3 bucket
``` bash
make create-bucket
```
### Build go programs
``` bash
make build
```
### Zip files for windows user
```
make win-zip
```
### Deploy stack
``` bash
make deploy
```
### Remove stack
``` bash
make remove
```
