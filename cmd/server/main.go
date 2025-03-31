package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	mailerhandlers "monolith-domain/internal/mailer/application/handlers"
	mailerservices "monolith-domain/internal/mailer/application/services"
	mailerinfra "monolith-domain/internal/mailer/infrastructure"
	newsletterhandlers "monolith-domain/internal/newsletter/application/handlers"
	newsletterservices "monolith-domain/internal/newsletter/application/services"
	newsletterinfra "monolith-domain/internal/newsletter/infrastructure"
	"monolith-domain/internal/newsletter/domain"
	resourcehandlers "monolith-domain/internal/resources/application/handlers"
	resourceservices "monolith-domain/internal/resources/application/services"
	resourceinfra "monolith-domain/internal/resources/infrastructure"
	resourcedomain "monolith-domain/internal/resources/domain"
	"monolith-domain/pkg/config"
	"monolith-domain/pkg/observability"
	"monolith-domain/pkg/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// Prometheus metrics
var (
	httpRequestsTotal    *prometheus.CounterVec
	httpRequestDuration  *prometheus.HistogramVec
	httpRequestsInFlight prometheus.Gauge
	registry            *prometheus.Registry
)

func initPrometheusMetrics() {
	registry = prometheus.NewRegistry()

	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"method", "path"},
	)

	httpRequestsInFlight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Current number of HTTP requests being served",
		},
	)

	registry.MustRegister(httpRequestsTotal)
	registry.MustRegister(httpRequestDuration)
	registry.MustRegister(httpRequestsInFlight)
}

func setupMiddlewares(app *fiber.App) {
	app.Use(recover.New())
	app.Use(fiberlogger.New(fiberlogger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} | ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	}))
	
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
	}))
	
	app.Use(observability.MetricsMiddleware)
}

func initializeDatabase(cfg *config.Config, logger *zap.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	gormLog := gormlogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		gormlogger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  gormlogger.Error,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLog,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}
	
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	var result int
	if err := db.Raw("SELECT 1").Scan(&result).Error; err != nil {
		logger.Error("Database connection test failed", zap.Error(err))
		return nil, fmt.Errorf("database connection test failed: %w", err)
	}
	logger.Info("Database connection verified")

	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return nil, fmt.Errorf("failed to create uuid extension: %w", err)
	}

	migrator := db.Migrator()
	newsletterTableExists := migrator.HasTable(&domain.Newsletter{})
	resourceTableExists := migrator.HasTable(&resourcedomain.Resource{})

	if !newsletterTableExists {
		logger.Info("Starting newsletter table migration...")
		if err := db.AutoMigrate(&domain.Newsletter{}); err != nil {
			logger.Error("Newsletter migration failed", zap.Error(err))
			return nil, fmt.Errorf("failed to migrate newsletter table: %w", err)
		}
		logger.Info("Newsletter table migration completed successfully")
	}

	if !resourceTableExists {
		logger.Info("Starting resource table migration...")
		if err := db.AutoMigrate(&resourcedomain.Resource{}); err != nil {
			logger.Error("Resource migration failed", zap.Error(err))
			return nil, fmt.Errorf("failed to migrate resource table: %w", err)
		}
		logger.Info("Resource table migration completed successfully")
	}

	if newsletterTableExists && resourceTableExists {
		logger.Info("Database schema is already up to date")
	}

	return db, nil
}

func initializeServices(cfg *config.Config, logger *zap.Logger) (*mailerservices.MailerService, *newsletterservices.NewsletterService, *resourceservices.ResourceService, error) {
	smtpMailer, err := mailerinfra.NewSMTPMailer(
		cfg.SMTP.Host,
		cfg.SMTP.Port,
		cfg.SMTP.User,
		cfg.SMTP.Pass,
		cfg.SMTP.From,
		cfg.SMTP.Secure,
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to initialize mailer: %w", err)
	}

	db, err := initializeDatabase(cfg, logger)
	if err != nil {
		return nil, nil, nil, err
	}

	mailerService := mailerservices.NewMailerService(smtpMailer)
	newsletterRepo := newsletterinfra.NewPostgresRepository(db)
	newsletterService := newsletterservices.NewNewsletterService(newsletterRepo)
	resourceRepo := resourceinfra.NewPostgresRepository(db)
	resourceService := resourceservices.NewResourceService(resourceRepo)

	return mailerService, newsletterService, resourceService, nil
}

func setupApplication(cfg *config.Config, logger *zap.Logger) (*fiber.App, error) {
	app := fiber.New(fiber.Config{
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Concurrency:  256 * 1024,
	})

	setupMiddlewares(app)

	mailerService, newsletterService, resourceService, err := initializeServices(cfg, logger)
	if err != nil {
		return nil, err
	}

	healthHandler := mailerhandlers.NewHealthCheckHandler()
	mailerHandler := mailerhandlers.NewMailerHandler(mailerService)
	newsletterHandler := newsletterhandlers.NewNewsletterHandler(newsletterService)
	resourceHandler := resourcehandlers.NewResourceHandler(resourceService)

	router.SetupRoutes(app, healthHandler, mailerHandler, newsletterHandler, resourceHandler)
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.HandlerFor(registry, promhttp.HandlerOpts{})))

	return app, nil
}

func gracefulShutdown(app *fiber.App, logger *zap.Logger) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	logger.Info("Initiating graceful shutdown...")

	if err := app.ShutdownWithTimeout(5 * time.Second); err != nil {
		logger.Error("Error during server shutdown", zap.Error(err))
	}

	logger.Info("Server shutdown completed")
}

func main() {
	if err := observability.InitLogger("development"); err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer observability.Sync()

	logger := observability.GetLogger()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	initPrometheusMetrics()
	logger.Info("Starting application...")

	app, err := setupApplication(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to setup application", zap.Error(err))
	}

	logger.Info("Server starting", zap.String("port", cfg.Server.Port))
	go func() {
		if err := app.Listen(":" + cfg.Server.Port); err != nil {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	gracefulShutdown(app, logger)
}
