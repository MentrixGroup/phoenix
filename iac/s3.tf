module "s3_bucket" {
  source = "terraform-aws-modules/s3-bucket/aws"
  for_each = var.s3_buckets

  bucket = each.key
  acl    = each.value.acl

  versioning = {
    enabled = each.value.versioning
  }

}