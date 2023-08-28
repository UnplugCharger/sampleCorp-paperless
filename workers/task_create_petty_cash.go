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

type CreatePettyCashPayload struct {
	// Petty cash id
	PettyCashID int32 `json:"petty_cash_id"`
}

func (distributor *RedisTaskDistributor) DistributeTaskCreatePettyCashPdf(ctx context.Context, payload *CreatePettyCashPayload, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("could not marshal payload: %w", err)

	}

	task := asynq.NewTask(TaskTypeCreatePettyCashPdf, jsonPayload, opts...)

	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("could not enqueue task: %w", err)

	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueued task")
	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskCreatePettyCashPdf(ctx context.Context, task *asynq.Task) error {
	var payload CreatePettyCashPayload
	conf := processor.conf

	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("could not unmarshal payload: %w", asynq.SkipRetry)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Msg("processing task")

	pettyCash, err := processor.db.GetPettyCash(ctx, payload.PettyCashID)
	if err != nil {
		return fmt.Errorf("could not get petty cash by id: %w", err)
	}

	// Create pdf
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	utils.PositionCompanyLogoAndDetailsAtTop(pdf)

	//pdf.ImageOptions(CompanyLogoFile, 0, 0, 60, 0, false, gofpdf.ImageOptions{ImageType: "JPG", ReadDpi: true}, 0, "")
	////pdf.SetY(25)
	//// Company Details next to the logo
	//pdf.SetX(120)
	//pdf.SetFont("Helvetica", "B", 11)
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

	name := fmt.Sprintf("%s.pdf", pettyCash.PettyCashNo)

	//err = pdf.OutputFileAndClose(fmt.Sprintf("%s/%s.pdf", PettyCashDir, name))
	var pdfBuffer bytes.Buffer

	err = pdf.Output(&pdfBuffer)
	if err != nil {
		return fmt.Errorf("could not create PDF: %w", err)
	}

	// Upload the PDF to S3.

	err = utils.UploadFileToS3Bucket(conf, name, PettyCashDir, pdfBuffer.Bytes())
	if err != nil {
		return fmt.Errorf("could not upload PDF to S3: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("PettyCash", name).Msg("processed task")

	return nil
}
