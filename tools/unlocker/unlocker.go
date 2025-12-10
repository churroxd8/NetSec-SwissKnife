package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run unlocker.go <KEY_HEX_STRING>")
		os.Exit(1)
	}

	keyHex := os.Args[1]
	key, err := hex.DecodeString(keyHex)
	if err != nil {
		fmt.Println("Error decoding key. Did you copy it correctly?")
		os.Exit(1)
	}

	targetDir := "./sandbox"
	fmt.Println("[*] Starting decryption...")

	filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ".locked" {
			fmt.Printf(" -> Decrypting: %s\n", path)
			decryptFile(path, key)
		}
		return nil
	})

	fmt.Println("\n[+] FILES RESTORED")
}

func decryptFile(filename string, key []byte) {
	ciphertext, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("	Error reading %s: %v\n", filename, err)
		return
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		fmt.Println("	Error: File too short to be valid")
		return
	}

	nonce, encryptedData := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		fmt.Printf("	Error decrypting (Wrong Key?): %v\n", err)
		return
	}

	originalName := strings.TrimSuffix(filename, ".locked")
	err = os.WriteFile(originalName, plaintext, 0644)
	if err != nil {
		fmt.Printf("	Error writing %s: %v\n", originalName, err)
		return
	}

	os.Remove(filename)
}
