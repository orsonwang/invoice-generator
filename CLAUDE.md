# 電子發票證明聯產生器

## 專案說明
台灣電子發票證明聯 CLI 工具，依照財政部格式二規範產生 A4 PDF 發票。

## 技術架構
- Node.js (ES Module)
- Puppeteer - HTML 轉 PDF
- Commander - CLI 參數處理

## 使用方式
```bash
node index.js <輸入JSON> -o <輸出PDF> [--html]
```

範例：
```bash
node index.js sample-invoice.json -o invoice.pdf --html
```

## 檔案結構
- `index.js` - CLI 主程式
- `template.html` - 發票 HTML 模板
- `sample-invoice.json` - JSON 輸入範例

## JSON 輸入格式
```json
{
  "invoiceNumber": "TK45937062",
  "invoiceDate": "2025-10-13",
  "formatCode": "25",
  "randomCode": "",
  "buyer": {
    "name": "公司名稱",
    "taxId": "統一編號",
    "address": "地址"
  },
  "seller": {
    "name": "公司名稱",
    "taxId": "統一編號",
    "address": "地址"
  },
  "items": [
    {
      "name": "品名",
      "quantity": 1,
      "unitPrice": 1000,
      "remark": "備註"
    }
  ],
  "taxType": "taxable",
  "salesAmount": 1000,
  "taxAmount": 50,
  "totalAmount": 1050
}
```

## 開發指令
```bash
pnpm install          # 安裝相依套件
node index.js --help  # 查看說明
```
