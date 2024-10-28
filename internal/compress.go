package internal

import (
	"archive/tar"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	dbm "github.com/cosmos/cosmos-db"
	"github.com/pierrec/lz4"
)

func Compress(homeDir string, chainID string, backendType dbm.BackendType, blockHeight int64, uploader Uploader, keepLocal bool, nodeType string) error {
	if nodeType != "" {
		nodeType = "_" + nodeType
	}
	filename := fmt.Sprintf("%s_%s%s_%d.tar.lz4", chainID, backendType, nodeType, blockHeight)

	var outputFile *os.File
	var err error
	if keepLocal || uploader == nil {
		outputFile, err = os.Create(filename)
		if err != nil {
			return err
		}
		defer outputFile.Close()
	}

	var writer io.Writer
	if uploader != nil {
		pr, pw := io.Pipe()
		defer pr.Close()

		go func() {
			defer pw.Close()
			if keepLocal {
				writer = io.MultiWriter(pw, outputFile)
			} else {
				writer = pw
			}
			err = compressToWriter(homeDir, writer)
			if err != nil {
				log.Printf("Error during compression: %v\n", err)
			}
		}()

		err = uploader.Upload(pr, filename)
		if err != nil {
			return err
		}
	} else {
		writer = outputFile
		err = compressToWriter(homeDir, writer)
		if err != nil {
			return err
		}
	}

	return nil
}

func compressToWriter(homeDir string, writer io.Writer) error {
	// Create an lz4 writer
	lz4Writer := lz4.NewWriter(writer)
	defer lz4Writer.Close()

	// Create a tar writer
	tarWriter := tar.NewWriter(lz4Writer)
	defer tarWriter.Close()

	// Define directories to include
	directories := []string{"data", "wasm"}

	for _, dir := range directories {
		dirPath := filepath.Join(homeDir, dir)
		// verify the directory exists
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			log.Printf("directory %s does not exist\n", dirPath)
			continue
		}

		err := filepath.Walk(dirPath, func(file string, fi os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Skip the wasm/cache directory
			if fi.IsDir() && filepath.Base(file) == "cache" && filepath.Dir(file) == filepath.Join(homeDir, "wasm", "wasm") {
				return filepath.SkipDir
			}

			// Create a tar header for the file
			header, err := tar.FileInfoHeader(fi, file)
			if err != nil {
				return err
			}
			header.Name, err = filepath.Rel(homeDir, file)
			if err != nil {
				return err
			}

			// Write the header to the tar archive
			if err := tarWriter.WriteHeader(header); err != nil {
				return err
			}

			// If it's a regular file, write its content to the tar archive
			if !fi.IsDir() {
				fileContent, err := os.Open(file)
				if err != nil {
					return err
				}
				defer fileContent.Close()

				if _, err := io.Copy(tarWriter, fileContent); err != nil {
					return err
				}
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}
