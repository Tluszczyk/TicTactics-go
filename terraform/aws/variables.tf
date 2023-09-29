variable "region" {
    type = string
    description = "AWS region"
    default     = "us-east-1"
}

variable "application" {
    type = string
    description = "Application name"
}

variable "deployment_type" {
    type = string
    description = "Must be one of the following: test, dev, prod"
    default     = "dev"

    validation {
        condition     = contains(["test", "dev", "prod"], var.deployment_type)
        error_message = "Valid values for var: deployment_type are (test, dev, prod)."
    } 
}

variable "files_path" {
    type = string
    description = "Path to files"
    default     = "./files"
}

variable "account_id" {
    type = string
    description = "id of the stack owner's account"
    default = "000000000000"
}

variable "endpoint_url" {
    type = string
    description = "endpoint url of the stack owner's account"
}

variable "database_deployment_option" {
    type = string
    description = "Must be one of the following: DYNAMO"
    default     = "DYNAMO"

    validation {
        condition     = contains(["DYNAMO"], var.database_deployment_option)
        error_message = "Valid values for var: deployment_type are (DYNAMO)."
    } 
}