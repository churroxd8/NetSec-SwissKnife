package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage:")
		fmt.Println("	1. Timestomp: go run ghost.go stomp <TARGET_FILE> <YYYY-MM-DD>")
		fmt.Println("	2. Shred:	  go run ghost.go shred <TARGET_FILE>")
		os.Exit(1)
	}

	mode := os.Args[1]
	target := os.Args[2]

	if mode == "stomp" {
		if len(os.Args) < 4 {
			fmt.Println("Usage: go run ghost.go stomp <FILE> <YYYY-MM-DD>")
			os.Exit(1)
		}
		newDateStr := os.Args[3]
		timestomp(target, newDateStr)

	} else if mode == "shred" {
		fmt.Printf("[*] Shredding %s (3 passes)...\n", target)
		shredFile(target)
		fmt.Println("[+] File obliterated")
	} else {
		fmt.Println("Unknown mode")
	}
}

func timestomp(filename string, dateStr string) {
	layout := "2006-01-02"
	newTime, err := time.ParseInLocation(layout, dateStr, time.Local)
	if err != nil {
		fmt.Println("Error parsing date. Use YYYY-MM-DD format")
		return
	}

	err = os.Chtimes(filename, newTime, newTime)
	if err != nil {
		fmt.Printf("Error changing timestamps: %v\n", err)
		return
	}

	fmt.Printf("[+] Flashback! %s is now dated %s\n", filename, dateStr)
}

func shredFile(filename string) {
	info, err := os.Stat(filename)
	if err != nil {
		fmt.Println("File not found")
		return
	}

	fileSize := info.Size()

	f, err := os.OpenFile(filename, os.O_RDWR, 0)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	for i := 1; i <= 3; i++ {
		fmt.Printf(" -> Pass %d: Overwriting bytes...\n", i)

		f.Seek(0, 0)

		garbage := make([]byte, fileSize)
		rand.Read(garbage)

		f.Write(garbage)
		f.Sync()
	}

	f.Close()
	os.Remove(filename)
}
