module "s3_buckets" {
  source   = "terraform-aws-modules/s3-bucket/aws"
  for_each = var.s3_buckets

  bucket = "${var.project}-${each.key}"
  acl    = each.value.acl

  versioning = {
    enabled = each.value.versioning
  }

  // Make objects private if the bucket is 
  block_public_acls       = each.value.block_public
  block_public_policy     = each.value.block_public
  ignore_public_acls      = each.value.block_public
  restrict_public_buckets = each.value.block_public
}