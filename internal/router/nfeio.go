package router

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

type Nfeio struct {
	Log *log.Logger
	App *fiber.App
}

