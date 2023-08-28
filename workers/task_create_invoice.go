package workers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/phpdave11/gofpdf"
	"github.com/qwetu_petro/backend/utils"
	"github.com/rs/zerolog/log"
	"strings"
)

const (
	PettyCashDir      = "generated_files/petty-cash"
	PaymentRequestDir = "generated_files/payment-requests"
	InvoicesDir       = "generated_files/invoices"
	QuotationsDir     = "generated_files/quotations"
	PurchaseOrderDir  = "generated_files/purchase-orders"
)

const (
	TaskTypeCreateInvoicePdf        = "task:create_invoice_pdf"
	TaskTypeCreateQuotationPdf      = "task:create_quotation_pdf"
	TaskTypeCreatePurchaseOrderPdf  = "task:create_purchase_order_pdf"
	TaskTypeCreatePaymentRequestPdf = "task:create_payment_request_pdf"
	TaskTypeCreatePettyCashPdf      = "task:create_petty_cash_pdf"
	TaskSendEmail                   = "task:send_email"
	TaskSendSMS                     = "task:send_sms"
)

type CreateInvoicePdfPayload struct {
	InvoiceID int32 `json:"invoice_id"`
}

func (distributor *RedisTaskDistributor) DistributeTaskCreateInvoicePdf(ctx context.Context, payload *CreateInvoicePdfPayload, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("could not marshal payload: %w", err)

	}

	task := asynq.NewTask(TaskTypeCreateInvoicePdf, jsonPayload, opts...)

	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("could not enqueue task: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueued task")

	return nil

}

func (processor *RedisTaskProcessor) ProcessTaskCreateInvoicePdf(ctx context.Context, task *asynq.Task) error {
	var payload CreateInvoicePdfPayload
	conf := processor.conf

	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("could not unmarshal payload: %w", asynq.SkipRetry)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Msg("processing task")

	invoice, err := processor.db.GetInvoiceById(ctx, payload.InvoiceID)
	if err != nil {
		return fmt.Errorf("could not get invoice by id: %w", err)
	}

	invoiceItems, err := processor.db.GetInvoiceItemsByInvoiceID(ctx, payload.InvoiceID)
	if err != nil {
		return fmt.Errorf("could not get invoice items by invoice id: %w", err)
	}

	bankInfo, err := processor.db.GetBankInfoByID(ctx, invoice.BankDetails)
	if err != nil {
		return fmt.Errorf("could not get bank info by id: %w", err)
	}

	signatory, err := processor.db.GetSignatoryById(ctx, int64(invoice.SignatoryID))
	if err != nil {
		return fmt.Errorf("could not get signatory by id: %w", err)
	}

	company, err := processor.db.GetCompanyByID(ctx, int64(invoice.CompanyID))
	if err != nil {
		return fmt.Errorf("could not get company by id: %w", err)
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	utils.PositionCompanyLogoAndDetailsAtTop(pdf)

	// Title
	pdf.SetFont("Helvetica", "B", 7) // Set a larger font size for the title
	title := "INVOICE"

	// Calculate width of the title to center it
	width, _ := pdf.GetPageSize()
	titleWidth := pdf.GetStringWidth(title) + 6 // 6 is additional padding
	offset := (width - titleWidth) / 2.0

	// Set X position to the calculated offset
	pdf.SetX(offset)

	// Print the title
	pdf.Cell(70, 10, title) // Increased height to 10 for better appearance with larger font size
	pdf.Ln(10)

	pdf.SetFont("Helvetica", "B", 7)
	pdf.Cell(30, 5, fmt.Sprintf("Attn:  %s", invoice.Attn))
	pdf.Ln(3)

	pdf.Cell(30, 5, fmt.Sprintf("Company:  %s", company.Name))
	pdf.Ln(3)

	pdf.Cell(30, 5, fmt.Sprintf("Site: %s", invoice.Site))
	pdf.Ln(3)

	pdf.Cell(30, 5, fmt.Sprintf("Date: %s", invoice.Date.Format("2006-01-02")))
	pdf.Ln(3)

	pdf.Cell(30, 5, fmt.Sprintf("Invoice Number: %s", invoice.InvoiceNo))
	pdf.Ln(3)

	pdf.Cell(30, 5, fmt.Sprintf("Purchase Order Number: %s", invoice.PurchaseOrderNumber))
	pdf.Ln(3)

	headers := []string{"Description", "UOM", "Qty", "Unit Price", "Net Price", "Currency"}
	headerWidths := []float64{80, 30, 15, 20, 20, 20} // these should add up to the width of your page

	// invoice items
	for i, header := range headers {
		pdf.SetFont("Helvetica", "", 6)
		pdf.CellFormat(headerWidths[i], 10, header, "0", 0, "L", false, 0, "")
	}
	pdf.Ln(-1)

	for _, item := range invoiceItems {
		pdf.SetFont("Helvetica", "", 6)

		// Replacing the plus signs with newline characters in the description
		descriptionWithNewLines := strings.Replace(item.Description, "+", "\n", -1)

		// Using MultiCell for the description so that the newline characters are effective
		pdf.MultiCell(headerWidths[0], 4, descriptionWithNewLines, "0", "L", false)

		// Getting current X and Y position after MultiCell
		currentX, currentY := pdf.GetX(), pdf.GetY()

		// Setting X and Y for the next cells
		pdf.SetXY(currentX+headerWidths[0], currentY-10)

		// Printing the rest of the cells
		pdf.CellFormat(headerWidths[1], 10, item.Uom, "0", 0, "L", false, 0, "")
		pdf.CellFormat(headerWidths[2], 10, fmt.Sprintf("%d", item.Qty), "0", 0, "L", false, 0, "")
		pdf.CellFormat(headerWidths[3], 10, fmt.Sprintf("%.2f", item.UnitPrice), "0", 0, "L", false, 0, "")
		pdf.CellFormat(headerWidths[4], 10, fmt.Sprintf("%.2f", float64(item.Qty)*item.UnitPrice), "0", 0, "L", false, 0, "")
		pdf.CellFormat(headerWidths[5], 10, item.Currency, "0", 0, "L", false, 0, "")
		pdf.Ln(-1)
	}

	utils.PositionSignatoryAndBankDetailsAtBottom(pdf, signatory, bankInfo)

	var pdfBuffer bytes.Buffer
	name := fmt.Sprintf("%s.pdf", invoice.InvoiceNo)

	err = pdf.Output(&pdfBuffer)
	if err != nil {
		return fmt.Errorf("could not create PDF: %w", err)
	}
	// Upload the PDF to S3.

	err = utils.UploadFileToS3Bucket(conf, name, InvoicesDir, pdfBuffer.Bytes())
	if err != nil {
		return fmt.Errorf("could not upload PDF to S3: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("Invoice", name).Msg("processed task")

	return nil
}
