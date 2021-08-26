terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.27"
    }
  }

  backend "remote" {
    organization = var.organization
    workspaces {
      name = var.project
    }
  }

  required_version = ">= 0.14.9"
}

provider "aws" {
  default_tags {
    tags = local.common_tags
  }
}

