package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load .env vars: %s", err.Error())
	}

	vaultPath := os.Getenv("VAULT_PATH")
	if vaultPath == "" {
		log.Fatalf("vault path is not present in .env file")
	}

	fmt.Printf("vault path: %s\n", vaultPath)
}
