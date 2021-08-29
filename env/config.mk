PHX_ACCOUNT_ID               = 717973338693
PHX_DEFAULT_REGION           = us-east-1
PHX_PREFIX                   = mntrxphnx

######
# SNS resources
######

# Topic that receives change events originating from the Wikimedia
# Event Streams service.
PHX_SNS_EVENT_STREAMS_BRIDGE      = $(PHX_PREFIX)_sns_event_streams_bridge

# Topic that receives events when new HTML is added to incoming (see
# PHX_S3_RAW_CONTENT_INCOMING).
PHX_SNS_RAW_CONTENT_INCOMING      = $(PHX_PREFIX)_sns_raw_content_incoming

# Topic that receives events when new linked-data (Wikidata) is added
# to the raw content store (see PHX_S3_RAW_CONTENT_WD_LINKED)
PHX_SNS_RAW_CONTENT_WD_LINKED     = $(PHX_PREFIX)_sns_raw_content_schemaorg

# Topic that receives events when new Node objects are added to the
# Structured Content Store
PHX_SNS_NODE_PUBLISHED            = $(PHX_PREFIX)_sns_node_published

# Topic that receives events when new Book objects are found in citations
# Structured Content Store
PHX_SNS_SOURCE_PARSE_PUBLISHED            = $(PHX_PREFIX)_sns_source_parse_published


######
# S3 resources
######

# The "raw content" bucket; Corresponds with uses of "raw content
# store" in the architecture documents.
PHX_S3_RAW_CONTENT_BUCKET        = $(PHX_PREFIX)-raw-content

# Folder where HTML documents of a corresponding revision are
# downloaded to after a change event is received.
PHX_S3_RAW_CONTENT_INCOMING      = incoming

# Folder where linked data (in the schema.org vocabulary) is stored.
PHX_S3_RAW_CONTENT_WD_LINKED     = schema.org

# Folder where HTML augmented with linked data is stored.
PHX_S3_RAW_CONTENT_LINKED_HTML   = linked-html

# The "structured content" bucket, where parsed and transformed data are
# stored in canonical format
PHX_S3_STRUCTURED_CONTENT_BUCKET = $(PHX_PREFIX)-structured-content


######
# Lambda resources
######

# Lambda function subscribed to Wikimedia Event Stream change events.
# Downloads the corresponding HTML (revision) and writes it to S3 (see
# PHX_S3_RAW_CONTENT_INCOMING)
PHX_LAMBDA_FETCH_CHANGED   = $(PHX_PREFIX)_lambda_fetch_changed

# Function invoked when new content has been added to incoming
# (PHX_SNS_RAW_CONTENT_INCOMING).  Downloads corresponding Wikidata
# information, constructs linked data (JSON-LD) in the schema.org
# vocabulary, and uploads to S3 (PHX_S3_RAW_CONTENT_WD_LINKED)
PHX_LAMBDA_FETCH_SCHEMAORG = $(PHX_PREFIX)_lambda_fetch_schemaorg

# Lambda subscribed to events that signal the creation of new Wikidata
# linked data (PHX_SNS_RAW_CONTENT_WD_LINKED).  Transforms the HTML
# from incoming (PHX_S3_RAW_CONTENT_INCOMING) to include the linked
# data (as JSON-LD), and uploads the result
# (PHX_S3_RAW_CONTENT_LINKED_HTML)
PHX_LAMBDA_MERGE_SCHEMAORG = $(PHX_PREFIX)_lambda_merge_schemaorg

# Lambda subscribed to events that signal the saving raw content to S3 
# storage, transforms raw content into canonical tructure and save to S3 
# storage (See PHX_S3_STRUCTURED_CONTENT_BUCKET)
PHX_LAMBDA_TRANSFORM_PARSOID = $(PHX_PREFIX)_lambda_transform_parsoid

# Lambda subscribed to events that signal that book source was found for
# a citation, gets book content into canonical tructure and save to S3 
# storage (See PHX_S3_STRUCTURED_CONTENT_BUCKET)
PHX_LAMBDA_BOOK_SOURCE = $(PHX_PREFIX)-lambda-book-source

# Lambda subscribed to events signaling that a new Node object has been stored.
# Retrieves related topic information for the Node, and stores the result.
PHX_LAMBDA_RELATED_TOPICS = $(PHX_PREFIX)-lambda-related-topics


######
# DynamoDB resources
######

# Table used to index page titles
PHX_DYNAMODB_PAGE_TITLES = $(PHX_PREFIX)_page_titles

# Table used to index node names
PHX_DYNAMODB_NODE_NAMES  = $(PHX_PREFIX)_node_names


######
# Rosette
######
PHX_ROSETTE_API_KEY = xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx


######
# Search index resources
######

# Use ../.config.mk to specify actual credentials
PHX_SEARCH_USERNAME   = wmfbooks-es-user
PHX_SEARCH_PASSWORD   = Wmfbooksespass_1

# Elasticsearch endpoint URL
PHX_SEARCH_ENDPOINT   = https://vpc-wmfbooks-es-e5tneoamgaiqiz67nbufjhrpla.us-east-1.es.amazonaws.com

# Elasticsearch index name for related topics indexing
PHX_SEARCH_IDX_TOPICS = topics