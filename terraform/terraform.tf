terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "6.13.0"
    }
  }

  backend "remote" {
    organization = "alexa-cycling-skill"
    workspaces {
      name = "workspace"
    }
  }

  required_version = "1.13.3"
}

provider "aws" {
  region = "eu-west-3"
}
