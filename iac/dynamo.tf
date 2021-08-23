module "dynamodb_tables" {
  source = "terraform-aws-modules/dynamodb-table/aws"

  for_each = var.dynamodb_tables

  name     = "${var.project}-${each.key}"
  hash_key = each.value.hash_key

  attributes = [
    {
      name = each.value.attribute_name
      type = each.value.attribute_type
    }
  ]

  tags = local.common_tags

}

output "dynamo_table_arns" {
  value = tomap({
    for key, dynamodb_table_arn in module.dynamodb_tables : key => dynamodb_table_arn
  })
}