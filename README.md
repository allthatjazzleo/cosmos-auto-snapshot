# Auto Snapshot

Auto Snapshot is a tool designed to run as `StatefulJob` crd in `cosmos-operator` to automate the process of creating and uploading snapshots of blockchain data. This tool is particularly useful for blockchain nodes that need to periodically back up their data.

## Features

- **Automated Compression**: Compresses blockchain data directories into a single `.lz4` file.
- **Uploader Support**: Supports uploading the compressed file to AWS S3/cloudflare R2.
- **Configurable**: Allows configuration through command-line flags and environment variables.
- **Storage Optimization**: Can reduces the storage size on local machines by compressing the data and uploading it to a remote storage via I/O pipeline without local storage footprint.


## Installation

To install the tool, clone the repository and build the binary:

```sh
git clone https://github.com/allthatjazzleo/auto-snapshot.git
cd auto-snapshot
go build -o auto-snapshot main.go
```

## Usage

Run the tool with the following command:
```sh
./auto-snapshot --chain-home <path_to_chain_home> [--keep-local] [--uploader <uploader_type>]
```

## Environment Variables

* AWS_S3_BUCKET: The S3 bucket where the snapshot will be uploaded.
* AWS_ACCESS_KEY_ID: (Optional) Your AWS access key ID.
* AWS_SECRET_ACCESS_KEY: (Optional) Your AWS secret access key.
* AWS_S3_API_ENDPOINT: (Optional) Custom S3 API endpoint.