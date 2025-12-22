package handler

import (
	"github.com/AnyGridTech/frappe-nfe-bridge/internal/service/frappe_invoice"

	"github.com/gofiber/fiber/v2"
)

type FrappeInvoiceHandler struct {
	svc service.IssuerService
}

func NewFrappeInvoiceHandler(svc service.IssuerService) *FrappeInvoiceHandler {
	return &FrappeInvoiceHandler{svc: svc}
}

// HandleWebhook triggers the issuance
// POST /webhook/invoices/issue
func (h *FrappeInvoiceHandler) CreateWebhook(c *fiber.Ctx) error {
	// Frappe webhooks usually send the DocType data in the body
	type WebhookPayload struct {
		Name string `json:"name"` // The Invoice ID
		// other event data...
	}

	var req WebhookPayload
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	// Call the service to process this invoice
	resp, err := h.svc.IssueNoteForFrappeInvoice(req.Name)
	if err != nil {
		// Log error but maybe don't return 500 to Frappe if you want to avoid retries loops
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Invoice Issued",
		"nfe_id":  resp.ID,
		"status":  resp.Status,
	})
}
