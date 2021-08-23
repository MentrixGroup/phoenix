module "s3_buckets" {
  source = "terraform-aws-modules/s3-bucket/aws"
  for_each = var.s3_buckets

  bucket = "${var.project}-${each.key}"
  acl    = each.value.acl

  versioning = {
    enabled = each.value.versioning
  }

  tags = local.common_tags
}