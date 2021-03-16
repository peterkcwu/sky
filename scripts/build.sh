#!/bin/bash

cd "$(dirname "$0")"

test -f bloboss && rm bloboss

echo "build..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bloboss ../cmd/bloboss/main.go

if test -f bloboss
then
    echo "upload file"
    rclone copy bloboss minio:bloboss/build
    echo "upload file end"
    echo "restart svr"
    md5sum bloboss
    #./build-on-dev.exe -ip 10.49.134.224 -cmd "cd /data/services/bloboss && sudo -u work bash reload.sh"
    exit 0
else
    echo 'build err' ; exit 50
fi
