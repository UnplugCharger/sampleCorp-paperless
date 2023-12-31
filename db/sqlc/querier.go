// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	ApprovePaymentRequest(ctx context.Context, arg ApprovePaymentRequestParams) (PaymentRequest, error)
	ApprovePettyCash(ctx context.Context, arg ApprovePettyCashParams) (PettyCash, error)
	ApprovePurchaseOrder(ctx context.Context, arg ApprovePurchaseOrderParams) (PurchaseOrder, error)
	CreateBankDetails(ctx context.Context, arg CreateBankDetailsParams) (BankDetail, error)
	CreateCar(ctx context.Context, arg CreateCarParams) (Car, error)
	CreateCompany(ctx context.Context, arg CreateCompanyParams) (Company, error)
	CreateFuelConsumption(ctx context.Context, arg CreateFuelConsumptionParams) (CarFuelConsumption, error)
	CreateInvoice(ctx context.Context, arg CreateInvoiceParams) (Invoice, error)
	CreateInvoiceItem(ctx context.Context, arg CreateInvoiceItemParams) (InvoiceItem, error)
	CreatePaymentRequest(ctx context.Context, arg CreatePaymentRequestParams) (PaymentRequest, error)
	CreatePettyCash(ctx context.Context, arg CreatePettyCashParams) (PettyCash, error)
	CreatePurchaseOrder(ctx context.Context, arg CreatePurchaseOrderParams) (PurchaseOrder, error)
	CreatePurchaseOrderItem(ctx context.Context, arg CreatePurchaseOrderItemParams) (PurchaseOrderItem, error)
	CreateQuotation(ctx context.Context, arg CreateQuotationParams) (Quotation, error)
	CreateQuotationItem(ctx context.Context, arg CreateQuotationItemParams) (QuotationItem, error)
	// description: Create a role
	CreateRole(ctx context.Context, arg CreateRoleParams) (Role, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	CreateSignatory(ctx context.Context, arg CreateSignatoryParams) (Signatory, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	CreateUserRoles(ctx context.Context, arg CreateUserRolesParams) (UserRole, error)
	DeleteCompanyByName(ctx context.Context, name string) error
	DeletePaymentRequest(ctx context.Context, requestID int32) error
	DeletePettyCash(ctx context.Context, transactionID int32) error
	// description: Delete a role
	DeleteRole(ctx context.Context, id int64) error
	DeleteSignatoryByName(ctx context.Context, name string) error
	DeleteUser(ctx context.Context, id int64) error
	DeleteUserRoles(ctx context.Context, userID int64) error
	DeleteUserRolesByRole(ctx context.Context, roleID int64) error
	GetBankDetailsByAccountNumber(ctx context.Context, accountNumber string) (BankDetail, error)
	GetBankInfoByID(ctx context.Context, id int32) (BankDetail, error)
	GetCarById(ctx context.Context, id int32) (Car, error)
	GetCars(ctx context.Context, arg GetCarsParams) ([]Car, error)
	GetCompanyByID(ctx context.Context, id int64) (Company, error)
	GetConsumptionByCar(ctx context.Context, carID int32) ([]CarFuelConsumption, error)
	GetConsumptionByCarAndDateRange(ctx context.Context, arg GetConsumptionByCarAndDateRangeParams) ([]CarFuelConsumption, error)
	GetConsumptions(ctx context.Context, arg GetConsumptionsParams) ([]CarFuelConsumption, error)
	GetInvoiceById(ctx context.Context, id int32) (Invoice, error)
	GetInvoiceItemsByInvoiceID(ctx context.Context, invoiceID int32) ([]InvoiceItem, error)
	GetInvoicesByPurchaseOrderNumber(ctx context.Context, purchaseOrderNumber string) ([]Invoice, error)
	GetPaymentRequest(ctx context.Context, requestID int32) (PaymentRequest, error)
	GetPettyCash(ctx context.Context, transactionID int32) (PettyCash, error)
	GetPurchaseOrder(ctx context.Context, id int32) (PurchaseOrder, error)
	GetPurchaseOrderItemsByPurchaseOrderID(ctx context.Context, purchaseOrderID int32) ([]PurchaseOrderItem, error)
	GetQuotationByID(ctx context.Context, id int32) (Quotation, error)
	GetQuotationItemsByQuotationID(ctx context.Context, quotationID int32) ([]QuotationItem, error)
	// description: Get a role
	GetRole(ctx context.Context, id int64) (Role, error)
	GetSessionByID(ctx context.Context, id uuid.UUID) (Session, error)
	GetSignatoryById(ctx context.Context, id int64) (Signatory, error)
	GetUserById(ctx context.Context, id int64) (User, error)
	GetUserByUserNameOrEmail(ctx context.Context, arg GetUserByUserNameOrEmailParams) (User, error)
	GetUserRole(ctx context.Context, arg GetUserRoleParams) (UserRole, error)
	GetUserRoles(ctx context.Context, userID int64) ([]UserRole, error)
	GetUsers(ctx context.Context, arg GetUsersParams) ([]User, error)
	ListBankDetailsByBankName(ctx context.Context, bankName string) ([]BankDetail, error)
	ListBanks(ctx context.Context) ([]BankDetail, error)
	ListCompanies(ctx context.Context, arg ListCompaniesParams) ([]Company, error)
	ListEmployeePaymentRequests(ctx context.Context, arg ListEmployeePaymentRequestsParams) ([]PaymentRequest, error)
	ListEmployeePettyCash(ctx context.Context, arg ListEmployeePettyCashParams) ([]PettyCash, error)
	ListInvoiceItemsByInvoiceID(ctx context.Context, invoiceID int32) ([]InvoiceItem, error)
	ListInvoices(ctx context.Context, arg ListInvoicesParams) ([]Invoice, error)
	ListPaymentRequests(ctx context.Context, arg ListPaymentRequestsParams) ([]PaymentRequest, error)
	ListPettyCash(ctx context.Context, arg ListPettyCashParams) ([]PettyCash, error)
	ListPurchaseOrderItemsByPurchaseOrderID(ctx context.Context, purchaseOrderID int32) ([]PurchaseOrderItem, error)
	ListPurchaseOrders(ctx context.Context, arg ListPurchaseOrdersParams) ([]PurchaseOrder, error)
	ListQuotationItemsByQuotationID(ctx context.Context, quotationID int32) ([]QuotationItem, error)
	ListQuotations(ctx context.Context, arg ListQuotationsParams) ([]Quotation, error)
	// description: List roles
	ListRoles(ctx context.Context, arg ListRolesParams) ([]Role, error)
	ListSignatories(ctx context.Context, arg ListSignatoriesParams) ([]Signatory, error)
	ListSiteInvoices(ctx context.Context, arg ListSiteInvoicesParams) ([]Invoice, error)
	ListUserCompanyInvoices(ctx context.Context, arg ListUserCompanyInvoicesParams) ([]Invoice, error)
	UpdatePaymentRequest(ctx context.Context, arg UpdatePaymentRequestParams) (PaymentRequest, error)
	UpdatePettyCash(ctx context.Context, arg UpdatePettyCashParams) (PettyCash, error)
	UpdatePurchaseOrder(ctx context.Context, arg UpdatePurchaseOrderParams) (PurchaseOrder, error)
	UpdatePurchaseOrderItem(ctx context.Context, arg UpdatePurchaseOrderItemParams) (PurchaseOrderItem, error)
	UpdateQuotation(ctx context.Context, arg UpdateQuotationParams) (Quotation, error)
	UpdateQuotationItem(ctx context.Context, arg UpdateQuotationItemParams) (QuotationItem, error)
	// description: Update a role
	UpdateRole(ctx context.Context, arg UpdateRoleParams) (Role, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
	UpdateUserRoles(ctx context.Context, arg UpdateUserRolesParams) (UserRole, error)
}

var _ Querier = (*Queries)(nil)
