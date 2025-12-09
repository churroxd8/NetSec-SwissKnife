// Project: dir_buster.go
// Author: churroxd8
// Description: Brute-forces web paths to find hidden admin panels or backdoors

package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

func worker(urlBase string, words chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for word := range words {
		targetURL := fmt.Sprintf("%s/%s", urlBase, word)
		resp, err := http.Get(targetURL)
		if err != nil {
			continue
		}
		if resp.StatusCode != 404 {
			fmt.Printf("[+] FOUND %s | Status %d\n", targetURL, resp.StatusCode)
		}
		resp.Body.Close()
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run dir_buster.go <http://target_ip>")
		os.Exit(1)
	}
	baseURL := os.Args[1]
	wordlistFile := "wordlist.txt"
	file, err := os.Open(wordlistFile)
	if err != nil {
		fmt.Println("Error opening wordlist.txt:", err)
		os.Exit(1)
	}
	defer file.Close()
	fmt.Printf("Starting Attack on %s...\n", baseURL)
	start := time.Now()
	words := make(chan string, 100)
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go worker(baseURL, words, &wg)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words <- scanner.Text()
	}
	close(words)
	wg.Wait()
	fmt.Printf("\nAttack completed in %s\n", time.Since(start))
}
