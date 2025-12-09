# NetSec Tools & Cloud Lab

## About the Project
This is a documentation of my journey on building custom cybersecurity tools from scratch. I wanted to understand how pre-made tools worked so I made my own versions in **Go (Golang)** and **Python** to understand how stuff like TCP, HTTP Ethernet work.

I also used **Terraform** to architect a custom "Vulnerable Lab" scenario in AWS to test these tools in a real-world cloud environment.

## The Toolkit
### 1. TCP/IP Sniffer (Python)
* **What it does:** Captures raw packets, decodes Ethernet/IP/TCP headers, and identifies flags (SYN, ACK, FIN).
* **Tech:** Raw Sockets, Struct unpacking.

### 2. Concurrent Port Scanner (Go)
* **What it does:** Scans a target subnet for open ports using high-concurrency Goroutines (a faster way than sequential scanning).
* **Features:** Service Banner Grabbing, Timeout handling.

### 3. HTTP Directory Buster (Go)
* **What it does:** Brute-force web paths to find a hidden admin panel or backdoor (e.g `backdoor.php`).
* **Tech:** HTTP Clinet, Worker Pools.

### 4. Remote Shell Client (Go)
* **What it does:** Connects to a remote PHP backdoor and provides an interactive terminal for Command Injection.

## Infrastructure as Code (Terraform)
I built two environments in AWS:
1. **The Fortress:** A locked-down EC2 instance allowing SSH only from my specific IP.
2. **The Target:** A deliberately vulnerable Web Server with an RCE injection point.

## Screenshots
![Sniffer Output](screenshots/sniffer_output.png)
*Capturing raw TCP handshake packets.*

![RCE Shell](screenshots/rce_shell.png)
*Remote Code Execution via custom Go client.*

## Disclaimer
These tools are for educational purposes and authorized security testing only. I am not responsible for your actions.