package helper

import (
	"bytes"
	"github.com/jung-kurt/gofpdf"
)

func GenerateMedicalRecordPDF(patientName, diagnosis, prescription, careSuggestion string) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.SetTextColor(51, 122, 183)
	pdf.CellFormat(0, 10, "Medical Record", "", 0, "C", false, 0, "")

	pdf.Ln(20)
	pdf.SetFont("Arial", "B", 14)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(0, 10, "Patient Information", "", 1, "L", false, 0, "")

	pdf.SetFillColor(240, 240, 240)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(40, 10, "Field", "1", 0, "", true, 0, "")
	pdf.CellFormat(0, 10, "Details", "1", 1, "", true, 0, "")

	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(40, 10, "Patient Name", "1", 0, "", false, 0, "")
	pdf.CellFormat(0, 10, patientName, "1", 1, "", false, 0, "")

	pdf.CellFormat(40, 10, "Diagnosis", "1", 0, "", false, 0, "")
	pdf.MultiCell(0, 10, diagnosis, "1", "L", false)

	pdf.CellFormat(40, 10, "Prescription", "1", 0, "", false, 0, "")
	pdf.MultiCell(0, 10, prescription, "1", "L", false)

	pdf.CellFormat(40, 10, "Care Suggestion", "1", 0, "", false, 0, "")
	pdf.MultiCell(0, 10, careSuggestion, "1", "L", false)

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
