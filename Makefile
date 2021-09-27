SHELL := /bin/bash

# Go and compilation related variables
BUILD_DIR ?= out
BINARY_NAME ?= proxy

# Add default target
.PHONY: default
default: install

.PHONY: install
install: $(SOURCES)
	go build -o $(BINARY_NAME)

$(BUILD_DIR)/macos-amd64/$(BINARY_NAME):
	GOARCH=amd64 GOOS=darwin go build -o $(BUILD_DIR)/macos-amd64/$(BINARY_NAME)

$(BUILD_DIR)/linux-amd64/$(BINARY_NAME):
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/linux-amd64/$(BINARY_NAME)

$(BUILD_DIR)/windows-amd64/$(BINARY_NAME).exe:
	GOARCH=amd64 GOOS=windows go build -o $(BUILD_DIR)/windows-amd64/$(BINARY_NAME).exe

.PHONY: cross ## Cross compiles all binaries
cross: $(BUILD_DIR)/macos-amd64/$(BINARY_NAME) $(BUILD_DIR)/linux-amd64/$(BINARY_NAME) $(BUILD_DIR)/windows-amd64/$(BINARY_NAME).exe

.PHONY: clean
clean:
	rm -fr $(BUILD_DIR)
	rm -fr $(BINARY_NAME)