.PHONY: build

build:
	go build -tags 'rocksdb pebbledb' -o bin/snapshot .
