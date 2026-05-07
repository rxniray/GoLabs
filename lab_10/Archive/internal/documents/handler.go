package documents

import (
	"bytes"
	"fmt"

	"github.com/go-pdf/fpdf"
	"github.com/gofiber/fiber/v2"
)

type generateRequest struct {
	Template string         `json:"template"`
	Data     map[string]any `json:"data"`
}

type Handler struct{}

func NewHandler() *Handler { return &Handler{} }

// POST /api/documents/generate
func (h *Handler) Generate(c *fiber.Ctx) error {
	var req generateRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	gen, ok := registry[req.Template]
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest,
			fmt.Sprintf("unknown template %q; available: invoice, report", req.Template))
	}

	pdf, err := gen(req.Data)
	if err != nil {
		return fmt.Errorf("generate pdf: %w", err)
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.pdf"`, req.Template))
	return c.Send(pdf)
}

// registry maps template name → generator func
var registry = map[string]func(map[string]any) ([]byte, error){
	"invoice": generateInvoice,
	"report":  generateReport,
}

func generateInvoice(data map[string]any) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "INVOICE")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)
	fields := []string{"number", "client", "amount"}
	for _, f := range fields {
		v := ""
		if val, ok := data[f]; ok {
			v = fmt.Sprintf("%v", val)
		}
		pdf.Cell(0, 8, fmt.Sprintf("%s: %s", f, v))
		pdf.Ln(8)
	}

	return pdfBytes(pdf)
}

func generateReport(data map[string]any) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "REPORT")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)
	for k, v := range data {
		pdf.Cell(0, 8, fmt.Sprintf("%s: %v", k, v))
		pdf.Ln(8)
	}

	return pdfBytes(pdf)
}

func pdfBytes(pdf *fpdf.Fpdf) ([]byte, error) {
	if err := pdf.Error(); err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
