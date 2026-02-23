# keygoard

TinyGo ベースのカスタムキーボードファームウェアフレームワーク。

- ターゲット MCU: RP2040（Waveshare RP2040-Zero）
- TinyGo ターゲット名: `waveshare-rp2040-zero`
- 開発方針: **最小構成から実機動作確認しながら段階的に積み上げる**

---

## テスト実機：zero-kb02（sago35/keyboards）

- リポジトリ: https://github.com/sago35/keyboards/tree/main/zero-kb02
- 別ファームウェアで動作確認済みのため、開発中のエラーはハード・マイコン故障ではなくファームウェアの問題とみなしてよい

### ハードウェア仕様（ピン配置）

| 機能 | GPIO |
|------|------|
| Matrix Col 0–3（出力） | GP5, GP6, GP7, GP8 |
| Matrix Row 0–2（入力） | GP9, GP10, GP11 |
| Matrix 方式 | COL2ROW（列をHighにして行をRead） |
| RGB LED（SK6812MINI-E） | GP1、12個 |
| Encoder A / B / Button | GP3 / GP4 / GP0 |
| OLED SDA / SCL | GP12 / GP13（I2C、SSD1306 128×64、addr=0x3C） |
| Joystick X / Y / Button | GP29(ADC) / GP28(ADC) / GP0 |

### テスト環境

- 主: macOS（USB Prober / System Information で認識確認）
- 副: 将来的に複数 OS で確認

---

## セットアップ

### 必要なツール

- [TinyGo](https://tinygo.org/) 0.40.1
- [Go](https://golang.org/) 1.25
- `goimports`（フォーマット用）
- `golangci-lint`（Lint 用）

### ビルド・書き込み

```sh
# ビルド確認
make build

# 実機への書き込み
make flash

# テスト（machine 非依存パッケージ）
go test ./keycode/... ./matrix/...

# フォーマット
make fmt

# Lint
make lint
```

---

## 開発フェーズ

詳細は [ROADMAP.md](ROADMAP.md) を参照。

| Phase | 内容 | 状態 |
|---|---|---|
| Phase 1 | 最小 HID キーボード（マトリクス・デバウンス・HID） | 進行中 |
| Phase 2 | レイヤーシステム（MO / TG / TT / LT） | 未着手 |
| Phase 3 | 周辺機器（エンコーダー・LED・OLED） | 未着手 |
| Phase 4 | Joystick + Gamepad | 未着手 |
| Phase 5 | Split 対応 | 未着手 |
| Phase 6 | 高度な機能（マクロ・Flash 永続化・NKRO） | 未着手 |
