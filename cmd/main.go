package main

import (
	"log"
	"obsidian_backuper/internals/utils"
	"obsidian_backuper/internals/zip"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func getAllFiles(root string) ([]os.DirEntry, error) {
	var allFiles []os.DirEntry

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() && d.Name() == ".obsidian" {
			return filepath.SkipDir
		}

		if d.Name() != ".obsidian" {
			relPath, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}

			customEntry := &customDirEntry{
				DirEntry: d,
				path:     relPath,
				fullPath: path,
			}

			allFiles = append(allFiles, customEntry)
		}

		return nil
	})

	return allFiles, err
}

type customDirEntry struct {
	os.DirEntry
	path     string
	fullPath string
}

func (c *customDirEntry) Name() string {
	return c.path
}

func (c *customDirEntry) FullPath() string {
	return c.fullPath
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load .env vars: %s\n", err.Error())
	}

	outputDir := os.Getenv("OUTPUT_DIRECTORY")
	vaultPath := os.Getenv("VAULT_PATH")
	if vaultPath == "" || outputDir == "" {
		log.Fatalf("vault path is not present in .env file\n")
	}

	fullPath, err := utils.ExpandTilde(vaultPath)
	if err != nil {
		log.Fatalf("failed to expand ~: %s\n", err.Error())
	}

	contents, err := getAllFiles(fullPath)
	if err != nil {
		log.Fatalf("failed to read vault dir contents: %s\n", err.Error())
	}

	if err := zip.Zip(contents, outputDir); err != nil {
		log.Fatalf("failed to zip: %s\n", err.Error())
	}
}
