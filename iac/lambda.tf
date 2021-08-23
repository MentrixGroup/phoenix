module "lambda_fetch-changed" {
  source = "./shared/lambda"

  function_name          = "${var.project}-fetch-schemaorg"
  file_path              = "../lambdas/fetch-changed/function.zip"
  sns_subscription_topic = module.sns_topics["event-streams-bridge"].sns_topic_arn
  sns_publish_topics     = [module.sns_topics["sns-raw-content-incoming"].sns_topic_arn]
  write_buckets          = [module.s3_buckets["raw-content"].s3_bucket_id]
  read_buckets           = []
  dynamodb_tables        = []
  common_tags            = local.common_tags
}

module "lambda_fetch-schemaorg" {
  source = "./shared/lambda"

  function_name          = "${var.project}-fetch-schemaorg"
  file_path              = "../lambdas/fetch-schema.org/function.zip"
  sns_subscription_topic = module.sns_topics["sns-raw-content-incoming"].sns_topic_arn
  sns_publish_topics     = [module.sns_topics["sns-raw-content-schemaorg"].sns_topic_arn]
  write_buckets          = [module.s3_buckets["raw-content"].s3_bucket_id]
  read_buckets           = [module.s3_buckets["raw-content"].s3_bucket_id]
  dynamodb_tables        = []
  common_tags            = local.common_tags
}

module "lambda_merge_schemaorg" {
  source = "./shared/lambda"

  function_name          = "${var.project}-merge-schemaorg"
  file_path              = "../lambdas/merge-schema.org/function.zip"
  sns_subscription_topic = module.sns_topics["sns-raw-content-schemaorg"].sns_topic_arn
  sns_publish_topics     = []
  write_buckets          = [module.s3_buckets["raw-content"].s3_bucket_id]
  read_buckets           = [module.s3_buckets["raw-content"].s3_bucket_id]
  dynamodb_tables        = []
  common_tags            = local.common_tags
}

module "lambda_parsoid" {
  source = "./shared/lambda"

  function_name          = "${var.project}-transform-parsoid"
  file_path              = "../lambdas/transform-parsoid/function.zip"
  sns_subscription_topic = module.sns_topics["sns-raw-content-schemaorg"].sns_topic_arn
  sns_publish_topics     = [module.sns_topics["sns-node-published"].sns_topic_arn]
  write_buckets          = [module.s3_buckets["structured-content"].s3_bucket_id, module.s3_buckets["raw-content"].s3_bucket_id]
  read_buckets           = [module.s3_buckets["raw-content"].s3_bucket_id, module.s3_buckets["structured-content"].s3_bucket_id]
  dynamodb_tables        = [module.dynamodb_tables["node-names"].dynamodb_table_arn, module.dynamodb_tables["page-titles"].dynamodb_table_arn]
  common_tags            = local.common_tags
}