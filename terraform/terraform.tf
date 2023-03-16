terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "4.58.0"
    }
  }

  backend "remote" {
    organization = "alexa-cycling-skill"
    workspaces {
      name = "workspace"
    }
  }

  required_version = "1.4.2"
}

provider "aws" {
  region = "eu-west-3"
}