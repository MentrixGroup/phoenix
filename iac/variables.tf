
# var.project is a prefix for all required resources and this should match the one from ".config.mk" in the root directory

variable "dynamodb_tables" {
  type = map(any)
  default = {
    node_names = {
      hash_key       = "Title",
      attribute_name = "Title",
      attribute_type = "S"
    },
    page_titles = {
      hash_key       = "Name",
      attribute_name = "Name",
      attribute_type = "S"
    }
  }
}

variable "s3_buckets" {
  type = map(any)
  default = {
    raw_content = {
      acl        = "private",
      versioning = false
    },
    structured_content = {
      acl        = "private",
      versioning = false
    }
  }
}

variable "sns_topics" {
  type = map(any)
  default = {
    sns_raw_content_schemaorg = {
      fifo_topic = false
    },
    sns_node_published = {
      fifo_topic = false
    },
    sns_raw_content_schemaorg = {
      fifo_topic = false
    },
    sns_raw_content_incoming = {
      fifo_topic = false
    },
    sns_event_streams_bridge = {
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
  default = ""
}

locals {
  # Common tags to be assigned to all resources
  common_tags = {
    Project     = var.project
    ManagedBy   = "Terraform"
    Environment = var.env_tag
  }

}