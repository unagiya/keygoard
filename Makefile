TARGET   := waveshare-rp2040-zero
EXAMPLE  := ./examples/zero-kb02/

.PHONY: test build flash fmt lint

# テスト
# 1. TinyGo: RP2040 ターゲットでのコンパイル検証（build 依存）
# 2. 標準 Go: machine 非依存パッケージのユニットテスト
test: build
	go test ./keycode/... ./internal/...

# TinyGo でのビルド確認
build:
	tinygo build -target=$(TARGET) $(EXAMPLE)

# RP2040 への書き込み
flash:
	tinygo flash -target=$(TARGET) $(EXAMPLE)

# コードフォーマット
fmt:
	goimports -w .

# Lint
lint:
	golangci-lint run
