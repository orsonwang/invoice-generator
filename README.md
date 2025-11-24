# 電子發票產生器 (Taiwan E-Invoice Generator)
因為使用中某平台的電子發票使用太多的需要另外下載的中國字型導致沒安裝字型的話產出的發票不忍猝睹，提了幾次功能建議都講送給研發團隊考慮，幾年都沒下文，所以就自己寫一下。因為是自用所以只有A4格式二。

台灣電子發票證明聯產生工具，依照財政部電子發票實施作業要點「格式二」規範，將 JSON 格式的發票資料轉換為標準 A4 格式的 PDF 電子發票。

## 功能特色

- 符合財政部電子發票「格式二」規範
- 固定 A4 尺寸排版 (21cm × 29.7cm)
- 支援應稅、零稅率、免稅等課稅類別
- 自動計算金額並轉換為中文大寫
- 品項區自動填充空白列至固定行數
- 命令列介面 (CLI)，方便整合至自動化流程

## 安裝

### 系統需求

- Node.js 18 或以上版本
- pnpm (推薦) 或 npm

### 安裝步驟

```bash
# 複製本專案
git clone https://github.com/orsonwang/invoice-generator.git
cd invoice-generator

# 安裝專案使用的套件
pnpm install

# 安裝 Chrome 瀏覽器（Puppeteer 需要）
npx puppeteer browsers install chrome
```

## 使用方式

### 基本用法

```bash
node index.js <輸入JSON檔案> -o <輸出PDF檔案>
```

### 範例

```bash
# 使用範例資料產生發票
node index.js sample-invoice.json -o invoice.pdf

# 同時輸出 HTML 檔案（便於除錯）
node index.js sample-invoice.json -o invoice.pdf --html
```

### 命令列選項

- `<input>` - JSON 輸入檔案路徑（必填）
- `-o, --output <path>` - PDF 輸出檔案路徑（預設：`invoice.pdf`）
- `--html` - 同時輸出 HTML 檔案
- `-V, --version` - 顯示版本資訊
- `-h, --help` - 顯示說明

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

### 主要技術棧

- **Node.js** - 執行環境
- **Puppeteer** - HTML 轉 PDF（使用 Chrome 無頭瀏覽器）
- **Commander** - 命令列介面
- **ES Modules** - 使用現代 JavaScript 模組系統

### 設計特點

- 因為是用網頁template，可以自己任選喜歡的字型
- 品項區固定 16 行，資料不足時自動補空白列，沒有作規格中第2頁的功能
- 使用巢狀表格處理品項列表，確保品名換行時其他欄位同步增高，品項文字兩行內不會有問題
- 沒有任何資料正確性檢查
- 金額自動格式化（千分位逗號）
- 總計自動轉換為中文大寫金額
- 需要總備註可以善用沒有資料檢查的特性，直接寫只有備註的空行，一樣最多兩行，超過兩行放到下一個品項會比較好。

## 專案結構

```
.
├── index.js              # 主程式（CLI 入口）
├── template.html         # 發票 HTML 模板
├── sample-invoice.json   # 範例資料
├── package.json          # 套件設定
├── LICENSE               # MIT 授權條款
└── README.md            # 本說明文件
```

## 授權

MIT License - 詳見 [LICENSE](LICENSE) 檔案

## 相關資源

- [財政部電子發票實施作業要點](https://www.etax.nat.gov.tw/)
- [Puppeteer 文件](https://pptr.dev/)

## 問題回報

如有任何問題或建議，歡迎在 [GitHub Issues](https://github.com/orsonwang/invoice-generator/issues) 提出。
