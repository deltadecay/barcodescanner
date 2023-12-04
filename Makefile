#git_hash := $(shell git rev-parse --short HEAD || echo 'development')
#version = ${git_hash}
version = 0.1.3

# Get current date time in UTC
current_time = $(shell date -u +"%Y-%m-%dT%H:%M:%S%Z")

# Add linker flags
linker_flags = '-s -w -X main.buildTime=${current_time} -X main.version=${version}'

# Build binaries for current OS and Linux
build:
	@echo "Building binaries..."
	GOOS=darwin GOARCH=arm64 go build -ldflags=${linker_flags} -o=./build/barcodescanner.darwin.arm64 .
	GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o=./build/barcodescanner.linux.amd64 .

clean:
	rm -rf build/*

package: build
	tar --directory build -czvf ./build/barcodescanner-${version}.darwin.arm64.tar.gz barcodescanner.darwin.arm64
	tar --directory build -czvf ./build/barcodescanner-${version}.linux.amd64.tar.gz barcodescanner.linux.amd64


.PHONY: build clean package
