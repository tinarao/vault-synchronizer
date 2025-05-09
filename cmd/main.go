package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load .env vars: %s\n", err.Error())
	}

	vaultPath := os.Getenv("VAULT_PATH")
	if vaultPath == "" {
		log.Fatalf("vault path is not present in .env file\n")
	}

	fullPath, err := expandTilde(vaultPath)
	if err != nil {
		log.Fatalf("failed to expand ~: %s\n", err.Error())
	}

	contents, err := os.ReadDir(fullPath)
	if err != nil {
		log.Fatalf("failed to read vault dir contents: %s\n", err.Error())
	}

	for _, e := range contents {
		fmt.Printf("name: %s\n", e.Name())
	}
}

func expandTilde(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, path[1:]), nil
}
