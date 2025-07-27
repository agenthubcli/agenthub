# AGENTHUB CLI Build System

.PHONY: build test clean install cross-compile package help

# Variables
BINARY_NAME=agenthub
VERSION?=0.2.0
BUILD_DIR=dist
LDFLAGS=-ldflags "-X main.version=${VERSION}"

# Default target
all: build

# Build for current platform
build:
	@echo "Building agenthub..."
	go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME} main.go

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf ${BUILD_DIR}
	go clean

# Install locally
install: build
	@echo "Installing agenthub..."
	sudo cp ${BUILD_DIR}/${BINARY_NAME} /usr/local/bin/

# Cross-compile for all platforms
cross-compile: clean
	@echo "Cross-compiling for multiple platforms..."
	
	# macOS (Intel)
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/darwin-amd64/${BINARY_NAME} main.go
	
	# macOS (Apple Silicon)
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o ${BUILD_DIR}/darwin-arm64/${BINARY_NAME} main.go
	
	# Linux (Intel)
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/linux-amd64/${BINARY_NAME} main.go
	
	# Linux (ARM)
	GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o ${BUILD_DIR}/linux-arm64/${BINARY_NAME} main.go
	
	# Windows (Intel)
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/windows-amd64/${BINARY_NAME}.exe main.go
	
	@echo "Cross-compilation complete!"

# Package for distribution
package-unix: cross-compile
	@echo "Creating distribution packages..."
	
	# Create tarballs for Unix-like systems
	cd ${BUILD_DIR}/darwin-amd64 && tar -czf ../agenthub-${VERSION}-darwin-amd64.tar.gz ${BINARY_NAME}
	cd ${BUILD_DIR}/darwin-arm64 && tar -czf ../agenthub-${VERSION}-darwin-arm64.tar.gz ${BINARY_NAME}
	cd ${BUILD_DIR}/linux-amd64 && tar -czf ../agenthub-${VERSION}-linux-amd64.tar.gz ${BINARY_NAME}
	cd ${BUILD_DIR}/linux-arm64 && tar -czf ../agenthub-${VERSION}-linux-arm64.tar.gz ${BINARY_NAME}
	
	# Create zip for Windows
	cd ${BUILD_DIR}/windows-amd64 && zip ../agenthub-${VERSION}-windows-amd64.zip ${BINARY_NAME}.exe
	
	@echo "Packages created in ${BUILD_DIR}/"

package-release: cross-compile
	@mkdir -p release
	@echo "Packaging binaries for release..."
	@for os in darwin-amd64 darwin-arm64 linux-amd64 windows-amd64; do \
		zip -j release/${BINARY_NAME}-${VERSION}-$$os.zip ${BUILD_DIR}/$$os/${BINARY_NAME}*; \
		shasum -a 256 release/${BINARY_NAME}-${VERSION}-$$os.zip > release/${BINARY_NAME}-${VERSION}-$$os.zip.sha256; \
	done


# Development server (if applicable)
dev:
	@echo "Starting development mode..."
	go run main.go

# Show help
help:
	@echo "Available commands:"
	@echo "  build         - Build for current platform"
	@echo "  test          - Run tests"
	@echo "  clean         - Clean build artifacts"
	@echo "  install       - Install locally"
	@echo "  cross-compile - Build for all platforms"
	@echo "  package       - Create distribution packages"
	@echo "  dev           - Run in development mode"
	@echo "  help          - Show this help"
