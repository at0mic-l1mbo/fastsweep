package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"github.com/fatih/color"
)

func pingHost(host string, wg *sync.WaitGroup, arq *os.File) {
	defer wg.Done()

	full_command := fmt.Sprintf("ping -c 1 -w 1 %s | grep '64 bytes'", host)
	_, err := exec.Command("/bin/bash", "-c", full_command).Output()
	if err == nil {
		message := fmt.Sprintf("[+] Response from %s\n", host)
		data := []byte(message)
		arq.Write(data)
	}
}

func beginPing(host string, arq *os.File) {
	var wg sync.WaitGroup

	for octet := 1; octet <= 256; octet++ {
		full_host := fmt.Sprintf("%s.%d", host, octet)
		wg.Add(1)
		go pingHost(full_host, &wg, arq)
	}

	wg.Wait()
}

func main() {
	if len(os.Args) < 2 {
		color.Red("[-] Invalid usage!")
		fmt.Println("Usage: go fast.go 192.168.0")
	} else {
		arq, err := os.Create("results.txt")
		if err != nil {
			log.Fatal("Failed to create file, verify your OS permissions!")
		}
		host := os.Args[1]
		color.Green("[+] Started ping sweep, could take a minute!")
		beginPing(host, arq)
		defer arq.Close()
		color.Green("[+] Results are saved in {results.txt}")
	}
}
