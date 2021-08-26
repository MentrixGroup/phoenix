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
| <a name="module_dynamodb_tables"></a> [dynamodb\_tables](#module\_dynamodb\_tables) | terraform-aws-modules/dynamodb-table/aws | n/a |
| <a name="module_lambda_fetch-changed"></a> [lambda\_fetch-changed](#module\_lambda\_fetch-changed) | ./shared/lambda | n/a |
| <a name="module_lambda_fetch-schemaorg"></a> [lambda\_fetch-schemaorg](#module\_lambda\_fetch-schemaorg) | ./shared/lambda | n/a |
| <a name="module_lambda_merge_schemaorg"></a> [lambda\_merge\_schemaorg](#module\_lambda\_merge\_schemaorg) | ./shared/lambda | n/a |
| <a name="module_lambda_parsoid"></a> [lambda\_parsoid](#module\_lambda\_parsoid) | ./shared/lambda | n/a |
| <a name="module_s3_buckets"></a> [s3\_buckets](#module\_s3\_buckets) | terraform-aws-modules/s3-bucket/aws | n/a |
| <a name="module_sns_topics"></a> [sns\_topics](#module\_sns\_topics) | terraform-aws-modules/sns/aws | ~> 3.0 |

## Resources

No resources.

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_dynamodb_tables"></a> [dynamodb\_tables](#input\_dynamodb\_tables) | n/a | `map(any)` | <pre>{<br>  "node-names": {<br>    "attribute_name": "Title",<br>    "attribute_type": "S",<br>    "hash_key": "Title"<br>  },<br>  "page-titles": {<br>    "attribute_name": "Name",<br>    "attribute_type": "S",<br>    "hash_key": "Name"<br>  }<br>}</pre> | no |
| <a name="input_env_tag"></a> [env\_tag](#input\_env\_tag) | n/a | `string` | `"staging"` | no |
| <a name="input_organization"></a> [organization](#input\_organization) | n/a | `string` | `"mentrix"` | no |
| <a name="input_project"></a> [project](#input\_project) | n/a | `string` | `"phoenix"` | no |
| <a name="input_s3_buckets"></a> [s3\_buckets](#input\_s3\_buckets) | n/a | `map(any)` | <pre>{<br>  "raw-content": {<br>    "acl": "private",<br>    "versioning": false<br>  },<br>  "structured-content": {<br>    "acl": "private",<br>    "versioning": false<br>  }<br>}</pre> | no |
| <a name="input_sns_topics"></a> [sns\_topics](#input\_sns\_topics) | n/a | `map(any)` | <pre>{<br>  "event-streams-bridge": {<br>    "fifo_topic": false<br>  },<br>  "sns-node-published": {<br>    "fifo_topic": false<br>  },<br>  "sns-raw-content-incoming": {<br>    "fifo_topic": false<br>  },<br>  "sns-raw-content-schemaorg": {<br>    "fifo_topic": false<br>  }<br>}</pre> | no |

## Outputs

No outputs.
