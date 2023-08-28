package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/qwetu_petro/backend/db/sqlc"
	"net/http"
)

type createCompanyRequest struct {
	Name     string  `json:"name" binding:"required"`
	Initials string  `json:"initials" binding:"required"`
	Address  *string `json:"address"`
}

func (server *Server) createCompany(ctx *gin.Context) {
	var req createCompanyRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.CreateCompanyParams{
		Name:     req.Name,
		Initials: req.Initials,
		Address:  req.Address,
	}

	company, err := server.store.CreateCompany(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, company)
}

type getCompanyRequest struct {
	Name string `uri:"name" binding:"required"`
}

type deleteCompany struct {
	Name string `json:"name"`
}

func (server *Server) deleteCompany(ctx *gin.Context) {
	var req deleteCompany

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.store.DeleteCompanyByName(ctx, req.Name)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "Company deleted successfully")

}

type listCompaniesRequest struct {
	PageID   int32 `form:"page_id" binding:"required"`
	PageSize int32 `form:"page_size" binding:"required"`
}

// List all companies
func (server *Server) listCompanies(ctx *gin.Context) {
	var req listCompaniesRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListCompaniesParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	companies, err := server.store.ListCompanies(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, companies)

}
