terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.27"
    }
  }

  backend "remote" {
    organization = "mentrix"
    workspaces {
      name = "phoenix"
    }
  }

  required_version = ">= 0.14.9"
}

provider "aws" {}

module "lambdas" {
  source = "./shared/lambda"
  count  = length(var.lambdas)

  function_name          = var.lambdas[count.index].name
  file_path              = var.lambdas[count.index].path
  sns_subscription_topic = var.lambdas[count.index].sns_subscription_topic
  sns_publish_topics     = var.lambdas[count.index].sns_publish_topics
  write_buckets          = var.lambdas[count.index].write_buckets
  read_buckets           = var.lambdas[count.index].read_buckets
  dynamodb_tables        = var.lambdas[count.index].dynamodb_tables
}

module "lambda_parsoid" {
  source = "./shared/lambda"

  function_name          = "phoenix-transform-parsoid"
  file_path              = "../lambdas/transform-parsoid/function.zip"
  sns_subscription_topic = "phoenix-sns-raw-content-schemaorg"
  sns_publish_topics     = ["phoenix-sns-node-published"]
  write_buckets          = ["phoenix-structured-content", "phoenix-raw-content"]
  read_buckets           = ["phoenix-raw-content", "phoenix-structured-content"]
  dynamodb_tables        = [module.dynamodb_tables["phoenix-node-names"].dynamodb_table_arn, module.dynamodb_tables["phoenix-page-titles"].dynamodb_table_arn]
}