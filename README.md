# Cosmos Auto Snapshot

Auto Snapshot is a tool designed to run as `StatefulJob` crd in `cosmos-operator` to automate the process of creating and uploading snapshots of blockchain data. This tool is particularly useful for blockchain nodes that need to periodically back up their data.

## Features

- **Automated Compression**: Compresses blockchain data directories into a single `.tar.lz4` file.
- **Uploader Support**: Supports uploading the compressed file to AWS S3/cloudflare R2/GCS/Azure.
- **Storage Optimization**: Can reduces the storage size on local machines by compressing the data and uploading it to a remote storage via I/O pipeline without local storage footprint.


## Installation

To install the tool, clone the repository and build the binary:

```sh
git clone https://github.com/allthatjazzleo/cosmos-auto-snapshot.git
cd cosmos-auto-snapshot
go build -o cosmos-auto-snapshot main.go
```

## Usage

Run the tool with the following command:
```sh
./cosmos-auto-snapshot --chain-home <path_to_chain_home> [--keep-local] [--uploader <uploader_type>] [--prefix <prefix>]
```

## Environment Variables

* AWS_S3_BUCKET: The S3 bucket where the snapshot will be uploaded if you are using S3 uploader.
* GCS_BUCKET: The GCS bucket where the snapshot will be uploaded if you are using GCS uploader.
* AZURE_STORAGE_ACCOUNT_URL: The Azure storage account URL if you are using Azure uploader.
* AZURE_STORAGE_CONTAINER: The Azure storage container name if you are using Azure uploader.
