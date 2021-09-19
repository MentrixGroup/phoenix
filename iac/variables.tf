
# var.project is a prefix for all required resources and this should match the one from ".config.mk" in the root directory

variable "dynamodb_tables" {
  type = map(any)
  default = {
    node_names = {
      hash_key       = "Title"
      range_key      = "Authority"
      hk_attribute_name = "Title"
      hk_attribute_type = "S"
      rk_attribute_name = "Authority"
      rk_attribute_type = "S"
      billing_mode   = "PROVISIONED"
      read_capacity  = 2
      write_capacity = 2
    },
    page_titles = {
      hash_key       = "Name"
      hk_attribute_name = "Name"
      hk_attribute_type = "S"
      rk_attribute_name = "Authority"
      rk_attribute_type = "S"
      billing_mode   = "PROVISIONED"
      read_capacity  = 2
      write_capacity = 2
    }
  }
}

variable "s3_buckets" {
  type = map(any)
  default = {
    raw-content = {
      acl          = "private",
      block_public = true,
      versioning   = false
    },
    structured-content = {
      acl          = "private",
      block_public = true,
      versioning   = false
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