package handlers

import (
	"monolith-domain/internal/mailer/application/services"

	"github.com/gofiber/fiber/v2"
)

// MailerHandler handles email-related requests
type MailerHandler struct {
	mailerService *services.MailerService
}

// NewMailerHandler initializes a new MailerHandler
func NewMailerHandler(mailerService *services.MailerService) *MailerHandler {
	return &MailerHandler{
		mailerService: mailerService,
	}
}

type SendMailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
	IsHTML  bool   `json:"is_html"`
}

type SendBulkEmailsRequest struct {
	Recipients []string `json:"recipients"`
	Subject    string   `json:"subject"`
	Body       string   `json:"body"`
	IsHTML     bool     `json:"is_html"`
}

// SendMail handles sending a single email
func (h *MailerHandler) SendMail(c *fiber.Ctx) error {
	var req SendMailRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.mailerService.SendMail(req.To, req.Subject, req.Body, req.IsHTML); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send email",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Email sent successfully",
	})
}

// SendBulkEmails handles sending emails to multiple recipients
func (h *MailerHandler) SendBulkEmails(c *fiber.Ctx) error {
	var req SendBulkEmailsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.mailerService.SendBulkEmails(req.Recipients, req.Subject, req.Body, req.IsHTML); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send bulk emails",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Bulk emails sent successfully",
	})
}
