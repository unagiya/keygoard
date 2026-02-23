TARGET   := waveshare-rp2040-zero
EXAMPLE  := ./examples/zero-kb02/

.PHONY: test build flash fmt lint

# テスト
# 1. 標準 Go: machine 非依存パッケージ（keycode, matrix/debounce）
# 2. TinyGo:  RP2040 ターゲットで全パッケージのビルド検証（実機接続が必要）
test:
	go test ./keycode/... ./matrix/...
	tinygo test -target=$(TARGET) ./...

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
