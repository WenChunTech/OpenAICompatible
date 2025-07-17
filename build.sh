#!/usr/bin/env bash

GOOS=(darwin linux windows)
GOARCH=(amd64 arm64)
for os in ${GOOS[@]}; do
    for arch in ${GOARCH[@]}; do
        echo "Building for $os $arch"
        if [ $os = "windows" ]; then
            CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -ldflags "-s -w" -o openai_proxy_${os}_${arch}.exe main.go
        else
            CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -ldflags "-s -w" -o openai_proxy_${os}_${arch} main.go
        fi
    done
done