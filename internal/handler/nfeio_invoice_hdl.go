package handler

import (
	service "github.com/AnyGridTech/frappe-nfe-bridge/internal/service/nfeio_invoice"
	"github.com/gofiber/fiber/v2"
)

type NfeIoInvoiceHandler struct {
	svc service.FrappeService
}

func NewNfeIoInvoiceHandler(svc service.FrappeService) *NfeIoInvoiceHandler {
	return &NfeIoInvoiceHandler{svc: svc}
}

func (h *NfeIoInvoiceHandler) ProcessResponseWebhook(c *fiber.Ctx) error {
	// Implementation here...
	return nil
}