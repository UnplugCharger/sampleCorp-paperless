package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/qwetu_petro/backend/db/sqlc"
	"net/http"
)

type createBankDetailsRequest struct {
	BankName      string `json:"bank_name" binding:"required"`
	AccountName   string `json:"account_name" binding:"required"`
	AccountNumber string `json:"account_number" binding:"required"`
	Branch        string `json:"branch" binding:"required"`
	SwiftCode     string `json:"swift_code" binding:"required"`
	Address       string `json:"address" binding:"required"`
	Country       string `json:"country" binding:"required"`
	Currency      string `json:"currency" binding:"required"`
	AccountType   string `json:"account_type" binding:"required"`
	CompanyID     int32  `json:"company_id" binding:"required"`
}

type createBankDetailsResponse struct {
	BankName    string `json:"bank_name"`
	AccountName string `json:"account_name"`
	AccountType string `json:"account_type"`
}

func (server *Server) createBankDetails(ctx *gin.Context) {
	// TODO: check if the user is an admin or authorized to create a bank detail
	var req createBankDetailsRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.CreateBankDetailsParams{
		BankName:      req.BankName,
		AccountName:   req.AccountName,
		AccountNumber: req.AccountNumber,
		Branch:        req.Branch,
		SwiftCode:     req.SwiftCode,
		Address:       req.Address,
		Country:       req.Country,
		Currency:      req.Currency,
		AccountType:   req.AccountType,
		CompanyID:     req.CompanyID,
	}

	bankDetails, err := server.store.CreateBankDetails(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, createBankDetailsResponse{
		BankName:    bankDetails.BankName,
		AccountName: bankDetails.AccountName,
		AccountType: bankDetails.AccountType,
	})
}

type getBankDetailsRequest struct {
	AccountNumber string `uri:"account_number" binding:"required,min=1"`
}

type getBankDetailsResponse struct {
	BankName      string `json:"bank_name"`
	AccountName   string `json:"account_name"`
	AccountNumber string `json:"account_number"`
	Branch        string `json:"branch"`
	SwiftCode     string `json:"swift_code"`
	Address       string `json:"address"`
	Country       string `json:"country"`
	Currency      string `json:"currency"`
	AccountType   string `json:"account_type"`
	CompanyID     int32  `json:"company_id"`
}

func (server *Server) getBankDetails(ctx *gin.Context) {
	var req getBankDetailsRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	bankDetails, err := server.store.GetBankDetailsByAccountNumber(ctx, req.AccountNumber)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, getBankDetailsResponse{
		BankName:      bankDetails.BankName,
		AccountName:   bankDetails.AccountName,
		AccountNumber: bankDetails.AccountNumber,
		Branch:        bankDetails.Branch,
		SwiftCode:     bankDetails.SwiftCode,
		Address:       bankDetails.Address,
		Country:       bankDetails.Country,
		Currency:      bankDetails.Currency,
		AccountType:   bankDetails.AccountType,
		CompanyID:     bankDetails.CompanyID,
	})
}

// listBanks returns a list of banks
func (server *Server) listBanks(ctx *gin.Context) {
	banks, err := server.store.ListBanks(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, banks)
}
