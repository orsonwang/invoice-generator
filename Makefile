.PHONY: all build clean release linux windows darwin help

# 版本資訊
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)

# 輸出目錄
BUILD_DIR := build
RELEASE_DIR := release

# 應用程式名稱
APP_NAME := invoice-generator

# 必要附加檔案
REQUIRED_FILES := template_go.html sample-invoice.json README.md LICENSE

all: help

help:
	@echo "可用的 Make 目標："
	@echo "  build          - 編譯當前平台的執行檔"
	@echo "  linux          - 編譯 Linux (amd64) 版本"
	@echo "  windows        - 編譯 Windows (amd64) 版本"
	@echo "  darwin         - 編譯 macOS (amd64, arm64) 版本"
	@echo "  release        - 編譯所有平台並打包 release"
	@echo "  clean          - 清除編譯產物"

build:
	@echo "正在編譯 $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME) main.go
	@echo "✓ 編譯完成: $(BUILD_DIR)/$(APP_NAME)"

linux:
	@echo "正在編譯 Linux amd64 版本..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 main.go
	@echo "✓ Linux 版本編譯完成"

windows:
	@echo "正在編譯 Windows amd64 版本..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe main.go
	@echo "✓ Windows 版本編譯完成"

darwin:
	@echo "正在編譯 macOS amd64 版本..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 main.go
	@echo "✓ macOS amd64 版本編譯完成"
	@echo "正在編譯 macOS arm64 版本..."
	GOOS=darwin GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 main.go
	@echo "✓ macOS arm64 版本編譯完成"

release: clean linux windows darwin
	@echo "正在打包 release..."
	@mkdir -p $(RELEASE_DIR)

	# Linux amd64
	@echo "打包 Linux amd64..."
	@mkdir -p $(RELEASE_DIR)/$(APP_NAME)-$(VERSION)-linux-amd64
	@cp $(BUILD_DIR)/$(APP_NAME)-linux-amd64 $(RELEASE_DIR)/$(APP_NAME)-$(VERSION)-linux-amd64/$(APP_NAME)
	@cp $(REQUIRED_FILES) $(RELEASE_DIR)/$(APP_NAME)-$(VERSION)-linux-amd64/
	@cd $(RELEASE_DIR) && tar -czf $(APP_NAME)-$(VERSION)-linux-amd64.tar.gz $(APP_NAME)-$(VERSION)-linux-amd64
	@rm -rf $(RELEASE_DIR)/$(APP_NAME)-$(VERSION)-linux-amd64
	@echo "✓ Linux amd64 打包完成"

	# Windows amd64
	@echo "打包 Windows amd64..."
	@mkdir -p $(RELEASE_DIR)/$(APP_NAME)-$(VERSION)-windows-amd64
	@cp $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe $(RELEASE_DIR)/$(APP_NAME)-$(VERSION)-windows-amd64/$(APP_NAME).exe
	@cp $(REQUIRED_FILES) $(RELEASE_DIR)/$(APP_NAME)-$(VERSION)-windows-amd64/
	@cd $(RELEASE_DIR) && zip -q -r $(APP_NAME)-$(VERSION)-windows-amd64.zip $(APP_NAME)-$(VERSION)-windows-amd64
	@rm -rf $(RELEASE_DIR)/$(APP_NAME)-$(VERSION)-windows-amd64
	@echo "✓ Windows amd64 打包完成"

	# macOS amd64
	@echo "打包 macOS amd64..."
	@mkdir -p $(RELEASE_DIR)/$(APP_NAME)-$(VERSION)-darwin-amd64
	@cp $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 $(RELEASE_DIR)/$(APP_NAME)-$(VERSION)-darwin-amd64/$(APP_NAME)
	@cp $(REQUIRED_FILES) $(RELEASE_DIR)/$(APP_NAME)-$(VERSION)-darwin-amd64/
	@cd $(RELEASE_DIR) && tar -czf $(APP_NAME)-$(VERSION)-darwin-amd64.tar.gz $(APP_NAME)-$(VERSION)-darwin-amd64
	@rm -rf $(RELEASE_DIR)/$(APP_NAME)-$(VERSION)-darwin-amd64
	@echo "✓ macOS amd64 打包完成"

	# macOS arm64
	@echo "打包 macOS arm64..."
	@mkdir -p $(RELEASE_DIR)/$(APP_NAME)-$(VERSION)-darwin-arm64
	@cp $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 $(RELEASE_DIR)/$(APP_NAME)-$(VERSION)-darwin-arm64/$(APP_NAME)
	@cp $(REQUIRED_FILES) $(RELEASE_DIR)/$(APP_NAME)-$(VERSION)-darwin-arm64/
	@cd $(RELEASE_DIR) && tar -czf $(APP_NAME)-$(VERSION)-darwin-arm64.tar.gz $(APP_NAME)-$(VERSION)-darwin-arm64
	@rm -rf $(RELEASE_DIR)/$(APP_NAME)-$(VERSION)-darwin-arm64
	@echo "✓ macOS arm64 打包完成"

	@echo ""
	@echo "========================================="
	@echo "Release 打包完成！"
	@echo "版本: $(VERSION)"
	@echo "輸出目錄: $(RELEASE_DIR)/"
	@echo ""
	@ls -lh $(RELEASE_DIR)/*.{tar.gz,zip} 2>/dev/null || true
	@echo "========================================="

clean:
	@echo "清除編譯產物..."
	@rm -rf $(BUILD_DIR) $(RELEASE_DIR)
	@echo "✓ 清除完成"
