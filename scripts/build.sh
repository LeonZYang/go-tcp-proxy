#!/bin/bash

set -e
GOMAXPROCS=1 go test -timeout 90s ./...
GOMAXPROCS=4 go test -timeout 90s -race ./...
rm -fr ./bin
mkdir ./bin &>/dev/null
# no tests, but a build is something
for dir in $(find cmd/* -maxdepth 1 -type d); do
        echo "building $dir"
        GOOS=linux go build -o ./bin/$(basename $dir) ./$dir
done