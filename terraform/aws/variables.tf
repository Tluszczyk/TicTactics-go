variable "region" {
    type = string
    description = "AWS region"
    default     = "us-east-1"
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
    default     = "../../files"
}