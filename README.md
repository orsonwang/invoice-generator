# 電子發票產生器 (Taiwan E-Invoice Generator)

因為使用中某平台的電子發票使用太多的需要另外下載的中國字型導致沒安裝字型的話產出的發票不忍猝睹，提了幾次功能建議都講送給研發團隊考慮，幾年都沒下文，所以就自己寫一下。因為是自用所以只有A4格式二。

台灣電子發票證明聯產生工具，依照財政部電子發票實施作業要點「格式二」規範，將 JSON 格式的發票資料轉換為標準 A4 格式的 PDF 電子發票。

## 版本說明

本專案提供兩種實作版本：

### Go 版本（推薦，golang-version 分支）
- ✅ **編譯成單一執行檔**，無需安裝 Node.js 或其他依賴
- ✅ **效能更好**，啟動快速、記憶體占用低
- ✅ **跨平台**，支援 Windows、macOS、Linux
- ✅ **使用 chromedp** 直接操作 Chrome DevTools Protocol

### Node.js 版本（master 分支）
- 使用 Puppeteer + Node.js
- 適合已有 Node.js 環境的專案

**建議使用 Go 版本**，以下說明以 Go 版本為主。

## 功能特色

- 符合財政部電子發票「格式二」規範
- 固定 A4 尺寸排版 (21cm × 29.7cm)
- 支援應稅、零稅率、免稅等課稅類別
- 自動計算金額並轉換為中文大寫
- 品項區自動填充空白列至固定行數
- 命令列介面 (CLI)，方便整合至自動化流程

## 安裝（Go 版本）

### 系統需求

- Go 1.19 或以上版本
- Chrome 或 Chromium 瀏覽器（chromedp 會自動下載）

### 安裝步驟

```bash
# 複製本專案（Go 版本分支）
git clone -b golang-version https://github.com/orsonwang/invoice-generator.git
cd invoice-generator

# 編譯成執行檔
go build -o invoice-generator main.go

# 或直接執行（不編譯）
go run main.go -i sample-invoice.json -o invoice.pdf
```

### 下載預編譯執行檔

（未來會提供 Releases）

## 使用方式

### 基本用法

```bash
# 使用編譯好的執行檔
./invoice-generator -i <輸入JSON檔案> -o <輸出PDF檔案>

# 或使用 go run
go run main.go -i <輸入JSON檔案> -o <輸出PDF檔案>
```

### 範例

```bash
# 使用範例資料產生發票
./invoice-generator -i sample-invoice.json -o invoice.pdf

# 同時輸出 HTML 檔案（便於除錯）
./invoice-generator -i sample-invoice.json -o invoice.pdf --html
```

### 命令列選項

- `-i <path>` - JSON 輸入檔案路徑（必填，或直接作為第一個參數）
- `-o <path>` - PDF 輸出檔案路徑（預設：`invoice.pdf`）
- `--html` - 同時輸出 HTML 檔案

## JSON 資料格式

```json
{
  "invoiceNumber": "TK45937999",
  "invoiceDate": "2027-10-13",
  "formatCode": "25",
  "randomCode": "",
  "seller": {
    "name": "雲杉科技股份有限公司",
    "taxId": "52449873",
    "address": "台北市信義區基隆路一段172巷1號3樓"
  },
  "buyer": {
    "name": "XXXX股份有限公司",
    "taxId": "92760123",
    "address": "臺北市中正區重慶南路１７號６樓"
  },
  "items": [
    {
      "name": "韌體安全保護系統",
      "quantity": 1,
      "unitPrice": 7188601,
      "remark": "2025Q3"
    }
  ],
  "taxType": "taxable",
  "salesAmount": 7188601,
  "taxAmount": 35942,
  "totalAmount": 7548023
}
```

### 欄位說明

#### 基本資訊
- `invoiceNumber` - 發票號碼
- `invoiceDate` - 開立日期
- `formatCode` - 格式代號
- `randomCode` - 隨機碼（B2B不用填）

#### 賣方資訊 (seller)
- `name` - 公司名稱
- `taxId` - 統一編號
- `address` - 地址

#### 買方資訊 (buyer)
- `name` - 公司名稱
- `taxId` - 統一編號
- `address` - 地址

#### 品項清單 (items)
- `name` - 品名
- `quantity` - 數量
- `unitPrice` - 單價
- `remark` - 備註（選填）

#### 金額資訊
- `taxType` - 課稅類別：`"taxable"`（應稅）、`"zeroTax"`（零稅率）、`"taxFree"`（免稅）
- `salesAmount` - 銷售額合計
- `taxAmount` - 營業稅
- `totalAmount` - 總計

## 技術說明

### Go 版本技術棧

- **Go** - 執行環境
- **chromedp** - Chrome DevTools Protocol 客戶端
- **html/template** - HTML 模板引擎

### Node.js 版本技術棧（master 分支）

- **Node.js** - 執行環境
- **Puppeteer** - HTML 轉 PDF
- **Commander** - 命令列介面

### 設計特點

- 因為是用網頁template，可以自己任選喜歡的字型
- 品項區固定 16 行，資料不足時自動補空白列，沒有作規格中第2頁的功能
- 使用巢狀表格處理品項列表，確保品名換行時其他欄位同步增高，品項文字兩行內不會有問題
- 沒有任何資料正確性檢查
- 金額自動格式化（千分位逗號）
- 總計自動轉換為中文大寫金額
- 需要總備註可以善用沒有資料檢查的特性，直接寫只有備註的空行，一樣最多兩行，超過兩行放到下一個品項會比較好

## 專案結構

### Go 版本（golang-version 分支）
```
.
├── main.go               # 主程式
├── template_go.html      # 發票 HTML 模板
├── sample-invoice.json   # 範例資料
├── go.mod                # Go 模組定義
├── go.sum                # Go 依賴鎖定檔
├── LICENSE               # MIT 授權條款
└── README.md            # 本說明文件
```

### Node.js 版本（master 分支）
```
.
├── index.js              # 主程式
├── template.html         # 發票 HTML 模板
├── sample-invoice.json   # 範例資料
├── package.json          # 套件設定
├── LICENSE               # MIT 授權條款
└── README.md            # 本說明文件
```

## 效能比較

| 項目 | Go 版本 | Node.js 版本 |
|------|---------|-------------|
| 啟動時間 | ~0.5s | ~2s |
| 記憶體占用 | ~80MB | ~150MB |
| 執行檔大小 | ~20MB（單一檔案） | 需要 node_modules（~200MB） |
| 安裝依賴 | 無需安裝 | 需要 npm/pnpm install |

## 授權

MIT License - 詳見 [LICENSE](LICENSE) 檔案

## 相關資源

- [財政部電子發票實施作業要點](https://www.etax.nat.gov.tw/)
- [chromedp 文件](https://github.com/chromedp/chromedp)
- [Puppeteer 文件](https://pptr.dev/) (Node.js 版本)

## 問題回報

如有任何問題或建議，歡迎在 [GitHub Issues](https://github.com/orsonwang/invoice-generator/issues) 提出。
