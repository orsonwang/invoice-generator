#!/usr/bin/env node

import { program } from "commander";
import puppeteer from "puppeteer";
import fs from "fs";
import path from "path";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// 數字轉中文大寫
function numberToChinese(num) {
  const digits = ["零", "壹", "貳", "參", "肆", "伍", "陸", "柒", "捌", "玖"];
  const units = ["", "拾", "佰", "仟"];
  const bigUnits = ["", "萬", "億"];

  if (num === 0) return "零元整";

  const numStr = Math.floor(num).toString();
  const len = numStr.length;
  let result = "";
  let zeroFlag = false;

  for (let i = 0; i < len; i++) {
    const digit = parseInt(numStr[i]);
    const pos = len - i - 1;
    const unitPos = pos % 4;
    const bigUnitPos = Math.floor(pos / 4);

    if (digit === 0) {
      zeroFlag = true;
      if (unitPos === 0 && bigUnitPos > 0) {
        result += bigUnits[bigUnitPos];
      }
    } else {
      if (zeroFlag) {
        result += "零";
        zeroFlag = false;
      }
      result += digits[digit] + units[unitPos];
      if (unitPos === 0 && bigUnitPos > 0) {
        result += bigUnits[bigUnitPos];
      }
    }
  }

  return result + "元整";
}

// 格式化數字（千分位）
function formatNumber(num) {
  return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
}

// 產生 HTML
function generateHTML(data) {
  const templatePath = path.join(__dirname, "template.html");
  let html = fs.readFileSync(templatePath, "utf-8");

  // 產生明細表格行（每個商品一行，無分隔線）
  let itemsTableRows = "";
  for (const item of data.items) {
    const amount = item.quantity * item.unitPrice;
    itemsTableRows += `
      <tr>
        <td style="width: 32%; text-align: left; border-top:none;border-bottom:none;border-left:none;">${item.name}</td>
        <td style="width: 8%; text-align: center; border-top:none;border-bottom:none;border-left:none;">${item.quantity}</td>
        <td style="width: 15%; text-align: right; border-top:none;border-bottom:none;border-left:none;">${formatNumber(item.unitPrice)}</td>
        <td style="width: 15%; text-align: right; border-top:none;border-bottom:none;border-left:none;">${formatNumber(amount)}</td>
        <td style="width: 30%; text-align: left; border-top:none;border-bottom:none;border-left:none;border-right:none;">${item.remark || ""}</td>
      </tr>
    `;
  }

  // 課稅類別
  const taxableChecked = data.taxType === "taxable" ? "checked" : "";
  const zeroTaxChecked = data.taxType === "zeroTax" ? "checked" : "";
  const taxFreeChecked = data.taxType === "taxFree" ? "checked" : "";

  // 替換模板變數
  const replacements = {
    "{{sellerName}}": data.seller.name,
    "{{invoiceDate}}": data.invoiceDate,
    "{{currentPage}}": "1",
    "{{totalPages}}": "1",
    "{{invoiceNumber}}": data.invoiceNumber,
    "{{buyerName}}": data.buyer.name,
    "{{buyerTaxId}}": data.buyer.taxId,
    "{{buyerAddress}}": data.buyer.address,
    "{{formatCode}}": data.formatCode,
    "{{randomCode}}": data.randomCode || "",
    "{{itemsTableRows}}": itemsTableRows,
    "{{salesAmount}}": formatNumber(data.salesAmount),
    "{{taxAmount}}": formatNumber(data.taxAmount),
    "{{totalAmount}}": formatNumber(data.totalAmount),
    "{{totalAmountChinese}}": numberToChinese(data.totalAmount),
    "{{taxableChecked}}": taxableChecked,
    "{{zeroTaxChecked}}": zeroTaxChecked,
    "{{taxFreeChecked}}": taxFreeChecked,
    "{{sellerTaxId}}": data.seller.taxId,
    "{{sellerAddress}}": data.seller.address,
  };

  for (const [key, value] of Object.entries(replacements)) {
    html = html.replace(new RegExp(key, "g"), value);
  }

  return html;
}

// 產生 PDF
async function generatePDF(htmlContent, outputPath) {
  const browser = await puppeteer.launch({
    headless: "new",
    args: ["--no-sandbox", "--disable-setuid-sandbox"],
  });

  const page = await browser.newPage();
  await page.setContent(htmlContent, { waitUntil: "networkidle0" });

  await page.pdf({
    path: outputPath,
    format: "A4",
    printBackground: true,
    margin: { top: 0, right: 0, bottom: 0, left: 0 },
  });

  await browser.close();
}

// CLI 主程式
program
  .name("einvoice")
  .description("台灣電子發票證明聯產生器 (格式二)")
  .version("1.0.0")
  .argument("<input>", "JSON 輸入檔案路徑")
  .option("-o, --output <path>", "PDF 輸出檔案路徑", "invoice.pdf")
  .option("--html", "同時輸出 HTML 檔案")
  .action(async (input, options) => {
    try {
      // 讀取 JSON
      const jsonPath = path.resolve(input);
      if (!fs.existsSync(jsonPath)) {
        console.error(`錯誤：找不到輸入檔案 ${jsonPath}`);
        process.exit(1);
      }

      const data = JSON.parse(fs.readFileSync(jsonPath, "utf-8"));

      // 產生 HTML
      const html = generateHTML(data);

      // 輸出 HTML（如果有指定）
      if (options.html) {
        const htmlPath = options.output.replace(".pdf", ".html");
        fs.writeFileSync(htmlPath, html);
        console.log(`HTML 已輸出：${htmlPath}`);
      }

      // 產生 PDF
      const outputPath = path.resolve(options.output);
      console.log(`正在產生 PDF...`);
      await generatePDF(html, outputPath);
      console.log(`PDF 已輸出：${outputPath}`);
    } catch (error) {
      console.error(`錯誤：${error.message}`);
      process.exit(1);
    }
  });

program.parse();
