package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/qwetu_petro/backend/db/sqlc"
)

type createSignatoryRequest struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}

func (server *Server) createSignatory(ctx *gin.Context) {
	var req createSignatoryRequest

	// validate request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	signatoryParams := db.CreateSignatoryParams{
		Name:  req.Name,
		Title: req.Title,
	}

	signatory, err := server.store.CreateSignatory(ctx, signatoryParams)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, signatory)
}

type listSignatoriesRequest struct {
	PageID   int32 `form:"page_id" binding:"required"`
	PageSize int32 `form:"page_size" binding:"required"`
}

func (server *Server) listSignatories(ctx *gin.Context) {
	var req listSignatoriesRequest

	// validate request
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	args := db.ListSignatoriesParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	signatories, err := server.store.ListSignatories(ctx, args)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, signatories)
}
