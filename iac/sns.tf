module "sns_topics" {
  source   = "terraform-aws-modules/sns/aws"
  version  = "~> 3.0"
  for_each = var.sns_topics
  name     = "${var.project}_${each.key}"
}