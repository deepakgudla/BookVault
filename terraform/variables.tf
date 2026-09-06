variable "aws_region" {
  description = "AWS region for BookVault resources."
  type        = string
  default     = "ap-south-2"
}

variable "bucket_name" {
  description = "Globally unique S3 bucket name for product uploads."
  type        = string
}

variable "object_prefix" {
  description = "S3 key prefix the application may access."
  type        = string
  default     = "products/"
}

variable "queue_name" {
  description = "SQS queue name consumed by the notifier."
  type        = string
  default     = "bookvault-events"
}

variable "iam_role_name" {
  description = "Optional existing IAM role to receive the generated least-privilege policies."
  type        = string
  default     = ""
}