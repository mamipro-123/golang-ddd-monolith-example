package router

import (
	mailerhandlers "monolith-domain/internal/mailer/application/handlers"
	newsletterhandlers "monolith-domain/internal/newsletter/application/handlers"
	resourcehandlers "monolith-domain/internal/resources/application/handlers"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes registers all routes
func SetupRoutes(app *fiber.App, healthHandler *mailerhandlers.HealthCheckHandler, mailerHandler *mailerhandlers.MailerHandler, newsletterHandler *newsletterhandlers.NewsletterHandler, resourceHandler *resourcehandlers.ResourceHandler) {
	app.Get("/health", healthHandler.Handle)
	app.Post("/send-email", mailerHandler.SendMail)
	app.Post("/send-bulk-email", mailerHandler.SendBulkEmails)
	app.Post("/newsletter/subscribe", newsletterHandler.Subscribe)
	app.Post("/newsletter/unsubscribe", newsletterHandler.Unsubscribe)
	app.Get("/newsletter/subscribers", newsletterHandler.GetAllActiveSubscribers)
	app.Post("/resource", resourceHandler.CreateResource)
	app.Put("/resource/:id", resourceHandler.UpdateResource)
	app.Delete("/resource/:id", resourceHandler.DeleteResource)
	app.Get("/resource/:id", resourceHandler.GetResourceByID)
	app.Get("/resource", resourceHandler.GetResourceByKeyAndLang)
	app.Get("/resource/lang/:lang_code", resourceHandler.GetAllResourcesByLang)
	app.Get("/resources", resourceHandler.GetAllResources)
}
