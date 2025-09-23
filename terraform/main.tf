# Configure the AWS provider.
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = "us-east-1"
}

# Simulate a Virtual Private Cloud (VPC) for our services
resource "aws_vpc" "app_vpc" {
  cidr_block = "10.0.0.0/16"
  tags = {
    Name = "my-testing-vpc"
  }
}

# Simulate an EC2 instance where our Kubernetes cluster would run.
resource "aws_instance" "go_api_instance" {
  ami           = "ami-0c55b159cbfafe1f0" # A common AMI for Amazon Linux 2
  instance_type = "t2.micro"
  vpc_security_group_ids = [aws_security_group.app_sg.id]

  tags = {
    Name = "go-api-server"
    Project = "Go-Microservice"
  }
}

# Simulate a security group that allows incoming traffic on specific ports.
resource "aws_security_group" "app_sg" {
  name        = "go_api_sg"
  vpc_id      = aws_vpc.app_vpc.id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}