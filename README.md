# Modular Golang Service Platform

## Project Overview

This is a scalable Golang project designed with a flexible, domain-driven architecture to support multiple services and future expansions.

## Architecture Overview

### Design Principles

- **Modular Architecture**: Easily extendable structure
- **Domain-Driven Design (DDD)**: Clean separation of concerns
- **Bounded Context**: Isolated service implementations
- **Clean Architecture**: Clear separation of layers
- **Observability**: Built-in logging and metrics

## Project Structure

```plaintext
.
├── .deploy/                # Deployment configurations
│   ├── prometheus/        # Prometheus configurations
│   │   └── prometheus.yml
│   └── grafana/          # Grafana configurations
│       └── dashboard.json
├── api-docs/             # API documentation
│   └── api_docs.json     # API specification
├── cmd/                  # Application entry points
│   └── server
│       └── main.go       # Main application startup
├── config/              # Application configuration
│   └── config.yaml      # YAML configuration file
├── internal/            # Internal application packages
│   ├── mailer/         # Mailer bounded context
│   │   ├── application # Application layer
│   │   │   ├── handlers    # HTTP request handlers
│   │   │   └── services    # Business logic services
│   │   ├── domain      # Domain layer
│   │   │   ├── errors.go
│   │   │   └── repositories
│   │   └── infrastructure # Infrastructure layer
│   │       └── smtp_mailer.go
│   ├── newsletter/     # Newsletter bounded context
│   │   ├── application
│   │   │   ├── handlers
│   │   │   └── services
│   │   ├── domain
│   │   └── infrastructure
│   ├── resources/      # Dynamic Resource Management bounded context
│   │   ├── application
│   │   │   ├── handlers
│   │   │   └── services
│   │   ├── domain
│   │   └── infrastructure
│   └── sharedkernel/   # Shared kernel
└── pkg/                # External packages
    ├── config/        # Configuration management
    ├── observability/ # Logging and metrics
    └── router/        # Routing configuration
```

## Architectural Approach

### Extensibility

The project is structured to allow seamless addition of new services. Here's how to add a new bounded context:

1. Create a new directory in `internal/` for your bounded context:

```plaintext
internal/
└── newservice/
    ├── application/
    │   ├── handlers/
    │   └── services/
    ├── domain/
    └── infrastructure/
```

2. Update the service initialization in `main.go`:

```go
func initializeServices(cfg *config.Config, logger *zap.Logger) (*services.MailerService, *services.NewService, error) {
    // Existing service initialization...
    
    // New service initialization
    newService := services.NewService(...)
    
    return mailerService, newService, nil
}
```

3. Add handlers:

```go
func initializeHandlers(mailerService *services.MailerService, newService *services.NewService) (*handlers.HealthCheckHandler, *handlers.MailerHandler, *handlers.NewHandler) {
    return handlers.NewHealthCheckHandler(),
           handlers.NewMailerHandler(mailerService),
           handlers.NewHandler(newService)
}
```

4. Configure routes:

```go
func SetupRoutes(app *fiber.App, healthHandler *handlers.HealthCheckHandler, mailerHandler *handlers.MailerHandler, newHandler *handlers.NewHandler) {
    // Existing routes...
    
    // New service routes
    app.Post("/new-service", newHandler.Create)
    app.Get("/new-service/:id", newHandler.Get)
}
```

### Monitoring and Observability

The application includes built-in monitoring capabilities:

1. **Prometheus Metrics**:
   - HTTP request counts
   - Request durations
   - In-flight requests
   - Custom metrics for each bounded context

2. **Grafana Dashboards**:
   - Pre-configured dashboards for monitoring
   - Customizable metrics visualization

3. **Structured Logging**:
   - JSON-formatted logs
   - Contextual information
   - Error tracking

## Available Services

### Mailer Service
- Send individual emails
- Send bulk emails
- SMTP configuration support
- Email templates

### Newsletter Service
- Subscribe to newsletter
- Unsubscribe from newsletter
- Get all active subscribers with pagination
- Email verification (Basic)

### Resource Management Service
- Dynamic content management
- Multi-language support
- CRUD operations for resources
- Key-value based content storage
- Language-specific content retrieval
- Pagination support for listing resources

## API Endpoints

### Mailer Endpoints
- `POST /send-email`: Send individual email
- `POST /send-bulk-email`: Send bulk emails

### Newsletter Endpoints
- `POST /newsletter/subscribe`: Subscribe to newsletter
- `POST /newsletter/unsubscribe`: Unsubscribe from newsletter
- `GET /newsletter/subscribers`: Get all active subscribers

### Resource Endpoints
- `POST /resource`: Create new resource
- `PUT /resource/:id`: Update existing resource
- `DELETE /resource/:id`: Delete resource
- `GET /resource/:id`: Get resource by ID
- `GET /resource`: Get resource by key and language
- `GET /resource/lang/:lang_code`: Get all resources by language
- `GET /resources`: Get all resources with pagination

## Getting Started

### Prerequisites

- Go 1.20+
- Docker and Docker Compose
- Basic understanding of DDD principles

### Configuration

The application uses YAML configuration:

```yaml
smtp:
  host: smtp.example.com
  port: 587
  user: your_username@example.com
  pass: your_secure_password
  from: sender@example.com
  secure: false

database:
  host: localhost
  port: 5432
  user: postgres
  password: your_password
  dbname: your_database
  sslmode: disable
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: 1h

server:
  port: 3000
  metrics_port: 9090
```

### Setup

1. Clone the repository
2. Configure `config/config.yaml`
3. Run `go mod tidy`
4. Start monitoring services:
   ```bash
   docker compose up -d
   ```
5. Start the application:
   ```bash
   go run cmd/server/main.go
   ```

## Architecture Benefits

- **Decoupled Service Design**: Each bounded context is isolated
- **Easy to Add New Features**: Clear structure for adding new services
- **Maintainable Codebase**: Organized and well-documented
- **Clear Separation of Concerns**: Domain, application, and infrastructure layers
- **Built-in Monitoring**: Prometheus and Grafana integration
- **Structured Logging**: JSON-formatted logs with context
- **Configuration Management**: YAML-based configuration
- **Graceful Shutdown**: Proper handling of application termination
- **Multi-language Support**: Built-in support for multiple languages
- **Pagination Support**: Efficient handling of large datasets
- **UUID Support**: Built-in UUID generation and management
- **Database Migrations**: Automatic schema management

## Planned Service Expansion

- Authentication Service
- User Management
- Notification Service
- Analytics Service
- Audit Service
- Content Management Service
- Media Management Service
- Cache Service
- Search Service


