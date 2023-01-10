variable "access_key" {
  type = string
}

variable "secret_key" {
  type = string
}

variable "region" {
  type = string
}

variable "ami" {
  type = string
}

variable "instance_type" {
  type = string
}



# Specificare il provider di cloud
provider "aws" {
  access_key = var.access_key
  secret_key = var.secret_key
  region     = var.region
}

# Creare una risorsa di tipo "aws_instance"
resource "aws_instance" "web" {
  # Scegliere un immagine AMI
  ami           = var.ami
  instance_type = var.instance_type

  # Aprire le porte SSH e HTTPS
  security_groups = ["ssh", "https"]

  # Assegnare un nome di host
  tags = {
    Name = "web"
  }
}

# Creare una risorsa di tipo "aws_security_group" per aprire la porta SSH
resource "aws_security_group" "ssh" {
  name        = "ssh"
  description = "Allow SSH access"

  # Aprire la porta 22 per il traffico in entrata
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Aprire anche la porta 6443
  ingress {
    from_port   = 6443
    to_port     = 6443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# Creare una risorsa di tipo "aws_security_group" per aprire la porta HTTPS
resource "aws_security_group" "https" {
  name        = "https"
  description = "Allow HTTPS access"

  # Aprire la porta 443 per il traffico in entrata
  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

output "public_ip" {
  value = aws_instance.web.public_ip
}
