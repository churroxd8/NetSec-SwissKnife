# Name: SecureFortress
# Author: churroxd8
# Description: A locked-down EC2 instance that allows only SSH from specific IP (Whitelisting)

terraform {
    required_providers {
      aws = {
        source = "hashicorp/aws"
        version = "~> 4.16"
      }
      http = {
        source = "hashicorp/http"
        version = "~> 3.0"
      }
    }
}

# 1. Configure the Provider (Targeting US East 1)
provider "aws" {
    region = "us-east-1"
}

# 2. Get Your Own Public IP Automatically
# This reaches out to an API to ask "What is my IP?"
data "http" "my_ip" {
    url = "http://ipv4.icanhazip.com"
}

# 3. The Digital Wall (Security Group)
resource "aws_security_group" "fortress_wall" {
    name = "fortress_security_group"
    description = "Allow SSH only from my specific IP"

    # Inbound Rule: Allow SSH (22) ONLY from your IP
    ingress {
        from_port = 22
        to_port = 22
        protocol = "tcp"
        cidr_blocks = ["${chomp(data.http.my_ip.response_body)}/32"] # <--- The Magic
    }

    # Outbound Rule: Allow the server to talk to the world (for updates)
    egress {
        from_port = 0
        to_port = 0
        protocol = "-1"
        cidr_blocks = ["0.0.0.0/0"]
    }
}

# 4. The Server (EC2 Instance)
resource "aws_instance" "app_server" {
    ami = "ami-053b0d53c279acc90" # Ubuntu Server 22.04 LTS (US-East-1)
    instance_type = "t3.micro" # Free Tier Eligible

    # Attach our firewall
    vpc_security_group_ids = [aws_security_group.fortress_wall.id]

    tags = {
        Name = "MySecureFortress"
    }
}

# 5. Output the Server's IP when done
output "instance_public_ip" {
    value = aws_instance.app_server.public_ip
}