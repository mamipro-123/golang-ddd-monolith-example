
# Modular Golang Service Platform

## Project Overview

This is a scalable Golang project designed with a flexible, domain-driven architecture to support multiple services and future expansions.

## Architecture Overview

### Design Principles

-   **Modular Architecture**: Easily extendable structure
-   **Domain-Driven Design (DDD)**: Clean separation of concerns
-   **Bounded Context**: Isolated service implementations

## Project Structure

```
.
├── .env                    # Environment configuration
├── api-docs                # API documentation
│   └── api_docs.json       # API specification
├── cmd                     # Application entry points
│   └── server
│       └── main.go         # Main application startup
├── go.mod                  # Go module dependencies
├── go.sum                  # Dependency checksums
├── internal                # Internal application packages
│   ├── mailer              # Current service bounded context
│   │   ├── application     # Application layer
│   │   │   ├── handlers    # HTTP request handlers
│   │   │   └── services    # Business logic services
│   │   ├── domain          # Domain models and logic
│   │   │   ├── errors.go
│   │   │   └── repositories# Domain repositories
│   │   └── infrastructure  # Infrastructure implementations
│   │   │   ├── smtp_mailer.go
│   └── sharedkernel        # Shared kernel for cross-cutting concerns
└── pkg                     # External packages
    └── router              # Routing configuration

```

## Architectural Approach

### Extensibility

The current project is structured to allow seamless addition of new services:

-   Each service resides in its own bounded context
-   Shared kernel for common functionalities
-   Independent domain, application, and infrastructure layers

## Planned Service Expansion

-   Authentication Service
-   User Management
-   Logging and Monitoring Services

## Getting Started

### Prerequisites

-   Go 1.20+
-   Basic understanding of DDD principles


## Example .env
```ini
# SMTP Server Configuration
# Replace these values with your actual SMTP server details

# SMTP Host Address
SMTP_HOST=smtp.example.com

# SMTP Port (typically 587 for TLS or 465 for SSL)
SMTP_PORT=587

# SMTP Authentication Username
SMTP_USER=your_username@example.com

# SMTP Authentication Password
SMTP_PASS=your_secure_password

# Sender Email Address
SMTP_FROM=sender@example.com

# Enable/Disable Secure Connection
# Use 'true' for secure connections, 'false' for non-secure
SECURE_EMAIL=false

# Additional Notes:
# - Ensure your SMTP credentials are kept confidential
# - Port 587 is commonly used for STARTTLS
# - Port 465 is used for SSL/TLS
# - SECURE_EMAIL determines whether to use a secure connection

```

### Setup

1.  Clone the repository
2.  Configure environment variables
3.  Run  `go mod tidy`
4.  Start the server with  `go run cmd/server/main.go`

## Architecture Benefits

-   Decoupled Service Design
-   Easy to Add New Features
-   Maintainable Codebase
-   Clear Separation of Concerns


