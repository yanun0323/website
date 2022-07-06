# 在 VSCode 設置 Swift 開發環境

紀錄如何在 VSCode 上設置 Swift 開發所需的環境

[link](#apple-swift-format)

## 安裝 VSCode Extension

- [**Swift**](https://marketplace.visualstudio.com/items?itemName=sswg.swift-lang) - apple 開發者維護的 extension
- [**apple-swift-format**](https://marketplace.visualstudio.com/items?itemName=vknabel.vscode-apple-swift-format) - swift 自動 format 功能

- (選擇) [**Android iOS Emulator**](https://marketplace.visualstudio.com/items?itemName=DiemasMichiels.emulate) - 運行模擬器用

## apple-swift-format 安裝流程

1. 安裝 Xcode 並開啟安裝 Xcode-command-line
2. 安裝 mint

```
brew install mint
```

3. 指定 xcode-selector 位置

```
sudo xcode-select -s /Applications/Xcode.app/Contents/Developer
```

4. 安裝 apple 官方 format

```go
mint install apple/swift-format@release/{VERSION}
/* {VERSION}= 5.6 or others */
```

5. 找到 apple/swift-format 的存放路徑

```
mint which swift-format
```

6. 到 VSCode preference 找到 'Apple-swift-format: Path'，將 apple/swift-format 的路徑貼上去
7. 完成！在任何 .swift 檔案案儲存就會自動 format 了
