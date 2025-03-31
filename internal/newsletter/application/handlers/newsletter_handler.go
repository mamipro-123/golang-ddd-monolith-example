package handlers

import (
	"monolith-domain/internal/newsletter/application/services"
	"github.com/gofiber/fiber/v2"
)

type NewsletterHandler struct {
	service *services.NewsletterService
}

func NewNewsletterHandler(service *services.NewsletterService) *NewsletterHandler {
	return &NewsletterHandler{service: service}
}

type SubscribeRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type UnsubscribeRequest struct {
	Token string `json:"token" validate:"required"`
}

type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Total      int64      `json:"total"`
	Page       int        `json:"page"`
	Size       int        `json:"size"`
	TotalPages int        `json:"total_pages"`
}

func (h *NewsletterHandler) Subscribe(c *fiber.Ctx) error {
	var req SubscribeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	newsletter, err := h.service.Subscribe(req.Email)
	if err != nil {
		if err.Error() == "email already subscribed" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Email already subscribed",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to subscribe",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Successfully subscribed",
		"data":    newsletter,
	})
}

func (h *NewsletterHandler) Unsubscribe(c *fiber.Ctx) error {
	var req UnsubscribeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.service.Unsubscribe(req.Token); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to unsubscribe",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Successfully unsubscribed",
	})
}

func (h *NewsletterHandler) GetAllActiveSubscribers(c *fiber.Ctx) error {
	// Sayfalama parametrelerini al
	page := c.QueryInt("page", 1)
	size := c.QueryInt("size", 10)

	// Sayfa ve boyut kontrolü
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	// Aboneleri getir
	newsletters, total, err := h.service.GetAllActiveSubscribers(page, size)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch subscribers",
		})
	}

	// Toplam sayfa sayısını hesapla
	totalPages := int((total + int64(size) - 1) / int64(size))

	response := PaginationResponse{
		Data:       newsletters,
		Total:      total,
		Page:       page,
		Size:       size,
		TotalPages: totalPages,
	}

	return c.JSON(response)
} 

