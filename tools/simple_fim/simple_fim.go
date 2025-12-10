package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("	1. Learn:	go run simple_fim.go baseline <DIR>")
		fmt.Println("	2. Check:	go run simple_fim.go check <DIR>")
		os.Exit(1)
	}

	mode := os.Args[1]
	targetDir := "./sandbox" // Hardcodad for now
	if len(os.Args) > 2 {
		targetDir = os.Args[2]
	}

	baselineFile := "baseline.txt"

	if mode == "baseline" {
		fmt.Println("[*] Calculating baseline for %s...\n", targetDir)
		hashes := hashDirectory(targetDir)
		saveBaseline(baselineFile, hashes)
		fmt.Println("[+] Baseline saved to baseline.txt")

	} else if mode == "check" {
		fmt.Println("[*] Checking file integrity...")

		baseline := loadBaseline(baselineFile)
		current := hashDirectory(targetDir)

		changesFound := false

		for path, oldHash := range baseline {
			newHash, exists := current[path]
			if !exists {
				fmt.Printf("[!] ALERT: File DELETED: %s\n", path)
				changesFound = true
			} else if oldHash != newHash {
				fmt.Printf("[!] ALERT: File CHANGES: %s\n", path)
				changesFound = true
			}
		}

		if !changesFound {
			fmt.Println("[+] System Secure. No changes detected")
		}
	} else {
		fmt.Println("Unknown mode. Use `baseline` or `check`")
	}
}

func calculateHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func hashDirectory(root string) map[string]string {
	hashes := make(map[string]string)

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || strings.Contains(path, "baseline.txt") {
			return nil
		}

		hash, err := calculateHash(path)
		if err == nil {
			hashes[path] = hash
		}
		return nil
	})
	return hashes
}

func saveBaseline(filename string, hashes map[string]string) {
	f, _ := os.Create(filename)
	defer f.Close()
	for path, hash := range hashes {
		fmt.Fprintf(f, "%s|%s\n", path, hash)
	}
}

func loadBaseline(filename string) map[string]string {
	hashes := make(map[string]string)
	data, _ := os.ReadFile(filename)
	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) == 2 {
			hashes[parts[0]] = parts[1]
		}
	}
	return hashes
}
