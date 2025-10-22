package controller

import (
	"github.com/Qodarrz/fiber-app/middleware"
	"github.com/Qodarrz/fiber-app/service"
	"github.com/gofiber/fiber/v2"
)

type GeminiController struct {
	service *service.GeminiService
}

func NewGeminiController(s *service.GeminiService) *GeminiController {
	return &GeminiController{service: s}
}

func (gc *GeminiController) Generate(c *fiber.Ctx) error {
	var req struct {
		Prompt string `json:"prompt"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	result, err := gc.service.Generate(req.Prompt)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"result": result})
}

func InitGeminiController(app *fiber.App, s *service.GeminiService, mw *middleware.Middlewares) {
	gc := NewGeminiController(s)

	api := app.Group("/api/gemini")
	api.Post("/generate", gc.Generate)
}
