package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	port := "8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	fmt.Printf("[*] HoneyPort active on port %s\n", port)
	fmt.Println("[*] Waiting for intruders...")

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Printf("Error starting listener: %v\n", err)
		os.Exit(1)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting: %v\n", err)
			continue
		}

		go handleIntruder(conn)
	}
}

func handleIntruder(conn net.Conn) {
	defer conn.Close()

	remoteAddr := conn.RemoteAddr().String()
	fmt.Printf("\n[!] ALERT: Intrusion detected from %s\n", remoteAddr)
	fakeBanner := "Admin Console v1.0 \r\nLogin: "
	conn.Write([]byte(fakeBanner))
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	buffer := make([]byte, 1024)

	n, _ := conn.Read(buffer)
	if n > 0 {
		fmt.Printf("	-> Attacker tried: %s\n", string(buffer[:n]))
	}
	fmt.Printf("		-> Connection closed for %s\n", remoteAddr)
}
