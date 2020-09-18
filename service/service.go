package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/handlers"
	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/rs/cors"
	"github.com/wikimedia/phoenix/common"
	"github.com/wikimedia/phoenix/storage"
)

var (
	errorLogName  = flag.String("error-log", "error.log", "Path to the error log.")
	accessLogName = flag.String("access-log", "-", "Path to the access log.")

	// These values are passed in at build-time using -ldflags (see: Makefile)
	awsRegion          string
	dynamoDBPageTitles string
	dynamoDBNodeNames  string
	s3Bucket           string
)

// True if err is an awserr.Error, AND its code is s3.ErrCodeNoSuchKey, false otherwise.
func isS3NotFound(err error) bool {
	var s3err awserr.Error
	if errors.As(err, &s3err) {
		if s3err.Code() == s3.ErrCodeNoSuchKey {
			return true
		}
	}
	return false
}

func isErrNotFound(err error) bool {
	var nerr *storage.ErrNotFound
	if errors.As(err, &nerr) {
		return true
	}
	return false
}

// PageNameInput corresponds to a GraphQL input used by the Page query
type PageNameInput struct {
	Authority string
	Name      string
}

// NodeNameInput corresponds to a GraphQL input used by the Page query
type NodeNameInput struct {
	Authority string
	PageName  string
	Name      string
}

// RootResolver is the top-level GraphQL resolver
type RootResolver struct {
	Repository *storage.Repository
	Logger     *common.Logger
}

// Page returns a Page given its ID
func (r *RootResolver) Page(args struct {
	ID   *string
	Name *PageNameInput
}) (*PageResolver, error) {
	var page *common.Page
	var err error

	if args.Name != nil {
		// A page name was supplied
		if page, err = r.Repository.GetPageByName(args.Name.Authority, args.Name.Name); err != nil {
			// If err is of type storage.ErrNotFound, then this is not an error per say
			if isErrNotFound(err) {
				return nil, nil
			}

			r.Logger.Error("Unable to retrieve Page (authority=%s, name=%s): %s", args.Name.Authority, args.Name.Name, err)
			return nil, err
		}

		// A name argument was supplied and a matching page was found.  If an ID was also specified but it does NOT match
		// the page returned, then we return nil on the basis that we have no results that match all of the predicates supplied.
		if args.ID != nil {
			if page.ID != *args.ID {
				return nil, nil
			}
		}
	} else if args.ID != nil {
		// The page ID was supplied
		if page, err = r.Repository.GetPage(*args.ID); err != nil {
			// If this was an error returned by S3 (it is an awserr.Error), and its code is s3.ErrCodeNoSuchKey
			// then the object was simply not found (read: this is not an error per say).
			if isS3NotFound(err) {
				return nil, nil
			}

			r.Logger.Error("Unable to retrieve Page (ID=%s): %s", *args.ID, err)
			return nil, err
		}
	} else {
		// Neither a page ID or a name was supplied
		return nil, nil
	}

	return &PageResolver{page, r.Repository}, nil
}

// Node returns a Node given its ID
func (r *RootResolver) Node(args struct {
	ID   *string
	Name *NodeNameInput
}) (*NodeResolver, error) {
	var node *common.Node
	var err error

	if args.Name != nil {
		if node, err = r.Repository.GetNodeByName(args.Name.Authority, args.Name.PageName, args.Name.Name); err != nil {
			// If err is of type storage.ErrNotFound, then this is not an error per say
			if isErrNotFound(err) {
				return nil, nil
			}

			r.Logger.Error("Unable to retrieve Node (authority=%s, pageName=%s, name=%s): %s", args.Name.Authority, args.Name.PageName, args.Name.Name, err)
			return nil, err
		}
	} else if args.ID != nil {
		if node, err = r.Repository.GetNode(*args.ID); err != nil {
			// If this was an error returned by S3 (it is an awserr.Error) and its code is s3.ErrCodeNoSuchKey
			// then the object was simply not found (read: this is not an error per say).
			if isS3NotFound(err) {
				return nil, nil
			}

			r.Logger.Error("Unable to retrieve Node (ID=%s): %s", *args.ID, err)
			return nil, err
		}
	} else {
		// Neither a node ID or a name was supplied
		return nil, nil
	}

	return &NodeResolver{node}, nil
}

// PageResolver resolves a GraphQL page type
type PageResolver struct {
	p    *common.Page
	repo *storage.Repository
}

// ID resolves a page id attribute
func (r *PageResolver) ID() graphql.ID {
	return graphql.ID(r.p.ID)
}

// Name resolves a page name attribute
func (r *PageResolver) Name() string {
	return r.p.Name
}

// URL resolves a page url attribute
func (r *PageResolver) URL() string {
	return r.p.URL
}

// DateModified resolves a page dateModified attribute
func (r *PageResolver) DateModified() string {
	return r.p.DateModified.Format(time.RFC3339)
}

// HasPart resolves a page hasPart attribute
func (r *PageResolver) HasPart(args struct {
	Limit  *int32
	Offset *int32
}) ([]*NodeResolver, error) {
	var err error
	var node *common.Node
	var offset int32 = 0
	var resolvers = make([]*NodeResolver, 0)

	if args.Offset != nil {
		offset = *args.Offset
	}

	// TODO: This is slow; Consider adding concurrency
	for i, id := range r.p.HasPart[offset:] {
		if args.Limit != nil && (int32(i)+1) > *args.Limit {
			break
		}
		if node, err = r.repo.GetNode(id); err != nil {
			// If this was an error returned by S3 (it is an awserr.Error) and its code is s3.ErrCodeNoSuchKey
			// then the object was simply not found (read: this is not an error per say).
			if isS3NotFound(err) {
				return nil, nil
			}
			return nil, err
		}
		resolvers = append(resolvers, &NodeResolver{node})
	}

	return resolvers, nil
}

// About resolves a page about attribute
func (r *PageResolver) About(args struct{ Key *string }) []*TupleResolver {
	res := make([]*TupleResolver, 0)

	for k, v := range r.p.About {
		// If a key predicate has been supplied, then try to match it and return.
		if args.Key != nil {
			if *args.Key == k {
				res = append(res, &TupleResolver{key: k, val: v})
				break
			}
		} else {
			res = append(res, &TupleResolver{key: k, val: v})
		}
	}

	return res
}

// TupleResolver resolves a GraphQL Tuple type
type TupleResolver struct {
	key string
	val string
}

// Key resolves a tuple key attribute
func (r *TupleResolver) Key() string {
	return r.key
}

// Val resolves a tuple val attribute
func (r *TupleResolver) Val() string {
	return r.val
}

// NodeResolver resolves a GraphQL node type
type NodeResolver struct {
	n *common.Node
}

// ID resolves a node id attribute
func (r *NodeResolver) ID() graphql.ID {
	return graphql.ID(r.n.ID)
}

// Name resolves a node name attribute
func (r *NodeResolver) Name() string {
	return r.n.Name
}

// IsPartOf resolves a node isPartOf attribute
func (r *NodeResolver) IsPartOf() []graphql.ID {
	parents := make([]graphql.ID, len(r.n.IsPartOf))
	for _, id := range r.n.IsPartOf {
		parents = append(parents, graphql.ID(id))
	}
	return parents
}

// DateModified resolves a node dateModified attribute
func (r *NodeResolver) DateModified() string {
	return r.n.DateModified.Format(time.RFC3339)
}

// Unsafe resolves a node unsafe attribute
func (r *NodeResolver) Unsafe() string {
	return r.n.Unsafe
}

// Return configuration variables that are the union of defaults, and any values passed in the environment
func config() (region, titlesTable, namesTable, bucket string) {
	// Retrieve environment variables
	env := func(name string, def string) string {
		if v := os.Getenv(name); v != "" {
			return v
		}
		return def
	}

	region = env("AWS_REGION", awsRegion)
	titlesTable = env("AWS_DYNAMODB_PAGE_TITLES_TABLE", dynamoDBPageTitles)
	namesTable = env("AWS_DYNAMODB_NODE_NAMES_TABLE", dynamoDBNodeNames)
	bucket = env("AWS_BUCKET", s3Bucket)

	return region, titlesTable, namesTable, bucket
}

func main() {
	var accessLog io.Writer
	var b []byte
	var err error
	var logger *common.Logger
	var resolver *RootResolver
	var schema *graphql.Schema

	var region, titlesTable, namesTable, bucket = config()
	var awsSession = session.New(&aws.Config{Region: aws.String(region)})

	flag.Parse()

	// Setup the error logger
	if logger, err = common.NewFileLogger("INFO", *errorLogName); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open error log: %s", err)
		os.Exit(1)
	}

	// Setup the access log
	if *accessLogName == "-" {
		accessLog = os.Stdout
	} else {
		if accessLog, err = os.OpenFile(*accessLogName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666); err != nil {
			fmt.Fprintf(os.Stderr, "Error opening %s for writing: %s\n", *accessLogName, err)
			os.Exit(1)
		}
	}

	// FIXME: Keeping the schema in a separate file could be convenient during early development (VS Code's
	// GraphQL editor addon has been pretty handy, for example), but at some point practical considerations
	// will probably dictate that we move it to a string variable in-code.
	if b, err = ioutil.ReadFile("schema.gql"); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse schema.gql: %s", err)
		os.Exit(1)
	}

	resolver = &RootResolver{
		Repository: &storage.Repository{
			Store:  s3.New(awsSession),
			Index:  &storage.DynamoDBIndex{Client: dynamodb.New(awsSession), TitlesTable: titlesTable, NamesTable: namesTable},
			Bucket: bucket,
		},
		Logger: logger,
	}

	if schema, err = graphql.ParseSchema(string(b), resolver, graphql.UseFieldResolvers()); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing schema: %s", err)
		os.Exit(1)
	}

	handler := cors.Default().Handler(&relay.Handler{Schema: schema})

	log.Fatal(http.ListenAndServe(":8080", handlers.LoggingHandler(accessLog, handler)))
}
