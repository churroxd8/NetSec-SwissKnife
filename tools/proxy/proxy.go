package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go proxy.go <LOCAL_PORT> <REMOTE_ADDR:PORT>")
		fmt.Println("Example: go proxy.go 8080 google.com:80")
		os.Exit(1)
	}

	localPort := os.Args[1]
	remoteAddr := os.Args[2]

	fmt.Printf("[*] Proxy starting on :%s -> Fowarding to %s\n", localPort, remoteAddr)

	listener, err := net.Listen("tcp", ":"+localPort)
	if err != nil {
		panic(err)
	}

	for {
		client, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}

		go handleConnection(client, remoteAddr)
	}
}

func handleConnection(client net.Conn, remoteAddr string) {
	target, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		fmt.Println("Remote connection failed:", err)
		client.Close()
		return
	}

	fmt.Printf("[+] New Connection: %s <-> %s <-> %s\n", client.RemoteAddr())

	go proxyCopy(client, target, "Client->Target")
	go proxyCopy(target, client, "Target->Client")
}

func proxyCopy(src, dst net.Conn, direction string) {
	defer src.Close()
	defer dst.Close()

	buf := make([]byte, 32*1024)

	for {
		n, err := src.Read(buf)
		if n > 0 {
			if n < 500 {
				fmt.Printf("\n--- %s (%d bytes) ---\n", direction, n)
				fmt.Println(hex.Dump(buf[:n]))
			}

			dst.Write(buf[:n])
		}

		if err != nil {
			if err != io.EOF {

			}
			break
		}
	}
}
