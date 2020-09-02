module github.com/wikimedia/phoenix/storage

go 1.14

replace github.com/wikimedia/phoenix/common => ../common

require (
	github.com/aws/aws-sdk-go v1.34.17
	github.com/elastic/go-elasticsearch/v7 v7.9.0
	github.com/google/uuid v1.1.2
	github.com/spaolacci/murmur3 v1.1.0
	github.com/stretchr/testify v1.6.1
	github.com/wikimedia/phoenix/common v0.0.0-20200902184122-a02da52a642d
)
