package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"

	"github.com/AnyGridTech/frappe-nfe-bridge/internal/handler"
	"github.com/AnyGridTech/frappe-nfe-bridge/internal/repository"
	FrappeInvoiceService "github.com/AnyGridTech/frappe-nfe-bridge/internal/service/frappe_invoice"
	NfeIoInvoiceService "github.com/AnyGridTech/frappe-nfe-bridge/internal/service/nfeio_invoice"
)

func main() {
	// 1. Load Configuration (Fail fast if missing)
	cfg := loadConfig()

	// 2. Initialize Repositories
	// Note: We inject the specific Custom DocType name here
	// This keeps the repo generic enough to handle other DocTypes if needed later.
	frappeRepo := repository.NewFrappeRepo(
		cfg.FrappeURL,
		cfg.FrappeKey,
		cfg.FrappeSecret,
		cfg.CustomDoctype, // exact name of your custom doctype
	)

	nfeRepo := repository.NewNFeRepo(
		cfg.NFeEndpoint,
		cfg.NFeEndpointConsult,
		cfg.NFeAPIKey,
	)

	// 3. Initialize Services
	// We inject the NFe Company ID here as it's a business rule constant for the issuer
	frappe_invoice_service := FrappeInvoiceService.NewIssuerService(frappeRepo, nfeRepo, cfg.NFeCompanyID)
	nfeio_invoice_service := NfeIoInvoiceService.NewFrappeService()

	// 4. Initialize Handlers
	FrappeInvoice := handler.NewFrappeInvoiceHandler(frappe_invoice_service)
	NfeIoInvoice := handler.NewNfeIoInvoiceHandler(nfeio_invoice_service)

	// 5. Setup Fiber
	app := fiber.New(fiber.Config{
		AppName: "Frappe-NFe Integration API",
	})

	// Middleware
	app.Use(logger.New())  // Request logging
	app.Use(recover.New()) // Prevent crashes from panics

	// 6. Define Routes
	// Grouping routes is good practice for versioning
	v1 := app.Group("/api/v1")

	// Webhook endpoint that Frappe will call
	v1.Post("/webhook/invoices/issue", FrappeInvoice.CreateWebhook)
	v1.Post("/webhook/nfeio/response", NfeIoInvoice.ProcessResponseWebhook)

	// Health check (always good to have)
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// 7. Start Server
	log.Printf("Starting server on port %s...", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}

// --- Configuration Helper ---

type Config struct {
	Port               string
	FrappeURL          string
	FrappeKey          string
	FrappeSecret       string
	NFeAPIKey          string
	NFeCompanyID       string
	NFeEndpoint        string
	NFeEndpointConsult string
	CustomDoctype      string
}

func loadConfig() Config {
	loadEnv()
	// Helper to get env with default
	get := func(key, fallback string) string {
		if val, ok := os.LookupEnv(key); ok {
			return val
		}
		return fallback
	}

	// In a real production app, consider using "github.com/spf13/viper"
	// or "github.com/joho/godotenv" here.
	cfg := Config{
		Port:               get("PORT", "3000"),
		FrappeURL:          os.Getenv("FRAPPE_URL"),
		FrappeKey:          os.Getenv("FRAPPE_API_KEY"),
		FrappeSecret:       os.Getenv("FRAPPE_API_SECRET"),
		NFeAPIKey:          os.Getenv("NFE_API_KEY"),
		NFeCompanyID:       os.Getenv("NFE_COMPANY_ID"),
		NFeEndpoint:        get("NFE_ENDPOINT", "https://api.nfe.io/v2"),
		NFeEndpointConsult: get("NFE_ENDPOINT_CONSULT", "https://api.nfe.io/v2"),
		CustomDoctype:      os.Getenv("CUSTOM_DOCTYPE"),
	}

	// Basic validation
	if cfg.FrappeURL == "" || cfg.NFeAPIKey == "" {
		log.Fatal("CRITICAL: Missing environment variables (FRAPPE_URL or NFE_API_KEY)")
	}

	return cfg
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		err := godotenv.Load("../.env")
		if err != nil {
			log.Print("Warning: Error loading .env file")
		}
	}
}
