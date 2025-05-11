package zip

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"obsidian_backuper/internals/utils"
	"os"
	"path/filepath"
	"time"
)

func Zip(files []os.DirEntry, outDir string) error {
	now := time.Now()
	localDateString := now.Format("2006-01-02T15:04:05")

	outFullDir, err := utils.ExpandTilde(outDir)
	if err != nil {
		return err
	}

	if _, err := os.Stat(outFullDir); os.IsNotExist(err) {
		if err := os.Mkdir(outFullDir, os.ModePerm); err != nil {
			log.Fatalf("output dir does not exists and it failed to create it. err: %s\n", err.Error())
		}
	}

	archiveFilePath := filepath.Join(outFullDir, fmt.Sprintf("bkg-%s", localDateString))

	archive, err := os.Create(archiveFilePath)
	if err != nil {
		return err
	}
	defer archive.Close()

	if err := createArchive(files, archive); err != nil {
		return err
	}

	fmt.Printf("archive path: %s\n", archiveFilePath)
	return nil
}

func Unzip(dir string) {}

//

func createArchive(files []os.DirEntry, buf io.Writer) error {
	gw := gzip.NewWriter(buf)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	for _, file := range files {
		err := addToArchive(tw, file)
		if err != nil {
			return err
		}
	}

	fmt.Printf("successfully created archive with %d files\n", len(files))

	return nil
}

func addToArchive(tw *tar.Writer, file os.DirEntry) error {
	info, err := file.Info()
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}

	header.Name = file.Name()

	err = tw.WriteHeader(header)
	if err != nil {
		return err
	}

	if file.IsDir() {
		return nil
	}

	var fullPath string
	if customEntry, ok := file.(interface{ FullPath() string }); ok {
		fullPath = customEntry.FullPath()
	} else {
		fullPath = filepath.Join(os.Getenv("VAULT_PATH"), file.Name())
	}

	fileHandle, err := os.Open(fullPath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", fullPath, err)
	}
	defer fileHandle.Close()

	_, err = io.Copy(tw, fileHandle)
	if err != nil {
		return fmt.Errorf("failed to copy file %s: %w", fullPath, err)
	}

	return nil
}
