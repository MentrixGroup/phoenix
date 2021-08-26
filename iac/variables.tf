
# "phoenix-" is a prefix for all required resources and this should match the one from ".config.mk" in the root directory

variable "dynamodb_tables" {
  type = map(any)
  default = {
    node-names = {
      hash_key       = "Title",
      attribute_name = "Title",
      attribute_type = "S"
    },
    page-titles = {
      hash_key       = "Name",
      attribute_name = "Name",
      attribute_type = "S"
    }
  }
}

variable "s3_buckets" {
  type = map(any)
  default = {
    raw-content = {
      acl        = "private",
      versioning = false
    },
    structured-content = {
      acl        = "private",
      versioning = false
    }
  }
}

variable "sns_topics" {
  type = map(any)
  default = {
    sns-raw-content-schemaorg = {
      fifo_topic = false
    },
    sns-node-published = {
      fifo_topic = false
    },
    sns-raw-content-schemaorg = {
      fifo_topic = false
    },
    sns-raw-content-incoming = {
      fifo_topic = false
    },
    event-streams-bridge = {
      fifo_topic = false
    }

  }
}
variable "env_tag" {
  type    = string
  default = "staging"
}

variable "project" {
  type    = string
  default = "phoenix"
}

locals {
  # Common tags to be assigned to all resources
  common_tags = {
    Project     = var.project
    ManagedBy   = "Terraform"
    Environment = var.env_tag
  }

}