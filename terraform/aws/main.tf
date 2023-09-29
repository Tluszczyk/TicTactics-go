terraform {
	required_providers {
		aws = {
			source  = "hashicorp/aws"
		}
	}
}

provider "aws" {
	region = "us-east-1"

    shared_credentials_files    = ["~/.aws/credentials"]
    profile                     = "localstack"

    default_tags {
        tags = {
            Application = var.application
            Deployment	= var.deployment_type
        }
    }
}