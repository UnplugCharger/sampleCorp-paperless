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

type CreateQuotationPayload struct {
	// Quotation id
	QuotationID int32 `json:"quotation_id"`
}

func (distributor *RedisTaskDistributor) DistributeTaskCreateQuotationPdf(ctx context.Context, payload *CreateQuotationPayload, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("could not marshal payload: %w", err)

	}
	task := asynq.NewTask(TaskTypeCreateQuotationPdf, jsonPayload, opts...)

	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("could not enqueue task: %w", err)
	}
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueued task")
	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskCreateQuotationPdf(ctx context.Context, task *asynq.Task) error {
	var payload CreateQuotationPayload
	conf := processor.conf

	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("could not unmarshal payload: %w", asynq.SkipRetry)
	}
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Msg("processing task")

	quotation, err := processor.db.GetQuotationByID(ctx, payload.QuotationID)
	if err != nil {
		return fmt.Errorf("could not get quotation by id: %w", err)
	}

	quotationItems, err := processor.db.GetQuotationItemsByQuotationID(ctx, quotation.ID)
	if err != nil {
		return fmt.Errorf("could not get quotation items by quotation id: %w", err)
	}

	signatory, err := processor.db.GetSignatoryById(ctx, int64(quotation.SignatoryID))
	if err != nil {
		return fmt.Errorf("could not get signatory by id: %w", err)
	}

	company, err := processor.db.GetCompanyByID(ctx, int64(quotation.CompanyID))
	if err != nil {
		return fmt.Errorf("could not get company by id: %w", err)
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	//pdf.ImageOptions(CompanyLogoFile, 0, 0, 60, 0, false, gofpdf.ImageOptions{ImageType: "JPG", ReadDpi: true}, 0, "")
	////pdf.SetY(25)
	//// Company Details next to the logo
	//pdf.SetX(120)
	//pdf.SetFont("Helvetica", "B", 7)
	//pdf.Cell(30, 5, companyLocation)
	//pdf.Ln(5)
	//pdf.SetX(120)
	//pdf.Cell(30, 5, CompanyAddress)
	//pdf.Ln(5)
	//pdf.SetX(120)
	//pdf.Cell(30, 5, CompanyEmail)
	//pdf.Ln(5)
	//pdf.SetX(120)
	//pdf.Cell(30, 5, fmt.Sprintf("Mobile: %s", CompanyPhone))
	//pdf.Ln(10)
	//pdf.SetX(120)
	//pdf.Cell(30, 5, fmt.Sprintf("PIN: %s", companyPin))
	//
	//// Adjust the position of the content to be after the logo.
	//pdf.SetY(50)
	utils.PositionCompanyLogoAndDetailsAtTop(pdf)
	// Title
	pdf.SetFont("Helvetica", "B", 7) // Set a larger font size for the title
	title := "QUOTATION"

	// Calculate width of the title to center it
	width, _ := pdf.GetPageSize()
	titleWidth := pdf.GetStringWidth(title) + 6 // 6 is additional padding
	offset := (width - titleWidth) / 2.0

	// Set X position to the calculated offset
	pdf.SetX(offset)

	// Print the title
	pdf.Cell(70, 10, title) // Increased height to 10 for better appearance with larger font size
	pdf.Ln(10)              // New line after title

	pdf.SetFont("Helvetica", "", 7)
	pdf.Cell(30, 5, fmt.Sprintf("Attn:  %s", quotation.Attn))
	pdf.Ln(3)

	pdf.Cell(30, 5, fmt.Sprintf("Company:  %s", company.Name))
	pdf.Ln(3)

	pdf.Cell(30, 5, fmt.Sprintf("Site: %s", quotation.Site))
	pdf.Ln(3)

	pdf.Cell(30, 5, fmt.Sprintf("Date: %s", quotation.Date.Format("2006-01-02")))
	pdf.Ln(3)

	pdf.Cell(30, 5, fmt.Sprintf("Quotation Number: %s", quotation.QuotationNo))
	pdf.Ln(5)

	// Table header
	headers := []string{"Description", "UoM", "Qty", "LeadTime", "ItemPrice", "Discount", "UnitPrice", "NetPrice", "Currency"}
	headerWidths := []float64{50, 15, 15, 15, 15, 15, 15, 15, 15}

	for i, header := range headers {
		pdf.SetFont("Helvetica", "", 7)
		pdf.CellFormat(headerWidths[i], 10, header, "0", 0, "L", false, 0, "")
	}
	pdf.Ln(-1)

	// Table body
	for _, item := range quotationItems {
		pdf.SetFont("Helvetica", "", 6)

		// Replacing the plus signs with newline characters in the description
		descriptionWithNewLines := strings.Replace(item.Description, "+", "\n", -1)

		// Using MultiCell for the description so that the newline characters are effective
		pdf.MultiCell(headerWidths[0], 4, descriptionWithNewLines, "0", "L", false)

		// Getting current X and Y position after MultiCell
		currentX, currentY := pdf.GetX(), pdf.GetY()

		// Setting X and Y for the next cells
		pdf.SetXY(currentX+headerWidths[0], currentY-10)

		pdf.CellFormat(headerWidths[1], 10, item.Uom, "0", 0, "L", false, 0, "")
		pdf.CellFormat(headerWidths[2], 10, fmt.Sprintf("%d", item.Qty), "0", 0, "L", false, 0, "")
		pdf.CellFormat(headerWidths[3], 10, fmt.Sprintf("%s", item.LeadTime), "0", 0, "L", false, 0, "")
		pdf.CellFormat(headerWidths[4], 10, fmt.Sprintf("%f", item.ItemPrice), "0", 0, "L", false, 0, "")
		pdf.CellFormat(headerWidths[5], 10, fmt.Sprintf("%f", item.Disc), "0", 0, "L", false, 0, "")
		pdf.CellFormat(headerWidths[6], 10, fmt.Sprintf("%f", item.UnitPrice), "0", 0, "L", false, 0, "")
		pdf.CellFormat(headerWidths[7], 10, fmt.Sprintf("%f", item.NetPrice), "0", 0, "L", false, 0, "")
		pdf.CellFormat(headerWidths[8], 10, item.Currency, "0", 0, "L", false, 0, "")

		pdf.Ln(-1)
	}

	utils.PositionSignatoryAtBottom(pdf, signatory)

	name := fmt.Sprintf("%s.pdf", quotation.QuotationNo)
	var pdfBuffer bytes.Buffer

	err = pdf.Output(&pdfBuffer)
	if err != nil {
		return fmt.Errorf("could not create PDF: %w", err)
	}

	// Upload the PDF to S3.

	err = utils.UploadFileToS3Bucket(conf, name, QuotationsDir, pdfBuffer.Bytes())
	if err != nil {
		return fmt.Errorf("could not upload PDF to S3: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("Quotation", name).Msg("processed task")

	return nil

}
