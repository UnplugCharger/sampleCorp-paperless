package utils

import (
	"fmt"
	"github.com/phpdave11/gofpdf"
	db "github.com/qwetu_petro/backend/db/sqlc"
)

const (
	CompanyLogoFile = "./assets/logo.jpeg"
	CompanyAddress  = "P.O Box 52868-00100, Nairobi, Kenya"
	CompanyEmail    = "info@qwetupetroleum.com"
	CompanyPhone    = "+254 722 253 775 | +255 684 123 735"
	CompanyLocation = "45 Africa Reit,Karen Road,Karen,Nairobi"
	CompanyPin      = "P051693104Y"
)

func PositionSignatoryAndBankDetailsAtBottom(pdf *gofpdf.Fpdf, signatory db.Signatory, bankInfo db.BankDetail) {
	totalContentHeight := pdf.GetY()
	sectionHeight := 46.5 // Height required for signatory and bank details
	_, pageHeight := pdf.GetPageSize()
	remainingSpace := pageHeight - totalContentHeight - sectionHeight
	pdf.Ln(remainingSpace)

	// Add Bank Details
	pdf.SetFont("Arial", "B", 7)
	pdf.CellFormat(80, 10, "Bank Details", "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.Cell(30, 5, fmt.Sprintf("Bank Name: %s", bankInfo.BankName))
	pdf.Ln(3)
	pdf.Cell(30, 5, fmt.Sprintf("Account Number: %s", bankInfo.AccountNumber))
	// Add some spacing between Signatory and Bank Details
	pdf.Ln(5)

	// Add Signatory
	pdf.SetFont("Arial", "B", 7)
	pdf.CellFormat(80, 10, "Signatory", "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.Cell(30, 5, fmt.Sprintf("Name: %s", signatory.Name))
	pdf.Ln(3)
	pdf.Cell(30, 5, fmt.Sprintf("Position: %s", signatory.Title))

}

func PositionCompanyLogoAndDetailsAtTop(pdf *gofpdf.Fpdf) {
	pdf.ImageOptions(CompanyLogoFile, 0, 0, 60, 0, false, gofpdf.ImageOptions{ImageType: "JPG", ReadDpi: true}, 0, "")
	//pdf.SetY(25)
	// Company Details next to the logo
	pdf.SetX(120)
	pdf.SetFont("Helvetica", "B", 7)
	pdf.Cell(30, 5, CompanyLocation)
	pdf.Ln(3)
	pdf.SetX(120)
	pdf.Cell(30, 5, CompanyAddress)
	pdf.Ln(3)
	pdf.SetX(120)
	pdf.Cell(30, 5, CompanyEmail)
	pdf.Ln(3)
	pdf.SetX(120)
	pdf.Cell(30, 5, fmt.Sprintf("Mobile: %s", CompanyPhone))
	pdf.Ln(6)
	pdf.SetX(120)
	pdf.Cell(30, 5, fmt.Sprintf("PIN: %s", CompanyPin))

	// Adjust the position of the content to be after the logo.
	pdf.SetY(50)

}

func PositionSignatoryAtBottom(pdf *gofpdf.Fpdf, signatory db.Signatory) {
	totalContentHeight := pdf.GetY()
	signatoryHeight := 46.5
	_, pageHeight := pdf.GetPageSize()
	remainingSpace := pageHeight - totalContentHeight - signatoryHeight + 3.5
	pdf.Ln(remainingSpace)

	// Add Signatory
	pdf.SetFont("Arial", "B", 7)
	pdf.CellFormat(80, 10, "Signatory", "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.Cell(30, 5, fmt.Sprintf("Name: %s", signatory.Name))
	pdf.Ln(3)
	pdf.Cell(30, 5, fmt.Sprintf("Position: %s", signatory.Title))
}
