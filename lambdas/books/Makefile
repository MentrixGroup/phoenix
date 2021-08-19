
# Repository configuration
include ../../env/config.mk
# User/dev overrides
-include ../../.config.mk

GOOS    := linux
BINARY  := main
SOURCES := main.go clientapi.go

# Configuration
LDFLAGS  = -X main.awsAccount=$(PHX_ACCOUNT_ID)
LDFLAGS += -X main.awsRegion=$(PHX_DEFAULT_REGION)
LDFLAGS += -X main.dynamoDBPageTitles=$(PHX_DYNAMODB_PAGE_TITLES)
LDFLAGS += -X main.dynamoDBNodeNames=$(PHX_DYNAMODB_NODE_NAMES)
LDFLAGS += -X main.s3StructuredContentBucket=$(PHX_S3_STRUCTURED_CONTENT_BUCKET)
LDFLAGS += -X main.snsSourceParse=$(PHX_SNS_SOURCE_PARSE_PUBLISHED)


build: clean
	GOOS=$(GOOS) go build -ldflags '$(LDFLAGS)' -o $(BINARY) $(SOURCES)
	zip function.zip main

deploy: build
	aws lambda update-function-code --function-name $(PHX_LAMBDA_BOOK_SOURCE) --zip-file fileb://function.zip

clean:
	rm -f function.zip $(BINARY)

test:
	go test

.PHONY: build clean deploy
