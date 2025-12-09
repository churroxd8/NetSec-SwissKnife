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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run banner_scanner.go <target_ip>")
		os.Exit(1)
	}

	target := os.Args[1]
	fmt.Printf("Scanning target: %s\n\n", target)

	var wg sync.WaitGroup

	ports := make(chan int, 100)

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go worker(ports, &wg, target)
	}

	for i := 1; i <= 9000; i++ {
		ports <- i
	}

	close(ports)
	wg.Wait()
	fmt.Println("\nScan complete.")
}

func worker(ports chan int, wg *sync.WaitGroup, target string) {
	defer wg.Done()

	for p := range ports {
		address := net.JoinHostPort(target, strconv.Itoa(p))

		conn, err := net.DialTimeout("tcp", address, 1*time.Second)

		if err != nil {
			continue
		}

		conn.SetReadDeadline(time.Now().Add(2 * time.Second))

		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		conn.Close()

		if err == nil && n > 0 {
			banner := strings.TrimSpace(string(buffer[:n]))
			fmt.Printf("[+] Port %d OPEN | Banner: %s\n", p, banner)
		} else {
			fmt.Printf("[+] Port %d OPEN | (No Banner)\n", p)
		}
	}
}
