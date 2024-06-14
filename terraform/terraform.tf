terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.54.0"
    }
  }

  backend "remote" {
    organization = "alexa-cycling-skill"
    workspaces {
      name = "workspace"
    }
  }

  required_version = "1.8.5"
}

provider "aws" {
  region = "eu-west-3"
}