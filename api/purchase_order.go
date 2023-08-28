package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/qwetu_petro/backend/db/sqlc"
	"github.com/qwetu_petro/backend/workers"
)

type purchaseOrder struct {
	Attn           string  `json:"attn"`
	CompanyID      int32   `json:"company_id"`
	Address        string  `json:"address"`
	SignatoryID    int32   `json:"signatory_id"`
	QuotationID    *int32  `json:"quotation_id"`
	SentOrReceived *string `json:"sent_or_received"`
}

type purchaseOrderItem struct {
	Description     string  `json:"description"`
	PartNo          string  `json:"part_no"`
	Uom             string  `json:"uom"`
	Qty             int32   `json:"qty"`
	ItemPrice       float64 `json:"item_price"`
	Discount        float64 `json:"discount"`
	NetPrice        float64 `json:"net_price"`
	NetValue        float64 `json:"net_value"`
	Currency        string  `json:"currency"`
	PurchaseOrderID int32   `json:"purchase_order_id"`
}

type createPurchaseOrderTxnRequest struct {
	PurchaseOrder purchaseOrder       `json:"purchase_order"`
	Items         []purchaseOrderItem `json:"items"`
}

func (server *Server) createPurchaseOrder(ctx *gin.Context) {
	var req createPurchaseOrderTxnRequest

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

	purchaseOrderParams := db.CreatePurchaseOrderParams{
		Attn:           req.PurchaseOrder.Attn,
		CompanyID:      req.PurchaseOrder.CompanyID,
		Address:        req.PurchaseOrder.Address,
		SignatoryID:    req.PurchaseOrder.SignatoryID,
		QuotationID:    req.PurchaseOrder.QuotationID,
		SentOrReceived: req.PurchaseOrder.SentOrReceived,
	}

	purchaseOrderItemsParams := make([]db.CreatePurchaseOrderItemParams, len(req.Items))
	for i, item := range req.Items {
		purchaseOrderItemsParams[i] = db.CreatePurchaseOrderItemParams{
			Description:     item.Description,
			PartNo:          item.PartNo,
			Uom:             item.Uom,
			Qty:             item.Qty,
			ItemPrice:       item.ItemPrice,
			Discount:        item.Discount,
			NetPrice:        item.NetPrice,
			NetValue:        item.NetValue,
			Currency:        item.Currency,
			PurchaseOrderID: item.PurchaseOrderID,
		}
	}

	args := db.CreatePurchaseOrderTxParams{
		PurchaseOrder:      purchaseOrderParams,
		PurchaseOrderItems: purchaseOrderItemsParams,
	}

	// execute transaction
	purchaseOrderTxnResult, err := server.store.CreatePurchaseOrderTxn(ctx, args, userID)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	// create task payload
	payload := &workers.CreatePurchaseOrderPayload{
		PurchaseOrderID: purchaseOrderTxnResult.PurchaseOrder.ID,
	}

	// create task
	err = server.taskDistributor.DistributeTaskCreatePurchaseOrderPdf(ctx, payload)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, purchaseOrderTxnResult)

}

type listPurchaseOrderRequest struct {
	PageSize int32 `form:"page_size"`
	PageID   int32 `form:"page_id"`
}

type dbPurchaseOrder struct {
	db.PurchaseOrder
	Items []db.PurchaseOrderItem `json:"items"`
}

type listPurchaseOrderResponse struct {
	PurchaseOrders []dbPurchaseOrder `json:"purchase_orders"`
}

func (server *Server) listPurchaseOrder(ctx *gin.Context) {
	// Commenting out user role check for faster development.
	// Later, don't forget to uncomment this for production.

	// userRoles := ctx.MustGet("roles").([]int64)
	// fmt.Println("user roles", userRoles)
	// if len(userRoles) == 0 {
	//     ctx.JSON(401, "you don't have permission to access this resource")
	//     return
	// }

	// switch {
	// case contains(userRoles, 1):
	var req listPurchaseOrderRequest

	// validate request
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	args := db.ListPurchaseOrdersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	purchaseOrders, err := server.store.ListPurchaseOrders(ctx, args)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	var dbPurchaseOrders []dbPurchaseOrder
	for _, Order := range purchaseOrders {
		items, err := server.store.ListPurchaseOrderItemsByPurchaseOrderID(ctx, Order.ID)
		if err != nil {
			ctx.JSON(500, errorResponse(err))
			return
		}

		dbPurchaseOrder := dbPurchaseOrder{
			PurchaseOrder: Order,
			Items:         items,
		}
		dbPurchaseOrders = append(dbPurchaseOrders, dbPurchaseOrder)
	}

	response := listPurchaseOrderResponse{
		PurchaseOrders: dbPurchaseOrders,
	}

	ctx.JSON(200, response)

	// } Uncomment this line when enabling the role check again
}

type UpdatePurchaseOrderRequest struct {
	Attn           string  `json:"attn"`
	CompanyID      int32   `json:"company_id"`
	Address        string  `json:"address"`
	SignatoryID    int32   `json:"signatory_id"`
	SentOrReceived *string `json:"sent_or_received"`
	ID             int32   `json:"id"`
}

type UpdatePurchaseOrderItemRequest struct {
	ID          int32   `json:"id"`
	Description string  `json:"description"`
	PartNo      string  `json:"part_no"`
	Uom         string  `json:"uom"`
	Qty         int32   `json:"qty"`
	ItemPrice   float64 `json:"item_price"`
	Discount    float64 `json:"discount"`
	NetPrice    float64 `json:"net_price"`
	NetValue    float64 `json:"net_value"`
	Currency    string  `json:"currency"`
}

type UpdatePurchaseOrderTxnRequest struct {
	PurchaseOrder      UpdatePurchaseOrderRequest       `json:"purchase_order"`
	PurchaseOrderItems []UpdatePurchaseOrderItemRequest `json:"items"`
	ID                 int32                            `json:"id"`
}

type UpdatePurchaseOrderTxnResponse struct {
	PurchaseOrderID int32  `json:"purchase_order_id"`
	Message         string `json:"message"`
}

// Update  this is a transaction that updates a purchase order
func (server *Server) updatePurchaseOrder(ctx *gin.Context) {
	var req UpdatePurchaseOrderTxnRequest

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

	purchaseOrderParams := db.UpdatePurchaseOrderParams{
		Attn:           req.PurchaseOrder.Attn,
		CompanyID:      req.PurchaseOrder.CompanyID,
		Address:        req.PurchaseOrder.Address,
		SignatoryID:    req.PurchaseOrder.SignatoryID,
		SentOrReceived: req.PurchaseOrder.SentOrReceived,
		ID:             req.PurchaseOrder.ID,
	}

	purchaseOrderItemsParams := make([]db.UpdatePurchaseOrderItemParams, len(req.PurchaseOrderItems))

	for i, item := range req.PurchaseOrderItems {
		purchaseOrderItemsParams[i] = db.UpdatePurchaseOrderItemParams{
			ID:          item.ID,
			Description: item.Description,
			PartNo:      item.PartNo,
			Uom:         item.Uom,
			Qty:         item.Qty,
			ItemPrice:   item.ItemPrice,
			Discount:    item.Discount,
			NetPrice:    item.NetPrice,
			NetValue:    item.NetValue,
			Currency:    item.Currency,
		}
	}

	args := db.UpdatePurchaseOrderTxParams{
		PurchaseOrder:      purchaseOrderParams,
		PurchaseOrderItems: purchaseOrderItemsParams,
	}

	res, err := server.store.UpdatePurchaseOrderTxn(ctx, args, userID)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	// create task payload
	taskPayload := &workers.CreatePurchaseOrderPayload{
		PurchaseOrderID: res.PurchaseOrder.ID,
	}

	// create task
	err = server.taskDistributor.DistributeTaskCreatePurchaseOrderPdf(ctx, taskPayload)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	response := UpdatePurchaseOrderTxnResponse{
		PurchaseOrderID: res.PurchaseOrder.ID,
		Message:         "Purchase order updated successfully",
	}

	ctx.JSON(200, response)

}

type approvePurchaseOrderRequest struct {
	PurchaseOrderID int32  `json:"purchase_order_id"`
	Status          string `json:"status"`
}

type approvePurchaseOrderResponse struct {
	PurchaseOrderID int32  `json:"purchase_order_id"`
	Message         string `json:"message"`
}

func (server *Server) approvePurchaseOrder(ctx *gin.Context) {
	var req approvePurchaseOrderRequest

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

	args := db.ApprovePurchaseOrderParams{
		ID:         req.PurchaseOrderID,
		PoStatus:   &req.Status,
		ApprovedBy: &userID,
	}

	res, err := server.store.ApprovePurchaseOrderTxn(ctx, args, userID)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	// create task payload
	payload := &workers.CreatePurchaseOrderPayload{
		PurchaseOrderID: res.PurchaseOrder.ID,
	}

	// create task
	err = server.taskDistributor.DistributeTaskCreatePurchaseOrderPdf(ctx, payload)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	response := approvePurchaseOrderResponse{
		PurchaseOrderID: res.PurchaseOrder.ID,
		Message:         "Purchase order approved successfully",
	}

	ctx.JSON(200, response)

}

func (server *Server) downloadPurchaseOrderPdf(ctx *gin.Context) {
	response := gin.H{
		"message": "Purchase order pdf created successfully",
	}

	ctx.JSON(200, response)
}
