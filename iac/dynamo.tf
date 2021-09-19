module "dynamodb_tables" {
  source = "terraform-aws-modules/dynamodb-table/aws"

  for_each = var.dynamodb_tables

  name     = "${var.project}-${each.key}"
  hash_key = each.value.hash_key
  range_key = each.value.range_key
  billing_mode   = each.value.billing_mode
  read_capacity  = each.value.read_capacity
  write_capacity = each.value.write_capacity


  attributes = [
    {
      name = each.value.hk_attribute_name
      type = each.value.hk_attribute_type
    },
    {
      name = each.value.rk_attribute_name
      type = each.value.rk_attribute_type
    }
  ]

}