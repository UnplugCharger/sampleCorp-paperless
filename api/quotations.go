package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	db "github.com/qwetu_petro/backend/db/sqlc"
	"github.com/qwetu_petro/backend/utils"
	"github.com/qwetu_petro/backend/workers"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
)

type quotation struct {
	Attn           string  `json:"attn"`
	CompanyID      int32   `json:"company_id"`
	Site           string  `json:"site"`
	Validity       int32   `json:"validity"`
	Warranty       int32   `json:"warranty"`
	PaymentTerms   string  `json:"payment_terms"`
	DeliveryTerms  string  `json:"delivery_terms"`
	SignatoryID    int32   `json:"signatory_id"`
	Status         string  `json:"status"`
	SentOrReceived *string `json:"sent_or_received"`
}

type quotationItem struct {
	Description string  `json:"description"`
	Uom         string  `json:"uom"`
	Qty         int32   `json:"qty"`
	LeadTime    string  `json:"lead_time"`
	ItemPrice   float64 `json:"item_price"`
	Disc        float64 `json:"disc"`
	UnitPrice   float64 `json:"unit_price"`
	NetPrice    float64 `json:"net_price"`
	Currency    string  `json:"currency"`
}

type createQuotationTxnRequest struct {
	Quotation quotation       `json:"quotation"`
	Items     []quotationItem `json:"items"`
}

func (server *Server) createQuotation(ctx *gin.Context) {
	var req createQuotationTxnRequest

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

	quotationParams := db.CreateQuotationParams{
		Attn:           req.Quotation.Attn,
		CompanyID:      req.Quotation.CompanyID,
		Site:           req.Quotation.Site,
		Validity:       req.Quotation.Validity,
		Warranty:       req.Quotation.Warranty,
		PaymentTerms:   req.Quotation.PaymentTerms,
		DeliveryTerms:  req.Quotation.DeliveryTerms,
		SignatoryID:    req.Quotation.SignatoryID,
		Status:         req.Quotation.Status,
		SentOrReceived: req.Quotation.SentOrReceived,
	}

	quotationItemsParams := make([]db.CreateQuotationItemParams, len(req.Items))

	for i, item := range req.Items {
		quotationItemsParams[i] = db.CreateQuotationItemParams{
			Description: item.Description,
			Uom:         item.Uom,
			Qty:         item.Qty,
			LeadTime:    item.LeadTime,
			ItemPrice:   item.ItemPrice,
			Disc:        item.Disc,
			UnitPrice:   item.UnitPrice,
			NetPrice:    item.NetPrice,
			Currency:    item.Currency,
		}
	}

	args := db.CreateQuotationTxParams{
		Quotation:      quotationParams,
		QuotationItems: quotationItemsParams,
	}

	// execute transaction
	quotationTxnResult, err := server.store.CreateQuotationTxn(ctx, userID, args)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	// create task payload
	taskPayload := &workers.CreateQuotationPayload{
		QuotationID: quotationTxnResult.Quotation.ID,
	}

	// create task
	err = server.taskDistributor.DistributeTaskCreateQuotationPdf(ctx, taskPayload)
	if err != nil {
		ctx.JSON(500, "failed to create quotation pdf ")
		return
	}

	ctx.JSON(200, quotationTxnResult)
}

type dbQuotation struct {
	Quotation db.Quotation `json:"quotation"`
	Items     []db.QuotationItem
}

// ListQuotationsRequest
type listQuotationsRequest struct {
	PageID   int32 `form:"page_id" binding:"required"`
	PageSize int32 `form:"page_size" binding:"required"`
}

// ListQuotationsResponse
type listQuotationsResponse struct {
	Quotations []dbQuotation `json:"quotations"`
}

func (server *Server) listQuotations(ctx *gin.Context) {
	// Commenting out user role check for faster development.
	// Later, don't forget to uncomment this for production.

	// userRoles := ctx.MustGet("roles").([]int64)
	// fmt.Println("user roles", userRoles)
	// if len(userRoles) == 0 {
	//     ctx.JSON(401, "you don't have permission to access this resource")
	//     return
	// }

	// switch {
	// case contains(userRoles, 1): // Assuming 1 is the role ID that has permission to view quotations
	var req listQuotationsRequest

	// validate request
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	args := db.ListQuotationsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	quotations, err := server.store.ListQuotations(ctx, args)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	var dbQuotations []dbQuotation
	for _, quotation := range quotations {
		items, err := server.store.ListQuotationItemsByQuotationID(ctx, quotation.ID) // Assuming you have a similar function for quotations
		if err != nil {
			ctx.JSON(500, errorResponse(err))
			return
		}

		dbQuotation := dbQuotation{
			Quotation: quotation,
			Items:     items,
		}
		dbQuotations = append(dbQuotations, dbQuotation)
	}

	response := listQuotationsResponse{
		Quotations: dbQuotations,
	}

	ctx.JSON(200, response)

	// } Uncomment this line when enabling the role check again
}

type updateQuotationItemRequest struct {
	QuotationID int32   `json:"quotation_id"`
	Description string  `json:"description"`
	Uom         string  `json:"uom"`
	Qty         int32   `json:"qty"`
	LeadTime    string  `json:"lead_time"`
	ItemPrice   float64 `json:"item_price"`
	Disc        float64 `json:"disc"`
	UnitPrice   float64 `json:"unit_price"`
	NetPrice    float64 `json:"net_price"`
	Currency    string  `json:"currency"`
	ID          int32   `json:"id" required:"true" binding:"required"`
}

type updateQuotationRequest struct {
	Attn           string  `json:"attn"`
	CompanyID      int32   `json:"company_id"`
	Site           string  `json:"site"`
	Validity       int32   `json:"validity"`
	Warranty       int32   `json:"warranty"`
	PaymentTerms   string  `json:"payment_terms"`
	DeliveryTerms  string  `json:"delivery_terms"`
	SignatoryID    int32   `json:"signatory_id"`
	Status         string  `json:"status"`
	SentOrReceived *string `json:"sent_or_received"`
	ID             int32   `json:"id"`
}

type updateQuotationTxnRequest struct {
	ID                    int32                        `json:"id" required:"true" binding:"required"`
	UpdatedQuotation      updateQuotationRequest       `json:"quotation"`
	UpdatedQuotationItems []updateQuotationItemRequest `json:"items"`
}

// updateQuotationRequest
func (server *Server) updateQuotation(ctx *gin.Context) {
	var req updateQuotationTxnRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	userID, err := getUserIdFromContext(ctx)
	if err != nil {
		ctx.JSON(401, errorResponse(err))
		return
	}

	quotationParams := db.UpdateQuotationParams{
		Attn:           req.UpdatedQuotation.Attn,
		CompanyID:      req.UpdatedQuotation.CompanyID,
		Site:           req.UpdatedQuotation.Site,
		Validity:       req.UpdatedQuotation.Validity,
		Warranty:       req.UpdatedQuotation.Warranty,
		PaymentTerms:   req.UpdatedQuotation.PaymentTerms,
		DeliveryTerms:  req.UpdatedQuotation.DeliveryTerms,
		SignatoryID:    req.UpdatedQuotation.SignatoryID,
		Status:         req.UpdatedQuotation.Status,
		SentOrReceived: req.UpdatedQuotation.SentOrReceived,
		ID:             req.ID,
	}

	quotationItemsParams := make([]db.UpdateQuotationItemParams, len(req.UpdatedQuotationItems))

	for i, item := range req.UpdatedQuotationItems {
		quotationItemsParams[i] = db.UpdateQuotationItemParams{
			QuotationID: req.ID,
			Description: item.Description,
			Uom:         item.Uom,
			Qty:         item.Qty,
			LeadTime:    item.LeadTime,
			ItemPrice:   item.ItemPrice,
			Disc:        item.Disc,
			UnitPrice:   item.UnitPrice,
			NetPrice:    item.NetPrice,
			Currency:    item.Currency,
			ID:          item.ID,
		}
	}

	args := db.UpdateQuotationTxParams{
		UpdateQuotation:      quotationParams,
		UpdateQuotationItems: quotationItemsParams,
	}

	quotationTxnResult, err := server.store.UpdateQuotationTxn(ctx, userID, args)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// create task payload
	taskPayload := &workers.CreateQuotationPayload{
		QuotationID: quotationTxnResult.Quotation.ID,
	}

	// create task
	err = server.taskDistributor.DistributeTaskCreateQuotationPdf(ctx, taskPayload)
	if err != nil {
		ctx.JSON(500, "failed to create quotation pdf ")
		return
	}

	ctx.JSON(200, quotationTxnResult)

}

type downloadQuotationRequest struct {
	ID int32 `json:"id" required:"true" binding:"required"`
}

// downloadQuotation
func (server *Server) downloadQuotationPdf(ctx *gin.Context) {
	conf := server.config
	folder := workers.QuotationsDir
	var req downloadQuotationRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	quotationDetails, err := server.store.GetQuotationByID(ctx, req.ID)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	key := fmt.Sprintf("%s.pdf", quotationDetails.QuotationNo)

	presignedUrl, err := utils.GeneratePresignedURL(conf, folder, key)
	if err != nil {
		log.Error().Err(err).Msg("error generating presigned url")
		ctx.JSON(500, errorResponse(err))
		return
	}
	// Get the data
	resp, err := http.Get(presignedUrl)
	if err != nil {
		log.Error().Err(err).Msg("error getting the pdf file from presigned url")
		ctx.JSON(500, errorResponse(err))
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error().Err(err).Msg("error closing response body")
		}
	}(resp.Body)

	ctx.Header("Content-Disposition", "attachment; filename="+key)
	ctx.Header("Content-Type", "application/pdf")
	io.Copy(ctx.Writer, resp.Body)

}
