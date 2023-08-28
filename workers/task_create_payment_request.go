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
)

type CreatePaymentRequestPayload struct {
	// PaymentRequest  id
	PaymentRequestID int32 `json:"payment_request_id"`
}

func (distributor *RedisTaskDistributor) DistributeTaskCreatePaymentRequestPdf(ctx context.Context, payload *CreatePaymentRequestPayload, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(TaskTypeCreatePaymentRequestPdf, jsonPayload, opts...)

	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueued task")

	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskCreatePaymentRequestPdf(ctx context.Context, task *asynq.Task) error {
	var payload CreatePaymentRequestPayload
	conf := processor.conf

	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Msg("processing task")

	paymentRequest, err := processor.db.GetPaymentRequest(ctx, payload.PaymentRequestID)
	if err != nil {
		return fmt.Errorf("failed to get payment request by id: %w", err)
	}

	employee, err := processor.db.GetUserById(ctx, int64(paymentRequest.EmployeeID))
	if err != nil {
		return fmt.Errorf("failed to get employee by id: %w", err)
	}
	var adminName string
	if paymentRequest.AdminID != nil {
		admin, err := processor.db.GetUserById(ctx, int64(*paymentRequest.AdminID))
		if err != nil {
			return fmt.Errorf("failed to get admin by id: %w", err)
		}

		adminName = admin.FullName

	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	utils.PositionCompanyLogoAndDetailsAtTop(pdf)

	// Title
	pdf.SetFont("Helvetica", "B", 7)
	pdf.Cell(30, 5, "To: DIRECTOR")
	pdf.SetX(120)
	pdf.Cell(30, 5, fmt.Sprintf("DATE : %s", paymentRequest.RequestDate.Format("02/01/2006")))
	pdf.Ln(10)
	//pdf.SetFont("Helvetica", "", 10)
	pdf.Cell(30, 5, fmt.Sprintf("REQUEST FOR  : %s ", paymentRequest.Currency))
	pdf.SetX(120)
	pdf.Cell(30, 5, fmt.Sprintf("No.  %s", paymentRequest.PaymentRequestNo))
	pdf.Ln(5)

	// Payment Request Details
	pdf.SetFont("Helvetica", "", 7)
	pdf.Cell(30, 5, "I request Cheque/Cash/MPESA payment of the above amount for the following charges :")
	pdf.Ln(15)

	headers := []string{"Ser No.", "Details", "KSH/USD", "CTS"}
	headersWidth := []float64{20, 100, 40, 30}

	for i, header := range headers {
		pdf.SetFont("Helvetica", "B", 7)
		pdf.CellFormat(headersWidth[i], 7, header, "1", 0, "C", false, 0, "")
	}

	// TODO get the payment request details
	pdf.Ln(-1)
	pdf.SetFont("Helvetica", "", 7)

	pdf.CellFormat(headersWidth[0], 7, fmt.Sprintf("%d", paymentRequest.RequestID), "1", 0, "L", false, 0, "")

	// Details
	details := fmt.Sprintf("%s", paymentRequest.Description)
	pdf.CellFormat(headersWidth[1], 7, details, "1", 0, "L", false, 0, "")

	// KSH/USD (I'm assuming this is amount - you might need to adjust this based on what this column represents)
	pdf.CellFormat(headersWidth[2], 7, fmt.Sprintf("%.2f", paymentRequest.Amount), "1", 0, "L", false, 0, "")

	// CTS (I don't know what this represents so leaving it empty)
	pdf.CellFormat(headersWidth[3], 7, "00", "1", 0, "L", false, 0, "")

	// Add a new line
	pdf.Ln(10)
	pdf.SetY(220)

	// Amount in words
	amountInWords := fmt.Sprintf("Amount in words: %s", paymentRequest.AmountInWords)
	pdf.CellFormat(0, 7, amountInWords, "", 0, "L", false, 0, "")
	pdf.Ln(-1)

	// Add Name of Applicant (I'm assuming this information should be fetched separately, placeholder used for now)
	applicant := fmt.Sprintf("Name of Applicant: %s", employee.FullName)
	pdf.CellFormat(0, 7, applicant, "", 0, "L", false, 0, "")

	// Add a new line
	pdf.Ln(-1)

	// Add Department
	pdf.CellFormat(0, 7, "Department: "+employee.Department, "", 0, "L", false, 0, "")

	// Add a new line
	pdf.Ln(-1)

	// Add Signature placeholder
	pdf.CellFormat(0, 7, "Signature: ", "", 0, "L", false, 0, "")

	// Add a new line
	pdf.Ln(-1)

	// Add Approval text
	pdf.CellFormat(0, 7, "Approval: The above requisition is approved", "", 0, "L", false, 0, "")

	// Add a new line
	pdf.Ln(-1)

	// Add Signature placeholder
	pdf.CellFormat(0, 7, "Signature: ", "", 0, "L", false, 0, "")

	// Add a new line
	pdf.Ln(-1)

	// Add Director
	director := fmt.Sprintf("Director: %s", adminName)
	pdf.CellFormat(0, 7, director, "", 0, "L", false, 0, "")

	// Add a new line
	pdf.Ln(-1)

	// Add Date (using ApprovalDate from PaymentRequest struct)
	if paymentRequest.ApprovalDate != nil {
		approvalDate := paymentRequest.ApprovalDate.Format("02-01-2006")
		pdf.CellFormat(0, 7, "Date: "+approvalDate, "", 0, "L", false, 0, "")
	} else {
		pdf.CellFormat(0, 7, "Date: ", "", 0, "L", false, 0, "")
	}

	name := fmt.Sprintf("%s.pdf", paymentRequest.PaymentRequestNo)
	//err = pdf.OutputFileAndClose(fmt.Sprintf("%s/%s.pdf", PaymentRequestDir, name))
	//if err != nil {
	//	return fmt.Errorf("could not create PDF: %w", err)
	//}
	var pdfBuffer bytes.Buffer

	err = pdf.Output(&pdfBuffer)
	if err != nil {
		return fmt.Errorf("could not create PDF: %w", err)
	}

	// Upload the PDF to S3.

	err = utils.UploadFileToS3Bucket(conf, name, PaymentRequestDir, pdfBuffer.Bytes())
	if err != nil {
		return fmt.Errorf("could not upload PDF to S3: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("PaymentRequest", name).Msg("processed task")

	return nil
}
