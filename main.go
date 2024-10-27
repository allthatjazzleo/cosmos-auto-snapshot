package main

import (
	"flag"
	"log"
	"os"
	"path"

	"github.com/allthatjazzleo/auto-snapshot/internal"
)

func main() {
	// Define flags
	chainHomeDir := flag.String("chain-home", "", "The chain home directory")
	keepLocal := flag.Bool("keep-local", false, "Keep the local compressed file")
	uploaderType := flag.String("uploader", "s3", "Uploader type (s3/none) - set none to disable upload")
	nodeType := flag.String("node-type", "", "Node type (archive/default)")

	// Parse flags
	flag.Parse()

	// Validate flags
	if *chainHomeDir == "" {
		log.Println("chain-home flag is required")
		flag.Usage()
		os.Exit(1)
	}
	dataDir := path.Join(*chainHomeDir, "data")

	// Call checkVersion
	height, backendType, err := internal.CheckVersionAndDB(dataDir)
	if err != nil {
		log.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Print the result
	log.Printf("Height: %d\n", height)

	// get chain_id
	chainID, err := internal.GetChainID(*chainHomeDir)
	if err != nil {
		log.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	log.Printf("Chain ID: %s\n", chainID)

	var uploader internal.Uploader
	if *uploaderType == "s3" {
		config, err := internal.LoadConfig()
		if err != nil {
			log.Printf("Error loading AWS config: %v\n", err)
			os.Exit(1)
		}
		s3Bucket := os.Getenv("AWS_S3_BUCKET")
		if s3Bucket == "" {
			log.Printf("AWS_S3_BUCKET environment variable not set")
			os.Exit(1)
		}

		uploader, err = internal.NewS3Uploader(s3Bucket, config)
		if err != nil {
			log.Printf("Error creating S3 uploader: %v\n", err)
			os.Exit(1)
		}
	}

	err = internal.Compress(*chainHomeDir, chainID, backendType, height, uploader, *keepLocal, *nodeType)
	if err != nil {
		log.Printf("Error during compression: %v\n", err)
		os.Exit(1)
	}
}
