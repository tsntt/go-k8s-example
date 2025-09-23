variable "region" {
  description = "The AWS region to deploy to."
  type        = string
  default     = "us-east-1"
}

variable "ami_id" {
  description = "The AMI ID for the EC2 instance."
  type        = string
  default     = "ami-0c55b159cbfafe1f0"
}