## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 0.14.9 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | ~> 3.27 |

## Providers

No providers.

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_lambdas"></a> [lambdas](#module\_lambdas) | ./shared/lambda | n/a |

## Resources

No resources.

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_lambdas"></a> [lambdas](#input\_lambdas) | "phoenix-" is a prefix for all required resources and this should match the one from ".config.mk" in the root directory | <pre>list(object({<br>    name=string,<br>    path=string,<br>    sns_subscription_topic=string,<br>    sns_publish_topics=list(string)<br>    write_buckets=list(string)<br>    read_buckets=list(string)<br>    dynamodb_tables=list(string)<br>  }))</pre> | <pre>[<br>  {<br>    "dynamodb_tables": [],<br>    "name": "phoenix-fetch-changed",<br>    "path": "../lambdas/fetch-changed/function.zip",<br>    "read_buckets": [],<br>    "sns_publish_topics": [<br>      "phoenix-sns-raw-content-incoming"<br>    ],<br>    "sns_subscription_topic": "phoenix-event-streams-bridge",<br>    "write_buckets": [<br>      "phoenix-raw-content"<br>    ]<br>  },<br>  {<br>    "dynamodb_tables": [],<br>    "name": "phoenix-fetch-schemaorg",<br>    "path": "../lambdas/fetch-schema.org/function.zip",<br>    "read_buckets": [<br>      "phoenix-raw-content"<br>    ],<br>    "sns_publish_topics": [<br>      "phoenix-sns-raw-content-schemaorg"<br>    ],<br>    "sns_subscription_topic": "phoenix-sns-raw-content-incoming",<br>    "write_buckets": [<br>      "phoenix-raw-content"<br>    ]<br>  },<br>  {<br>    "dynamodb_tables": [],<br>    "name": "phoenix-merge-schemaorg",<br>    "path": "../lambdas/merge-schema.org/function.zip",<br>    "read_buckets": [<br>      "phoenix-raw-content"<br>    ],<br>    "sns_publish_topics": [],<br>    "sns_subscription_topic": "phoenix-sns-raw-content-schemaorg",<br>    "write_buckets": [<br>      "phoenix-raw-content"<br>    ]<br>  },<br>  {<br>    "dynamodb_tables": [<br>      "phoenix-node-names",<br>      "phoenix-page-titles"<br>    ],<br>    "name": "phoenix-transform-parsoid",<br>    "path": "../lambdas/transform-parsoid/function.zip",<br>    "read_buckets": [<br>      "phoenix-raw-content",<br>      "phoenix-structured-content"<br>    ],<br>    "sns_publish_topics": [<br>      "phoenix-sns-node-published"<br>    ],<br>    "sns_subscription_topic": "phoenix-sns-raw-content-schemaorg",<br>    "write_buckets": [<br>      "phoenix-structured-content",<br>      "phoenix-raw-content"<br>    ]<br>  }<br>]</pre> | no |

## Outputs

No outputs.
