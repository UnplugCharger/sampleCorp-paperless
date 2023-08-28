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

type CreatePurchaseOrderPayload struct {
	// Purchase order id
	PurchaseOrderID int32 `json:"purchase_order_id"`
}

func (distributor *RedisTaskDistributor) DistributeTaskCreatePurchaseOrderPdf(ctx context.Context, payload *CreatePurchaseOrderPayload, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("could not marshal payload: %w", err)
	}

	task := asynq.NewTask(TaskTypeCreatePurchaseOrderPdf, jsonPayload, opts...)

	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("could not enqueue task: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueued task")

	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskCreatePurchaseOrderPdf(ctx context.Context, task *asynq.Task) error {
	var payload CreatePurchaseOrderPayload
	conf := processor.conf

	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("could not unmarshal payload: %w", asynq.SkipRetry)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Msg("processing task")

	purchaseOrder, err := processor.db.GetPurchaseOrder(ctx, payload.PurchaseOrderID)
	if err != nil {
		return fmt.Errorf("could not get purchase order by id: %w", err)
	}
	company, err := processor.db.GetCompanyByID(ctx, int64(purchaseOrder.CompanyID))
	if err != nil {
		return fmt.Errorf("could not get company by id: %w", err)
	}

	signatory, err := processor.db.GetSignatoryById(ctx, int64(purchaseOrder.SignatoryID))
	if err != nil {
		return fmt.Errorf("could not get signatory by id: %w", err)
	}

	purchaseOrderItems, err := processor.db.GetPurchaseOrderItemsByPurchaseOrderID(ctx, purchaseOrder.ID)
	if err != nil {
		return fmt.Errorf("could not get purchase order items by purchase order id: %w", err)
	}
	var QuotationID int32
	if purchaseOrder.QuotationID != nil {
		QuotationID = *purchaseOrder.QuotationID
	}

	// get the quotation number
	var quotationNumber string
	relatedQuotation, err := processor.db.GetQuotationByID(ctx, QuotationID)
	if err != nil {
		return fmt.Errorf("could not get quotation by id: %w", err)
	}
	quotationNumber = relatedQuotation.QuotationNo

	// Create PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	utils.PositionCompanyLogoAndDetailsAtTop(pdf)

	// Title
	pdf.SetFont("Helvetica", "B", 7) // Set a larger font size for the title
	title := "PURCHASE ORDER"

	// Calculate width of the title to center it
	width, _ := pdf.GetPageSize()
	titleWidth := pdf.GetStringWidth(title) + 6 // 6 is additional padding
	offset := (width - titleWidth) / 2.0

	// Set X position to the calculated offset
	pdf.SetX(offset)

	// Print the title
	pdf.Cell(70, 10, title) // Increased height to 10 for better appearance with larger font size
	pdf.Ln(10)

	// Purchase Order Details
	pdf.SetFont("Helvetica", "", 7)
	pdf.Cell(30, 5, fmt.Sprintf("Attn:  %s", purchaseOrder.Attn))
	pdf.Ln(3)

	pdf.Cell(30, 5, fmt.Sprintf("Company:  %s", company.Name))
	pdf.Ln(3)

	pdf.Cell(30, 5, fmt.Sprintf("Quotation: %s", quotationNumber))
	pdf.Ln(3)

	pdf.Cell(30, 5, fmt.Sprintf("Date: %s", purchaseOrder.Date.Format("2006-01-02")))
	pdf.Ln(3)

	pdf.Cell(30, 5, fmt.Sprintf("Purchase Order Number: %s", purchaseOrder.PoNo))
	pdf.Ln(5)

	// Purchase Order Items
	headers := []string{"Description", "Part No", "UOM", "Qty", "Item Price", "Discount", "Net Price", "Net Value", "Currency"}
	headerWidths := []float64{50, 15, 15, 15, 15, 15, 15, 15, 15}

	for i, header := range headers {
		pdf.SetFont("Helvetica", "", 6)
		pdf.CellFormat(headerWidths[i], 10, header, "0", 0, "L", false, 0, "")
	}
	pdf.Ln(-1)

	// Table body
	for _, item := range purchaseOrderItems {
		pdf.SetFont("Helvetica", "", 6)

		// Replacing the plus signs with newline characters in the description
		descriptionWithNewLines := strings.Replace(item.Description, "+", "\n", -1)

		// Using MultiCell for the description so that the newline characters are effective
		pdf.MultiCell(headerWidths[0], 4, descriptionWithNewLines, "0", "L", false)

		// Getting current X and Y position after MultiCell
		currentX, currentY := pdf.GetX(), pdf.GetY()

		// Setting X and Y for the next cells
		pdf.SetXY(currentX+headerWidths[0], currentY-10)

		// Assuming item has fields like PartNo, Uom, Qty, ItemPrice, Discount, NetPrice, NetValue, Currency
		pdf.CellFormat(headerWidths[1], 10, item.PartNo, "0", 0, "L", false, 0, "")
		pdf.CellFormat(headerWidths[2], 10, item.Uom, "0", 0, "L", false, 0, "")
		pdf.CellFormat(headerWidths[3], 10, fmt.Sprintf("%d", item.Qty), "0", 0, "L", false, 0, "")
		pdf.CellFormat(headerWidths[4], 10, fmt.Sprintf("%f", item.ItemPrice), "0", 0, "L", false, 0, "")
		pdf.CellFormat(headerWidths[5], 10, fmt.Sprintf("%v", item.Discount), "0", 0, "L", false, 0, "")
		pdf.CellFormat(headerWidths[6], 10, fmt.Sprintf("%f", item.NetPrice), "0", 0, "L", false, 0, "")
		pdf.CellFormat(headerWidths[7], 10, fmt.Sprintf("%f", item.NetValue), "0", 0, "L", false, 0, "")
		pdf.CellFormat(headerWidths[8], 10, item.Currency, "0", 0, "L", false, 0, "")

		pdf.Ln(-1)
	}

	//// Signatory
	//pdf.Ln(10)
	//pdf.SetFont("Arial", "", 7)
	//pdf.CellFormat(80, 10, "Signatory", "", 0, "L", false, 0, "")
	//pdf.Ln(10)
	//pdf.Cell(30, 5, fmt.Sprintf("Name: %s", signatory.Name))
	//pdf.Ln(5)
	//pdf.Cell(30, 5, fmt.Sprintf("Position: %s", signatory.Title))
	utils.PositionSignatoryAtBottom(pdf, signatory)
	name := fmt.Sprintf("%s.pdf", purchaseOrder.PoNo)
	var pdfBuffer bytes.Buffer

	err = pdf.Output(&pdfBuffer)
	if err != nil {
		return fmt.Errorf("could not create PDF: %w", err)
	}

	//err = pdf.OutputFileAndClose(fmt.Sprintf("%s/%s", QuotationsDir, name))
	//if err != nil {
	//	return fmt.Errorf("could not create PDF: %w", err)
	//}

	// Upload the PDF to S3.

	err = utils.UploadFileToS3Bucket(conf, name, PurchaseOrderDir, pdfBuffer.Bytes())
	if err != nil {
		return fmt.Errorf("could not upload PDF to S3: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("PurchaseOrder", name).Msg("processed task")

	return nil
}
