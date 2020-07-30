package main

import (
	"flag"
	"fmt"
	"github.com/vlamitin/secrets-provider/internal/persistence"
	"github.com/vlamitin/secrets-provider/internal/server"
	"golang.org/x/crypto/ssh/terminal"
	"strings"
	"syscall"
)

func main() {
	fmt.Println("Hello!")
	var port = flag.Int("port", 36663, "port for process to start on")
	var persist = flag.Bool("persist", false, "should secrets-provider save encrypted secrets to secrets.db file")

	flag.Parse()

	fmt.Println(*persist)

	cryptKey, err := readCryptKey()
	if err != nil {
		fmt.Printf("error when read key: %v\n", err)
	}

	persistence.SetCryptKey(cryptKey)
	server.InitAndListen("localhost", *port)
}

func readCryptKey() (string, error) {
	fmt.Print("Enter encryption/decryption key: ")
	byteKey, err := terminal.ReadPassword(syscall.Stdin)
	fmt.Print("\n")
	if err != nil {
		return "", fmt.Errorf("error when reading password from stdin: %w", err)
	}
	cryptKey := string(byteKey)

	return strings.TrimSpace(cryptKey), nil
}
