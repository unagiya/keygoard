TARGET   := waveshare-rp2040-zero
EXAMPLE  := ./examples/zero-kb02/

.PHONY: test test-native build flash fmt lint

# machine 非依存パッケージのテスト（標準 Go）
test:
	go test ./keycode/...

# TinyGo native ターゲットでのテスト（ホスト上で machine を含むパッケージも実行可）
test-native:
	tinygo test -target=native ./...

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
