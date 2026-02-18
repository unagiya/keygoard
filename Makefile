.PHONY: help install build build-pico build-pico2 build-example build-example-pico build-example-pico2 test fmt lint clean deps setup

# デフォルトターゲット
help:
	@echo "keygoard - Makefile commands"
	@echo ""
	@echo "Usage:"
	@echo "  make install               Install keygoard CLI tool (standard Go)"
	@echo "  make build                 Build all packages for both targets (TinyGo)"
	@echo "  make build-pico            Build all packages for RP2040 (TinyGo)"
	@echo "  make build-pico2           Build all packages for RP2350 (TinyGo)"
	@echo "  make build-example         Build all examples for both targets (TinyGo)"
	@echo "  make build-example-pico    Build all examples for RP2040 (TinyGo)"
	@echo "  make build-example-pico2   Build all examples for RP2350 (TinyGo)"
	@echo "  make test                  Run tests (TinyGo)"
	@echo "  make fmt                   Format code"
	@echo "  make lint                  Run linter"
	@echo "  make clean                 Clean build artifacts"
	@echo "  make deps                  Update dependencies"
	@echo ""

# CLIツールのインストール（標準Goを使用、machineパッケージに依存しないため）
install:
	@echo "Installing keygoard CLI..."
	go install ./cmd/keygoard

# 全パッケージのビルド確認（両ターゲット）
build: build-pico build-pico2

# 全パッケージのビルド確認（RP2040向け）
build-pico:
	@echo "Building all packages with TinyGo (target: pico)..."
	tinygo test -target=pico -run='^$$' ./...

# 全パッケージのビルド確認（RP2350向け）
build-pico2:
	@echo "Building all packages with TinyGo (target: pico2)..."
	tinygo test -target=pico2 -run='^$$' ./...

# サンプルのビルド（両ターゲット）
build-example: build-example-pico build-example-pico2

# サンプルのビルド（RP2040向け、各exampleは独自のgo.modを持つ）
build-example-pico:
	@echo "Building examples for RP2040 (target: pico)..."
	@for dir in examples/*/; do \
		if [ -f "$$dir/main.go" ]; then \
			name=$$(basename $$dir); \
			echo "  Building $$name..."; \
			(cd $$dir && tinygo build -target=pico -o firmware.uf2 .) || exit 1; \
		elif ls $$dir/*/main.go >/dev/null 2>&1; then \
			for subdir in $$dir/*/; do \
				if [ -f "$$subdir/main.go" ]; then \
					name=$$(basename $$dir)/$$(basename $$subdir); \
					echo "  Building $$name..."; \
					(cd $$subdir && tinygo build -target=pico -o firmware.uf2 .) || exit 1; \
				fi; \
			done; \
		fi; \
	done
	@echo "Done. Firmware files (pico) created in each example directory."

# サンプルのビルド（RP2350向け、各exampleは独自のgo.modを持つ）
build-example-pico2:
	@echo "Building examples for RP2350 (target: pico2)..."
	@for dir in examples/*/; do \
		if [ -f "$$dir/main.go" ]; then \
			name=$$(basename $$dir); \
			echo "  Building $$name..."; \
			(cd $$dir && tinygo build -target=pico2 -o firmware-pico2.uf2 .) || exit 1; \
		elif ls $$dir/*/main.go >/dev/null 2>&1; then \
			for subdir in $$dir/*/; do \
				if [ -f "$$subdir/main.go" ]; then \
					name=$$(basename $$dir)/$$(basename $$subdir); \
					echo "  Building $$name..."; \
					(cd $$subdir && tinygo build -target=pico2 -o firmware-pico2.uf2 .) || exit 1; \
				fi; \
			done; \
		fi; \
	done
	@echo "Done. Firmware files (pico2) created in each example directory."

# テスト実行（TinyGoを使用）
test:
	@echo "Running tests with TinyGo..."
	tinygo test ./...

# コードフォーマット（標準Goを使用、フォーマットはTinyGo非依存）
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint実行（golangci-lintが必要）
# 注意: machineパッケージに依存するパッケージはgolangci-lintでエラーになる場合がある
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint is not installed. Skipping..."; \
		echo "Install: https://golangci-lint.run/usage/install/"; \
	fi

# 依存関係の更新
deps:
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

# クリーンアップ
clean:
	@echo "Cleaning build artifacts..."
	@find examples -name "firmware.uf2" -delete 2>/dev/null || true
	@find examples -name "firmware-pico2.uf2" -delete 2>/dev/null || true
	rm -rf build/
	go clean

# 開発環境のセットアップ
setup:
	@echo "Setting up development environment..."
	@echo "1. Checking Go installation..."
	@go version
	@echo "2. Checking TinyGo installation..."
	@tinygo version || echo "TinyGo not found. Please install: https://tinygo.org/getting-started/install/"
	@echo "3. Installing dependencies..."
	@go mod download
	@echo "4. Installing keygoard CLI..."
	@go install ./cmd/keygoard
	@echo ""
	@echo "Setup complete! Try: keygoard version"
