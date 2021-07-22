// Package lambda to parse books
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/wikimedia/phoenix/common"
	"github.com/wikimedia/phoenix/storage"
)

var (
	// These values are passed in at build-time w/ -ldflags (see: Makefile)
	awsRegion                 string
	awsAccount                string
	dynamoDBPageTitles        string
	dynamoDBNodeNames         string
	s3StructuredContentBucket string
	s3RawBucket               string
	s3RawIncomeFolder         string
	s3RawLinkedFolder         string
	snsNodePublished          string

	debug bool = false
	log   *common.Logger
)

const (
	baseURL = "https://www.googleapis.com"
)

var logger *common.Logger

func snmatch(str string, matcher string) []string {
	match := make([]string, 0)
	rgx := regexp.MustCompile(fmt.Sprintf(`(?s)%s[0-9][ \t]*=[ \t]*(.*?)\|`, matcher))
	matches := rgx.FindAllStringSubmatch(str, -1)

	if len(matches) > 0 {
		for _, v := range matches {
			match = append(match, strings.TrimSpace(v[1]))
		}

		return match
	}

	return match
}

func smatch(s string, matcher string) string {
	r := regexp.MustCompile(fmt.Sprintf(`(?s)%s[ \t]*=[ \t]*(.*?)\|`, regexp.QuoteMeta(matcher)))
	matches := r.FindStringSubmatch(s)

	if len(matches) > 0 {
		return strings.TrimSpace(matches[1])
	}

	return ""
}

func Handler(ctx context.Context, event events.SNSEvent, cl *Client) {
	var book *common.Book
	var err error

	awsSession := session.New(&aws.Config{Region: aws.String(awsRegion)})
	s3client := s3.New(awsSession)

	repo := storage.Repository{
		Store:  s3client,
		Index:  &storage.DynamoDBIndex{Client: dynamodb.New(awsSession), TitlesTable: dynamoDBPageTitles, NamesTable: dynamoDBNodeNames},
		Bucket: s3StructuredContentBucket,
	}

	for _, record := range event.Records {
		msg := &common.SourceParseEvent{}
		if err := json.Unmarshal([]byte(record.SNS.Message), msg); err != nil {
			log.Error("Unable to deserialize message payload:", err)
			continue
		}

		if book, err = cl.GetBook(msg.isbn); err != nil {
			logger.Error("error making HTTP request: %w", err)
			return
		}

		if err = repo.PutBook(book); err != nil {
			logger.Error("error saving Book to the store: %w", err)
			return
		}
	}
}

func init() {
	// Determine logging level
	var level string = "ERROR"
	if v, ok := os.LookupEnv("LOG_LEVEL"); ok {
		level = v
	}

	// Initialize the logger
	log = common.NewLogger(level)
	log.Warn("%s LOGGING ENABLED (use LOG_LEVEL env var to configure)", common.LevelString(log.Level))

	log.Debug("AWS account ......................: %s", awsAccount)
	log.Debug("AWS region .......................: %s", awsRegion)
	log.Debug("DynamoDB page titles table .......: %s", dynamoDBPageTitles)
	log.Debug("DynamoDB node names table ........: %s", dynamoDBNodeNames)
	log.Debug("S3 structured content bucket .....: %s", s3StructuredContentBucket)
	log.Debug("S3 raw content bucket ............: %s", s3RawBucket)
	log.Debug("S3 raw content incoming folder ...: %s", s3RawIncomeFolder)
	log.Debug("S3 raw content linked folder .....: %s", s3RawLinkedFolder)
	log.Debug("SNS node published topic .........: %s", snsNodePublished)
}

func main() {
	lambda.Start(func(ctx context.Context, event events.SNSEvent) {
		cl := NewClient(baseURL)

		Handler(ctx, event, cl)
	})
}
