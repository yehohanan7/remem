#!/bin/bash
for GOOS in darwin linux; do
    export GOARCH=amd64
    export GOOS=$GOOS
    go build -o downloads/remem-$GOOS
done
