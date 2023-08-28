package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/qwetu_petro/backend/db/sqlc"
	"github.com/qwetu_petro/backend/workers"
)

type invoice struct {
	PurchaseOrderNumber string  `json:"purchase_order_number"`
	Attn                string  `json:"attn"`
	CompanyID           int32   `json:"company_id"`
	Site                string  `json:"site"`
	AmountDue           float64 `json:"amount_due"`
	BankDetails         int32   `json:"bank_details"`
	SignatoryID         int32   `json:"signatory_id"`
	SentOrReceived      string  `json:"sent_or_received"`
}

type invoiceItem struct {
	Description string  `json:"description"`
	Uom         string  `json:"uom"`
	Qty         int32   `json:"qty"`
	UnitPrice   float64 `json:"unit_price"`
	NetPrice    float64 `json:"net_price"`
	Currency    string  `json:"currency"`
}

type createInvoiceTxnRequest struct {
	Invoice invoice       `json:"invoice"`
	Items   []invoiceItem `json:"items"`
}

func (server *Server) createInvoice(ctx *gin.Context) {
	var req createInvoiceTxnRequest

	// validate request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	userID, err := getUserIdFromContext(ctx)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	invoiceParams := db.CreateInvoiceParams{
		PurchaseOrderNumber: req.Invoice.PurchaseOrderNumber,
		Attn:                req.Invoice.Attn,
		CompanyID:           req.Invoice.CompanyID,
		Site:                req.Invoice.Site,
		AmountDue:           req.Invoice.AmountDue,
		BankDetails:         req.Invoice.BankDetails,
		SignatoryID:         req.Invoice.SignatoryID,
		SentOrReceived:      req.Invoice.SentOrReceived,
	}
	invoiceItemsParams := make([]db.CreateInvoiceItemParams, len(req.Items))
	for i, item := range req.Items {
		invoiceItemsParams[i] = db.CreateInvoiceItemParams{
			Description: item.Description,
			Uom:         item.Uom,
			Qty:         item.Qty,
			UnitPrice:   item.UnitPrice,
			NetPrice:    item.NetPrice,
			Currency:    item.Currency,
		}
	}

	// create invoice
	invoiceTxnResult, err := server.store.CreateInvoiceTxn(ctx, userID, invoiceParams, invoiceItemsParams)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	// create task payload
	taskPayload := &workers.CreateInvoicePdfPayload{
		InvoiceID: invoiceTxnResult.Invoice.ID,
	}

	err = server.taskDistributor.DistributeTaskCreateInvoicePdf(ctx, taskPayload)
	if err != nil {
		ctx.JSON(500, "failed to create invoice pdf")
		return
	}

	ctx.JSON(200, invoiceTxnResult)

}

type listInvoiceRequest struct {
	PageID   int32 `form:"page_id" binding:"required"`
	PageSize int32 `form:"page_size" binding:"required"`
}
type dbInvoice struct {
	db.Invoice
	Items []db.InvoiceItem `json:"items"`
}

type listInvoiceResponse struct {
	Invoices []dbInvoice `json:"invoices"`
}

func (server *Server) listInvoice(ctx *gin.Context) {
	// Commenting out user role check for faster development.
	// Later, don't forget to uncomment this for production.

	// userRoles := ctx.MustGet("roles").([]int64)
	// fmt.Println("user roles", userRoles)
	// if len(userRoles) == 0 {
	//    ctx.JSON(401, "you don't have permission to access this resource")
	//    return
	// }

	// switch {
	// case contains(userRoles, 1):
	var req listInvoiceRequest

	// validate request
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	args := db.ListInvoicesParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	invoices, err := server.store.ListInvoices(ctx, args)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	var dbInvoices []dbInvoice
	for _, invoice := range invoices {
		items, err := server.store.ListInvoiceItemsByInvoiceID(ctx, invoice.ID)
		if err != nil {
			ctx.JSON(500, errorResponse(err))
			return
		}

		dbInvoice := dbInvoice{
			Invoice: invoice,
			Items:   items,
		}
		dbInvoices = append(dbInvoices, dbInvoice)
	}

	response := listInvoiceResponse{
		Invoices: dbInvoices,
	}

	ctx.JSON(200, response)

	// } Uncomment this line when enabling the role check again
}

type UpdateInvoiceRequest struct {
	InvoiceID           int64    `json:"invoice_id"`
	PurchaseOrderNumber *string  `json:"purchase_order_number"`
	Attn                *string  `json:"attn"`
	CompanyID           *int32   `json:"company_id"`
	Site                *string  `json:"site"`
	AmountDue           *float64 `json:"amount_due"`
	BankDetails         *int32   `json:"bank_details"`
	SignatoryID         *int32   `json:"signatory_id"`
	SentOrReceived      *string  `json:"sent_or_received"`
}

func (server *Server) updateInvoice(ctx *gin.Context) {
	ctx.JSON(200, "update invoice")
}

func (server *Server) approveInvoice(ctx *gin.Context) {
	ctx.JSON(200, "approve invoice")
}

func (server *Server) downloadInvoicePdf(ctx *gin.Context) {
	ctx.JSON(200, "download invoice")
}

func contains(slice []int64, val int64) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
