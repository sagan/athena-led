APP := athena-led
BUILD_DIR := build
GO_BUILD_FLAGS := -ldflags="-s -w" -trimpath

.PHONY: all linux-arm linux-arm64 clean

all: linux-arm linux-arm64

linux-arm:
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=arm go build $(GO_BUILD_FLAGS) -o $(BUILD_DIR)/$(APP)-linux-arm .

linux-arm64:
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=arm64 go build $(GO_BUILD_FLAGS) -o $(BUILD_DIR)/$(APP)-linux-arm64 .

clean:
	rm -rf $(BUILD_DIR)
