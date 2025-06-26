# See rocksdb/README.md for instructions to update rocksdb version
FROM ghcr.io/strangelove-ventures/rocksdb:v7.10.2 AS rocksdb

FROM --platform=$BUILDPLATFORM golang:1.23-alpine AS builder

RUN apk add --update --no-cache\
    gcc\
    libc-dev\
    git\
    make\
    bash\
    g++\
    linux-headers\
    perl\
    snappy-dev\
    zlib-dev\
    bzip2-dev\
    lz4-dev\
    zstd-dev

ARG TARGETARCH
ARG BUILDARCH

RUN if [ "${TARGETARCH}" = "arm64" ] && [ "${BUILDARCH}" != "arm64" ]; then \
        wget -c https://storage.googleapis.com/strangelove-public/musl/aarch64-linux-musl-cross.tgz -O - | tar -xzvv --strip-components 1 -C /usr; \
    elif [ "${TARGETARCH}" = "amd64" ] && [ "${BUILDARCH}" != "amd64" ]; then \
        wget -c https://storage.googleapis.com/strangelove-public/musl/x86_64-linux-musl-cross.tgz -O - | tar -xzvv --strip-components 1 -C /usr; \
    fi

RUN set -eux;\
    if [ "${TARGETARCH}" = "arm64" ] && [ "${BUILDARCH}" != "arm64" ]; then \
        echo aarch64 > /etc/apk/arch;\
    elif [ "${TARGETARCH}" = "amd64" ] && [ "${BUILDARCH}" != "amd64" ]; then \
        echo x86_64 > /etc/apk/arch;\
    fi;\
    apk add --update --no-cache\
    snappy-static\
    zlib-static\
    bzip2-static\
    lz4-static\
    zstd-static\
    --allow-untrusted

# Install RocksDB headers and static library
COPY --from=rocksdb /rocksdb /rocksdb

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY *.go .
COPY internal/ internal/

ARG VERSION

RUN set -eux;\
    if [ "${TARGETARCH}" = "arm64" ] && [ "${BUILDARCH}" != "arm64" ]; then\
        export CC=aarch64-linux-musl-gcc CXX=aarch64-linux-musl-g++;\
    elif [ "${TARGETARCH}" = "amd64" ] && [ "${BUILDARCH}" != "amd64" ]; then\
        export CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++;\
    fi;\
    export  GOOS=linux \
            GOARCH=$TARGETARCH \
            CGO_ENABLED=1 \
            LDFLAGS='-linkmode external -extldflags "-static"' \
            CGO_CFLAGS="-I/rocksdb/include" \
            CGO_LDFLAGS="-L/rocksdb -L/usr/lib -L/lib -lrocksdb -lstdc++ -lm -lz -lbz2 -lsnappy -llz4 -lzstd";\
    go build -tags 'rocksdb pebbledb' -ldflags "$LDFLAGS" -a -o snapshot .


# Use alpine to source the latest CA certificates
FROM alpine:3 as alpine-3

# Build final image from scratch
FROM scratch

# Install trusted CA certificates
COPY --from=alpine-3 /etc/ssl/cert.pem /etc/ssl/cert.pem

WORKDIR /
USER 1025:1025
COPY --from=builder --chown=1025:1025 /workspace/snapshot .

ENTRYPOINT ["/snapshot"]
