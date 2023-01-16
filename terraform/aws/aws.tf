variable "access_key" {
  type = string
}

variable "secret_key" {
  type = string
}

variable "region" {
  type = string
  default = "us-west-2"
}

variable "ami" {
  type = string
  default = "ami-0569ce8a44f2351be"
}

variable "instance_type" {
  type = string
  default = "c6in.2xlarge"
}


# Specify the cloud provider
provider "aws" {
  access_key = var.access_key
  secret_key = var.secret_key
  region     = var.region
}


resource "aws_instance" "web" {
  ami           = var.ami
  instance_type = var.instance_type
  security_groups = ["ssh", "https", "api"]
  tags = {
    Name = "web"
  }
}

# Create a resource of type "aws_security_group" to enable SSH
resource "aws_security_group" "ssh" {
  name        = "ssh"
  description = "Allow SSH access"

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 6443
    to_port     = 6443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group" "https" {
  name        = "https"
  description = "Allow HTTPS access"

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group" "API" {
  name        = "api"
  description = "Allow API access"

  ingress {
    from_port   = 6443 
    to_port     = 6443 
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}


output "public_ip" {
  value = aws_instance.web.public_ip
}
