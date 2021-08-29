module "lambda_fetch_changed" {
  source = "./shared/lambda"

  function_name          = "${var.project}_fetch_changed"
  file_path              = "../lambdas/fetch-changed/function.zip"
  sns_subscription_topic = module.sns_topics["sns_event_streams_bridge"].sns_topic_arn
  sns_publish_topics     = [module.sns_topics["sns_raw_content_incoming"].sns_topic_arn]
  write_buckets          = [module.s3_buckets["raw-content"].s3_bucket_id]
  read_buckets           = []
}

module "lambda_fetch_schemaorg" {
  source = "./shared/lambda"

  function_name          = "${var.project}_fetch_schemaorg"
  file_path              = "../lambdas/fetch-schema.org/function.zip"
  sns_subscription_topic = module.sns_topics["sns_raw_content_incoming"].sns_topic_arn
  sns_publish_topics     = [module.sns_topics["sns_raw_content_schemaorg"].sns_topic_arn]
  write_buckets          = [module.s3_buckets["raw-content"].s3_bucket_id]
  read_buckets           = [module.s3_buckets["raw-content"].s3_bucket_id]
}

module "lambda_merge_schemaorg" {
  source = "./shared/lambda"

  function_name          = "${var.project}_merge_schemaorg"
  file_path              = "../lambdas/merge-schema.org/function.zip"
  sns_subscription_topic = module.sns_topics["sns_raw_content_schemaorg"].sns_topic_arn
  sns_publish_topics     = []
  write_buckets          = [module.s3_buckets["raw-content"].s3_bucket_id]
  read_buckets           = [module.s3_buckets["raw-content"].s3_bucket_id]
}

module "lambda_parsoid" {
  source = "./shared/lambda"

  function_name          = "${var.project}_transform_parsoid"
  file_path              = "../lambdas/transform-parsoid/function.zip"
  sns_subscription_topic = module.sns_topics["sns_raw_content_schemaorg"].sns_topic_arn
  sns_publish_topics     = [module.sns_topics["sns_node_published"].sns_topic_arn]
  write_buckets          = [module.s3_buckets["structured-content"].s3_bucket_id, module.s3_buckets["raw-content"].s3_bucket_id]
  read_buckets           = [module.s3_buckets["raw-content"].s3_bucket_id, module.s3_buckets["structured-content"].s3_bucket_id]
  dynamodb_tables = [{
    name = module.dynamodb_tables["node_names"].dynamodb_table_id,
    arn  = (module.dynamodb_tables["node_names"].dynamodb_table_arn)
    },
    {
      name = module.dynamodb_tables["page_titles"].dynamodb_table_id
      arn  = (module.dynamodb_tables["page_titles"].dynamodb_table_arn)
  }]
}