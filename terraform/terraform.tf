terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "3.71.0"
    }
  }

  backend "remote" {
    organization = "alexa-cycling-skill"
    workspaces {
      name = "workspace"
    }
  }

  required_version = ">= 0.14.9"
}

provider "aws" {
  region = "eu-west-3"
}