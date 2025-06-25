package main

import (
	"flag"
	"log"
	"os"
	"path"

	"github.com/allthatjazzleo/cosmos-auto-snapshot/internal"
)

func main() {
	// Define flags
	chainHomeDir := flag.String("chain-home", "", "The chain home directory")
	keepLocal := flag.Bool("keep-local", false, "Keep the local compressed file")
	uploaderType := flag.String("uploader", "s3", "Uploader type (s3/gcs/azure/none) - set none to disable upload")
	nodeType := flag.String("node-type", "", "Node type (archive/default)")
	prefix := flag.String("prefix", "", "Optional prefix for the filename")

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
	switch *uploaderType {
	case "s3":
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
	case "gcs":
		gcsBucket := os.Getenv("GCS_BUCKET")
		if gcsBucket == "" {
			log.Printf("GCS_BUCKET environment variable not set")
			os.Exit(1)
		}

		uploader, err = internal.NewGCSUploader(gcsBucket)
		if err != nil {
			log.Printf("Error creating GCS uploader: %v\n", err)
			os.Exit(1)
		}
	case "azure":
		azureAccountURL := os.Getenv("AZURE_STORAGE_ACCOUNT_URL")
		azureContainer := os.Getenv("AZURE_STORAGE_CONTAINER")
		if azureAccountURL == "" || azureContainer == "" {
			log.Printf("AZURE_STORAGE_ACCOUNT_URL and AZURE_STORAGE_CONTAINER environment variables are required")
			os.Exit(1)
		}
		uploader, err = internal.NewAzureUploader(azureAccountURL, azureContainer)
		if err != nil {
			log.Printf("Error creating Azure uploader: %v\n", err)
			os.Exit(1)
		}
	}

	err = internal.Compress(*chainHomeDir, chainID, backendType, height, uploader, *keepLocal, *nodeType, *prefix)
	if err != nil {
		log.Printf("Error during compression: %v\n", err)
		os.Exit(1)
	}
}
