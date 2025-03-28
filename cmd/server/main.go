package main

import (
	"log"
	"monolith-domain/internal/mailer/infrastructure"
	"monolith-domain/internal/mailer/application/services"
	"monolith-domain/internal/mailer/application/handlers"
	"monolith-domain/pkg/router"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"regexp"
)
func ValidateEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func main() {
	// Fiber başlat
	app := fiber.New()

	// SMTP Mailer başlat
	smtpMailer, err := infrastructure.NewSMTPMailer()
	if err != nil {
		log.Fatalf("Failed to initialize mailer: %v", err)
	}

	// Mailer servisini oluştur
	mailerService := services.NewMailerService(smtpMailer)

	
	// Handlers
	healthHandler := handlers.NewHealthCheckHandler()
	mailerHandler := handlers.NewMailerHandler(mailerService)

	// Router'ı başlat
	router.SetupRoutes(app, healthHandler, mailerHandler)

	fmt.Println(ValidateEmail("muh4mmrd@protonmail.com"))  // true
    fmt.Println(ValidateEmail("invalid.email"))            // false
    fmt.Println(ValidateEmail("user@domain"))              // false
	// Sunucuyu başlat
	log.Fatal(app.Listen(":3000"))


	
	
}
