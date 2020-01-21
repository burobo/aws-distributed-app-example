ENVIRONMENT := development
NAME        := your-ms-$(ENVIRONMENT)
BUCKET      := your-lambda-deploy-$(ENVIRONMENT)
OUT_DIR     := out
REGION      := ap-northeast-1

build:
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -installsuffix cgo -o $(OUT_DIR)/inquiry/inquire ./src/handler/inquiry/inquire
	zip $(OUT_DIR)/inquiry/inquire.zip $(OUT_DIR)/inquiry/inquire
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -installsuffix cgo -o $(OUT_DIR)/inquiry/inquiry-event-publisher ./src/handler/inquiry/inquiry-event-publisher
	zip $(OUT_DIR)/inquiry/inquiry-event-publisher.zip $(OUT_DIR)/inquiry/inquiry-event-publisher
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -installsuffix cgo -o $(OUT_DIR)/notification/notify-inquirer-of-confirmation ./src/handler/notification/notify-inquirer-of-confirmation
	zip $(OUT_DIR)/notification/notify-inquirer-of-confirmation.zip $(OUT_DIR)/notification/notify-inquirer-of-confirmation

win-zip:
	${GOPATH}\bin\build-lambda-zip.exe -o $(OUT_DIR)/inquiry/inquire.zip $(OUT_DIR)/inquiry/inquire
	${GOPATH}\bin\build-lambda-zip.exe -o $(OUT_DIR)/inquiry/inquiry-event-publisher.zip $(OUT_DIR)/inquiry/inquiry-event-publisher
	${GOPATH}\bin\build-lambda-zip.exe -o $(OUT_DIR)/notification/notify-inquirer-of-confirmation.zip $(OUT_DIR)/notification/notify-inquirer-of-confirmation

package:
	aws cloudformation package \
		--template-file template/$(ENVIRONMENT).yml \
		--s3-bucket $(BUCKET) \
		--s3-prefix $(NAME) \
		--output-template-file $(OUT_DIR)/.template.yml

deploy: package
	aws cloudformation deploy \
		--template-file $(OUT_DIR)/.template.yml \
		--stack-name $(NAME) \
		--capabilities CAPABILITY_IAM \
		--region $(REGION)
	
remove: 
	aws cloudformation delete-stack \
		--stack-name $(NAME) \
		--region $(REGION)

create-bucket:
	aws s3api create-bucket --bucket $(BUCKET) --region $(REGION) --create-bucket-configuration LocationConstraint=$(REGION)
