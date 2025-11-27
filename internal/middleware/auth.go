package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"os"

	"github.com/gofiber/fiber/v2"
)

func verifyWebhookSignature(signature string, body []byte) bool {
	secret := os.Getenv("WEBHOOK_SECRET")
	h := hmac.New(sha256.New, []byte(secret))

	h.Write(body)
	expectedSignature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

func WebhookAuth(c *fiber.Ctx) error {
	signature := c.Get(os.Getenv("WEBHOOK_SIGNATURE"))
	if !verifyWebhookSignature(signature, c.Body()) {
		c.Status(fiber.StatusUnauthorized).SendString("invalid webhook signature")
		return nil
	}
	return c.Next()
}
