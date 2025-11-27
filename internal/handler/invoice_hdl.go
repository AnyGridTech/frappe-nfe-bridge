package handler

import (
	"github.com/AnyGridTech/frappe-nfe-bridge/internal/service"

	"github.com/gofiber/fiber/v2"
)

type InvoiceHandler struct {
	svc service.IssuerService
}

func NewInvoiceHandler(svc service.IssuerService) *InvoiceHandler {
	return &InvoiceHandler{svc: svc}
}

// HandleWebhook triggers the issuance
// POST /webhooks/frappe/invoice-submitted
func (h *InvoiceHandler) HandleWebhook(c *fiber.Ctx) error {
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
