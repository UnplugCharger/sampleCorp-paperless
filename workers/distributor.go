package workers

import (
	"context"
	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	// DistributeTaskCreateInvoicePdf DistributeTask distributes a task to a worker
	DistributeTaskCreateInvoicePdf(ctx context.Context, payload *CreateInvoicePdfPayload, opts ...asynq.Option) error
	// DistributeTaskCreatePurchaseOrderPdf DistributeTaskCreatePurchaseOrder DistributeTask distributes a task to a worker
	DistributeTaskCreatePurchaseOrderPdf(ctx context.Context, payload *CreatePurchaseOrderPayload, opts ...asynq.Option) error
	// DistributeTaskCreatePettyCashPdf DistributeTaskCreatePettyCash DistributeTask distributes a task to a worker
	DistributeTaskCreatePettyCashPdf(ctx context.Context, payload *CreatePettyCashPayload, opts ...asynq.Option) error
	// DistributeTaskCreatePaymentRequestPdf DistributeTaskCreatePaymentRequest DistributeTask distributes a task to a worker
	DistributeTaskCreatePaymentRequestPdf(ctx context.Context, payload *CreatePaymentRequestPayload, opts ...asynq.Option) error
	// DistributeTaskCreateQuotationPdf DistributeTaskCreateQuotation DistributeTask distributes a task to a worker
	DistributeTaskCreateQuotationPdf(ctx context.Context, payload *CreateQuotationPayload, opts ...asynq.Option) error
}

type RedisTaskDistributor struct {
	// redis connection
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOpts asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOpts)
	return &RedisTaskDistributor{client: client}
}
