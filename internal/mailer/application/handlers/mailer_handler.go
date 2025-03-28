package handlers

import (
	"monolith-domain/internal/mailer/application/services"

	"github.com/gofiber/fiber/v2"
)

// MailerHandler handles email-related requests
type MailerHandler struct {
	Service *services.MailerService
}

// NewMailerHandler initializes a new MailerHandler
func NewMailerHandler(service *services.MailerService) *MailerHandler {
	return &MailerHandler{Service: service}
}

// SendMail handles sending a single email
func (h *MailerHandler) SendMail(c *fiber.Ctx) error {
	var request struct {
		To      string `json:"to"`
		Subject string `json:"subject"`
		Body    string `json:"body"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	err := h.Service.SendMail(request.To, request.Subject, request.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Email sent successfully"})
}

// SendBulkEmails handles sending emails to multiple recipients
func (h *MailerHandler) SendBulkEmails(c *fiber.Ctx) error {
	var request struct {
		Recipients []string `json:"recipients"`
		Subject    string   `json:"subject"`
		Body       string   `json:"body"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	err := h.Service.SendBulkEmails(request.Recipients, request.Subject, request.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Bulk emails sent successfully"})
}
