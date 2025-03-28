package router

import (
	"monolith-domain/internal/mailer/application/handlers"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes registers all routes
func SetupRoutes(app *fiber.App, healthHandler *handlers.HealthCheckHandler, mailerHandler *handlers.MailerHandler) {
	app.Get("/health", healthHandler.Handle)
	app.Post("/send-email", mailerHandler.SendMail)
	app.Post("/send-bulk-email", mailerHandler.SendBulkEmails)
}
