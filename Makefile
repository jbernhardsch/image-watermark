BINARY_NAME := watermark
DIST_DIR    := dist

.PHONY: all build clean darwin darwin-amd64 darwin-arm64 linux linux-amd64 linux-arm64 windows windows-amd64

all: darwin linux windows

# native build for the current machine
build:
	go build -o $(BINARY_NAME) .

build-darwin: darwin-amd64 darwin-arm64

darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 .

darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 .

build-linux: linux-amd64 linux-arm64

linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 .

linux-arm64:
	GOOS=linux GOARCH=arm64 go build -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 .

build-windows: windows-amd64

windows-amd64:
	GOOS=windows GOARCH=amd64 go build -o $(DIST_DIR)/$(BINARY_NAME)-windows.exe .

clean:
	rm -rf $(DIST_DIR) $(BINARY_NAME)
