// Project: banner_scanner.go
// Author: churroxd8
// Description: Scans a target subnet for open ports using high-concurrency goroutines

package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func worker(ports chan int, wg *sync.WaitGroup, target string) {
	defer wg.Done()

	for p := range ports {
		address := net.JoinHostPort(target, strconv.Itoa(p))
		conn, err := net.DialTimeout("tcp", address, 1*time.Second)

		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		conn.Close()

		if err == nil && n > 0 {
			banner := string(buffer[:n])
			banner = strings.TrimSpace(banner)
			fmt.Printf("[+] Port %d OPEN | Banner: %s\n", p, banner)
		} else {
			fmt.Printf("[+] Port %d OPEN | (No Banner)\n", p)
		}

	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run banner_scanner.go <target_ip>")
		fmt.Println("Example: go run banner_scanner.go scanme.nmap.org")
		os.Exit(1)
	}

	target := os.Args[1]

	fmt.Printf("Scanning target: %s\n\n", target)

	start := time.Now()
	ports := make(chan int, 100)
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go worker(ports, &wg, target)
	}

	// 1024 baby!
	for i := 1; i <= 1024; i++ {
		ports <- i
	}

	close(ports)
	wg.Wait()

	fmt.Printf("\nScan completed in %s\n", time.Since(start))
}
