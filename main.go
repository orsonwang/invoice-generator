package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// Invoice 發票資料結構
type Invoice struct {
	InvoiceNumber string  `json:"invoiceNumber"`
	InvoiceDate   string  `json:"invoiceDate"`
	FormatCode    string  `json:"formatCode"`
	RandomCode    string  `json:"randomCode"`
	Seller        Company `json:"seller"`
	Buyer         Company `json:"buyer"`
	Items         []Item  `json:"items"`
	TaxType       string  `json:"taxType"`
	SalesAmount   int     `json:"salesAmount"`
	TaxAmount     int     `json:"taxAmount"`
	TotalAmount   int     `json:"totalAmount"`
}

// Company 公司資料
type Company struct {
	Name    string `json:"name"`
	TaxID   string `json:"taxId"`
	Address string `json:"address"`
}

// Item 品項資料
type Item struct {
	Name      string `json:"name"`
	Quantity  int    `json:"quantity"`
	UnitPrice int    `json:"unitPrice"`
	Remark    string `json:"remark"`
}

const maxItemRows = 16

// 數字轉中文大寫
func numberToChinese(num int) string {
	if num == 0 {
		return "零元整"
	}

	digits := []string{"零", "壹", "貳", "參", "肆", "伍", "陸", "柒", "捌", "玖"}
	units := []string{"", "拾", "佰", "仟"}
	bigUnits := []string{"", "萬", "億"}

	numStr := fmt.Sprintf("%d", num)
	length := len(numStr)
	result := ""
	zeroFlag := false

	for i, char := range numStr {
		digit := int(char - '0')
		pos := length - i - 1
		unitPos := pos % 4
		bigUnitPos := pos / 4

		if digit == 0 {
			zeroFlag = true
			if unitPos == 0 && bigUnitPos > 0 {
				result += bigUnits[bigUnitPos]
			}
		} else {
			if zeroFlag {
				result += "零"
				zeroFlag = false
			}
			result += digits[digit] + units[unitPos]
			if unitPos == 0 && bigUnitPos > 0 {
				result += bigUnits[bigUnitPos]
			}
		}
	}

	return result + "元整"
}

// 格式化數字（千分位）
func formatNumber(num int) string {
	s := fmt.Sprintf("%d", num)
	n := len(s)
	if n <= 3 {
		return s
	}

	var result strings.Builder
	for i, c := range s {
		if i > 0 && (n-i)%3 == 0 {
			result.WriteByte(',')
		}
		result.WriteRune(c)
	}
	return result.String()
}

// TemplateData 模板資料
type TemplateData struct {
	Invoice
	ItemsTableRows       template.HTML
	SalesAmountFormatted string
	TaxAmountFormatted   string
	TotalAmountFormatted string
	TotalAmountChinese   string
	TaxableChecked       string
	ZeroTaxChecked       string
	TaxFreeChecked       string
}

// 產生 HTML
func generateHTML(invoice Invoice, templatePath string) (string, error) {
	// 產生品項表格行
	var itemsRows strings.Builder
	for _, item := range invoice.Items {
		amount := item.Quantity * item.UnitPrice
		itemsRows.WriteString(fmt.Sprintf(`
      <tr>
        <td style="width: 32%%; text-align: left; border-top:none;border-bottom:none;border-left:none;">%s</td>
        <td style="width: 8%%; text-align: center; border-top:none;border-bottom:none;border-left:none;">%d</td>
        <td style="width: 15%%; text-align: right; border-top:none;border-bottom:none;border-left:none;">%s</td>
        <td style="width: 15%%; text-align: right; border-top:none;border-bottom:none;border-left:none;">%s</td>
        <td style="width: 30%%; text-align: left; border-top:none;border-bottom:none;border-left:none;border-right:none;">%s</td>
      </tr>
    `, item.Name, item.Quantity, formatNumber(item.UnitPrice), formatNumber(amount), item.Remark))
	}

	// 補滿空白列
	emptyRowsNeeded := maxItemRows - len(invoice.Items)
	for i := 0; i < emptyRowsNeeded; i++ {
		itemsRows.WriteString(`
      <tr>
        <td style="width: 32%; text-align: left; border-top:none;border-bottom:none;border-left:none;">&nbsp;</td>
        <td style="width: 8%; text-align: center; border-top:none;border-bottom:none;border-left:none;">&nbsp;</td>
        <td style="width: 15%; text-align: right; border-top:none;border-bottom:none;border-left:none;">&nbsp;</td>
        <td style="width: 15%; text-align: right; border-top:none;border-bottom:none;border-left:none;">&nbsp;</td>
        <td style="width: 30%; text-align: left; border-top:none;border-bottom:none;border-left:none;border-right:none;">&nbsp;</td>
      </tr>
    `)
	}

	// 準備模板資料
	data := TemplateData{
		Invoice:              invoice,
		ItemsTableRows:       template.HTML(itemsRows.String()),
		SalesAmountFormatted: formatNumber(invoice.SalesAmount),
		TaxAmountFormatted:   formatNumber(invoice.TaxAmount),
		TotalAmountFormatted: formatNumber(invoice.TotalAmount),
		TotalAmountChinese:   numberToChinese(invoice.TotalAmount),
	}

	// 課稅類別
	switch invoice.TaxType {
	case "taxable":
		data.TaxableChecked = "checked"
	case "zeroTax":
		data.ZeroTaxChecked = "checked"
	case "taxFree":
		data.TaxFreeChecked = "checked"
	}

	// 讀取並解析模板
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("解析模板失敗: %w", err)
	}

	var htmlBuilder strings.Builder
	if err := tmpl.Execute(&htmlBuilder, data); err != nil {
		return "", fmt.Errorf("執行模板失敗: %w", err)
	}

	return htmlBuilder.String(), nil
}

// 產生 PDF
func generatePDF(htmlContent, outputPath string) error {
	// 建立 chromedp 上下文
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// 設定超時
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var buf []byte
	if err := chromedp.Run(ctx,
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}
			return page.SetDocumentContent(frameTree.Frame.ID, htmlContent).Do(ctx)
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			buf, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				WithPaperWidth(8.27).   // A4 width in inches
				WithPaperHeight(11.69). // A4 height in inches
				WithMarginTop(0).
				WithMarginBottom(0).
				WithMarginLeft(0).
				WithMarginRight(0).
				Do(ctx)
			return err
		}),
	); err != nil {
		return fmt.Errorf("產生 PDF 失敗: %w", err)
	}

	if err := os.WriteFile(outputPath, buf, 0644); err != nil {
		return fmt.Errorf("寫入 PDF 檔案失敗: %w", err)
	}

	return nil
}

func main() {
	var (
		inputPath  string
		outputPath string
		htmlOutput bool
	)

	flag.StringVar(&inputPath, "i", "", "JSON 輸入檔案路徑（必填）")
	flag.StringVar(&outputPath, "o", "invoice.pdf", "PDF 輸出檔案路徑")
	flag.BoolVar(&htmlOutput, "html", false, "同時輸出 HTML 檔案")
	flag.Parse()

	// 檢查必要參數
	if inputPath == "" {
		if flag.NArg() > 0 {
			inputPath = flag.Arg(0)
		} else {
			fmt.Println("錯誤：必須指定 JSON 輸入檔案")
			fmt.Println("使用方式: invoice-generator -i <輸入檔案> [-o <輸出檔案>] [--html]")
			os.Exit(1)
		}
	}

	// 讀取 JSON 檔案
	jsonData, err := os.ReadFile(inputPath)
	if err != nil {
		log.Fatalf("錯誤：無法讀取檔案 %s: %v", inputPath, err)
	}

	// 解析 JSON
	var invoice Invoice
	if err := json.Unmarshal(jsonData, &invoice); err != nil {
		log.Fatalf("錯誤：JSON 格式錯誤: %v", err)
	}

	// 取得模板路徑
	templatePath := "template_go.html"

	// 產生 HTML
	htmlContent, err := generateHTML(invoice, templatePath)
	if err != nil {
		log.Fatalf("錯誤：%v", err)
	}

	// 輸出 HTML（如果有指定）
	if htmlOutput {
		htmlPath := strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + ".html"
		if err := os.WriteFile(htmlPath, []byte(htmlContent), 0644); err != nil {
			log.Fatalf("錯誤：無法寫入 HTML 檔案: %v", err)
		}
		fmt.Printf("HTML 已輸出：%s\n", htmlPath)
	}

	// 產生 PDF
	fmt.Println("正在產生 PDF...")
	absOutputPath, _ := filepath.Abs(outputPath)
	if err := generatePDF(htmlContent, outputPath); err != nil {
		log.Fatalf("錯誤：%v", err)
	}
	fmt.Printf("PDF 已輸出：%s\n", absOutputPath)
}
