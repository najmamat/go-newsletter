# Go Newsletter - Implementation Guide

## Overview

This document describes the architectural changes made to the Go Newsletter application and provides guidance for implementing additional endpoints. The application has been refactored to use dependency injection, OpenAPI code generation, and a clean architecture pattern.

## Table of Contents

- [Architecture Changes](#architecture-changes)
- [Build System](#build-system)
- [OpenAPI Code Generation](#openapi-code-generation)
- [Dependency Injection Pattern](#dependency-injection-pattern)
- [Implementing New Endpoints](#implementing-new-endpoints)
- [Database Integration](#database-integration)
- [Development Workflow](#development-workflow)
- [Project Structure](#project-structure)

## Architecture Changes

### Key Architectural Decisions

1. **Dependency Injection**: Moved from direct database access to a clean dependency injection pattern
2. **OpenAPI-First Development**: All API endpoints are defined in OpenAPI spec and code is generated
3. **Layered Architecture**: Clear separation between handlers, services, and repositories
4. **Type Safety**: Generated types from OpenAPI ensure consistency between spec and implementation

### Dependencies Added

The following key dependencies were added to support the new architecture:

```go
// OpenAPI and code generation
github.com/getkin/kin-openapi v0.127.0
github.com/oapi-codegen/runtime v1.1.1

// Database
github.com/jackc/pgx/v5 v5.7.5

// Utilities
github.com/google/uuid v1.5.0
github.com/joho/godotenv v1.5.1
```

## Build System

### Makefile Targets

A comprehensive Makefile has been added with the following targets:

```bash
# Development
make help          # Show available commands
make dev           # Run in development mode with auto-reload
make run           # Build and run the application
make build         # Build the application
make clean         # Clean build artifacts

# Code Generation
make generate      # Generate Go code from OpenAPI specification

# Code Quality
make fmt           # Format Go code
make lint          # Lint Go code
make test          # Run all tests
make tidy          # Clean up dependencies

# Development Setup
make dev-deps      # Install development dependencies
make setup         # Set up development environment

# Docker
make docker-build  # Build Docker image
make docker-run    # Build and run Docker container

# Database (placeholders for future implementation)
make migrate-up    # Run database migrations up
make migrate-down  # Run database migrations down
```

### Usage Examples

```bash
# First time setup
make setup

# Generate code from OpenAPI spec
make generate

# Run in development mode
make dev

# Build for production
make build
```

## OpenAPI Code Generation

### Configuration

The OpenAPI code generation is configured in `api/oapi-config.yaml`:

```yaml
package: generated
generate:
  models: true
  chi-server: true
  client: true
  embedded-spec: true
output: pkg/generated/api.go
output-options:
  skip-fmt: false
  skip-prune: false
```

### Generated Components

The code generation creates:

1. **Models**: All request/response structs defined in OpenAPI spec
2. **Server Interface**: Interface that handlers must implement
3. **Client**: HTTP client for API consumption
4. **Embedded Spec**: OpenAPI spec embedded in the binary

### Key Generated Types

```go
// Authentication
type AuthCredentials struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type AuthResponse struct {
    AccessToken string        `json:"access_token"`
    User        EditorProfile `json:"user"`
}

// Newsletter Management
type Newsletter struct {
    Id          string     `json:"id"`
    EditorId    string     `json:"editor_id"`
    Name        string     `json:"name"`
    Description *string    `json:"description,omitempty"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
}

type NewsletterCreate struct {
    Name        string  `json:"name"`
    Description *string `json:"description,omitempty"`
}

// Publishing
type PublishedPost struct {
    Id           string     `json:"id"`
    NewsletterId string     `json:"newsletter_id"`
    EditorId     string     `json:"editor_id"`
    Title        string     `json:"title"`
    ContentHtml  string     `json:"content_html"`
    ContentText  *string    `json:"content_text,omitempty"`
    Status       string     `json:"status"`
    ScheduledAt  *time.Time `json:"scheduled_at,omitempty"`
    PublishedAt  time.Time  `json:"published_at"`
    CreatedAt    time.Time  `json:"created_at"`
}
```

## Generated Code vs Implementation Pattern

### Why We Have Both `api.go` (Generated) and `server.go` (Implementation)

This is a fundamental Go pattern called **Interface Segregation** combined with **Code Generation**:

#### Generated Code (`pkg/generated/api.go`)
- **Auto-generated** from OpenAPI specification
- Contains the `ServerInterface` (contract/interface definition)
- Contains all type definitions (structs, enums, etc.)
- Contains HTTP client for consuming the API
- **Never edit manually** - gets overwritten on each `make generate`

#### Implementation Code (`internal/server/server.go`)
- **Hand-written** business logic implementation
- Implements the `ServerInterface` from generated code
- Contains service dependencies and orchestration
- Contains error handling, validation, authorization
- **Safe from regeneration** - maintained by developers

#### The Pattern in Action

```go
// GENERATED - pkg/generated/api.go
type ServerInterface interface {
    GetMe(w http.ResponseWriter, r *http.Request)
    PutMe(w http.ResponseWriter, r *http.Request)
    GetNewsletters(w http.ResponseWriter, r *http.Request)
    // ... all endpoints defined in OpenAPI
}

// IMPLEMENTATION - internal/server/server.go
type Server struct {
    profileService    *services.ProfileService
    newsletterService *services.NewsletterService
    logger           *slog.Logger
}

// Server must implement ALL methods from ServerInterface
func (s *Server) GetMe(w http.ResponseWriter, r *http.Request) {
    // Extract user from auth context
    userID := r.Context().Value("userID").(string)
    
    // Call business logic
    profile, err := s.profileService.GetByID(r.Context(), userID)
    if err != nil {
        s.handleError(w, r, err)
        return
    }
    
    // Return response
    s.respondJSON(w, http.StatusOK, profile)
}
```

#### Benefits of This Pattern

1. **Type Safety**: Compiler ensures we implement all required endpoints
2. **Contract Compliance**: Generated interface guarantees API spec compliance
3. **Regeneration Safety**: Can update OpenAPI spec without losing implementation
4. **Separation of Concerns**: Generated types vs business logic
5. **Testability**: Easy to mock the interface for testing
6. **Documentation Sync**: Implementation always matches the documented API

#### Alternative Approaches (Why We Don't Use Them)

```go
// ❌ BAD: Putting logic directly in generated code
// This gets overwritten on regeneration!

// ❌ BAD: Manual HTTP handlers without interface compliance
// No guarantee they match the OpenAPI spec

// ✅ GOOD: Our pattern
// Generated interface + separate implementation = best of both worlds
```

## Dependency Injection Pattern

### Main Application Structure

The main application now follows a dependency injection pattern:

```go
func main() {
    // Load environment
    err := godotenv.Load()
    if err != nil {
        slog.Warn("No .env file found")
    }

    // Initialize database
    dbpool := initDB()
    defer dbpool.Close()

    // Initialize dependencies
    profileRepo := repository.NewProfileRepository(dbpool)
    profileService := services.NewProfileService(profileRepo)
    apiServer := server.NewAPIServer(profileService)

    // Setup router
    router := setupRouter(apiServer)

    // Start server
    port := utils.GetEnvWithDefault("PORT", "8080")
    slog.Info("Starting server", "port", port)
    log.Fatal(http.ListenAndServe(":"+port, router))
}
```

### Layer Responsibilities

1. **Repository Layer**: Database access and data persistence
2. **Service Layer**: Business logic and domain rules
3. **Handler Layer**: HTTP request/response handling and validation
4. **Server Layer**: Orchestrates services and implements OpenAPI interface

### Server Interface Implementation

The server must implement the generated `ServerInterface`:

```go
type APIServer struct {
    profileService *services.ProfileService
    // Add other services as needed
}

func NewAPIServer(profileService *services.ProfileService) *APIServer {
    return &APIServer{
        profileService: profileService,
    }
}

// Implement all OpenAPI endpoints
func (s *APIServer) GetMe(w http.ResponseWriter, r *http.Request) {
    // Implementation
}

func (s *APIServer) PutMe(w http.ResponseWriter, r *http.Request) {
    // Implementation
}
```

## Implementing New Endpoints

### Step-by-Step Process

#### 1. Define in OpenAPI Spec

First, add your endpoint to `api/openapi.yaml`:

```yaml
/newsletters:
  get:
    summary: List Editor's Newsletters
    tags:
      - Newsletters
    security:
      - bearerAuth: []
    responses:
      '200':
        description: A list of newsletters
        content:
          application/json:
            schema:
              type: array
              items:
                $ref: '#/components/schemas/Newsletter'
```

#### 2. Generate Code

Run code generation to create the interface and types:

```bash
make generate
```

#### 3. Implement Repository Layer

Create repository methods for database access:

```go
// internal/repository/newsletter_repository.go
type NewsletterRepository struct {
    db *pgxpool.Pool
}

func NewNewsletterRepository(db *pgxpool.Pool) *NewsletterRepository {
    return &NewsletterRepository{db: db}
}

func (r *NewsletterRepository) GetByEditorID(ctx context.Context, editorID string) ([]Newsletter, error) {
    query := `
        SELECT id, editor_id, name, description, created_at, updated_at 
        FROM newsletters 
        WHERE editor_id = $1
    `
    
    rows, err := r.db.Query(ctx, query, editorID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var newsletters []Newsletter
    for rows.Next() {
        var newsletter Newsletter
        err := rows.Scan(
            &newsletter.ID,
            &newsletter.EditorID,
            &newsletter.Name,
            &newsletter.Description,
            &newsletter.CreatedAt,
            &newsletter.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        newsletters = append(newsletters, newsletter)
    }
    
    return newsletters, nil
}
```

#### 4. Implement Service Layer

Create service methods for business logic:

```go
// internal/services/newsletter_service.go
type NewsletterService struct {
    repo *repository.NewsletterRepository
}

func NewNewsletterService(repo *repository.NewsletterRepository) *NewsletterService {
    return &NewsletterService{repo: repo}
}

func (s *NewsletterService) GetNewslettersByEditor(ctx context.Context, editorID string) ([]generated.Newsletter, error) {
    newsletters, err := s.repo.GetByEditorID(ctx, editorID)
    if err != nil {
        return nil, err
    }
    
    // Convert to generated types
    var result []generated.Newsletter
    for _, newsletter := range newsletters {
        result = append(result, generated.Newsletter{
            Id:          newsletter.ID,
            EditorId:    newsletter.EditorID,
            Name:        newsletter.Name,
            Description: newsletter.Description,
            CreatedAt:   newsletter.CreatedAt,
            UpdatedAt:   newsletter.UpdatedAt,
        })
    }
    
    return result, nil
}
```

#### 5. Implement Handler

Add the handler method to your server:

```go
// internal/server/api_server.go
func (s *APIServer) GetNewsletters(w http.ResponseWriter, r *http.Request) {
    // Extract user ID from context (set by auth middleware)
    userID, ok := r.Context().Value("userID").(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    
    newsletters, err := s.newsletterService.GetNewslettersByEditor(r.Context(), userID)
    if err != nil {
        slog.Error("Failed to get newsletters", "error", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(newsletters)
}
```

#### 6. Update Dependency Injection

Add the new service to your dependency injection:

```go
// main.go
func main() {
    // ... existing code ...
    
    // Initialize repositories
    profileRepo := repository.NewProfileRepository(dbpool)
    newsletterRepo := repository.NewNewsletterRepository(dbpool)
    
    // Initialize services
    profileService := services.NewProfileService(profileRepo)
    newsletterService := services.NewNewsletterService(newsletterRepo)
    
    // Initialize server
    apiServer := server.NewAPIServer(profileService, newsletterService)
    
    // ... rest of the code ...
}
```

### Error Handling Pattern

Use consistent error handling across all endpoints:

```go
func (s *APIServer) SomeEndpoint(w http.ResponseWriter, r *http.Request) {
    // Validation
    if someCondition {
        utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid input")
        return
    }
    
    // Business logic
    result, err := s.someService.DoSomething(r.Context())
    if err != nil {
        slog.Error("Operation failed", "error", err)
        utils.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
        return
    }
    
    // Success response
    utils.WriteJSONResponse(w, http.StatusOK, result)
}
```

### Authentication Middleware

Endpoints requiring authentication should extract user information from context:

```go
func (s *APIServer) AuthenticatedEndpoint(w http.ResponseWriter, r *http.Request) {
    userID, ok := r.Context().Value("userID").(string)
    if !ok {
        utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
        return
    }
    
    // Use userID for authorization checks
    // ...
}
```

## Database Integration

### Connection Pattern

The application uses pgx/v5 for PostgreSQL connectivity:

```go
func initDB() *pgxpool.Pool {
    databaseURL := utils.GetEnvWithDefault("DATABASE_URL", "postgres://localhost/newsletter_dev")
    
    config, err := pgxpool.ParseConfig(databaseURL)
    if err != nil {
        log.Fatal("Failed to parse database config:", err)
    }
    
    dbpool, err := pgxpool.NewWithConfig(context.Background(), config)
    if err != nil {
        log.Fatal("Failed to create connection pool:", err)
    }
    
    return dbpool
}
```

### Query Patterns

Use context-aware queries with proper error handling:

```go
func (r *Repository) GetByID(ctx context.Context, id string) (*Model, error) {
    query := `SELECT id, name, created_at FROM table WHERE id = $1`
    
    var model Model
    err := r.db.QueryRow(ctx, query, id).Scan(
        &model.ID,
        &model.Name,
        &model.CreatedAt,
    )
    
    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, nil // Not found
        }
        return nil, err
    }
    
    return &model, nil
}
```

## Development Workflow

### Daily Development Process

1. **Update OpenAPI Spec**: Modify `api/openapi.yaml`
2. **Generate Code**: Run `make generate`
3. **Implement Layers**: Repository → Service → Handler
4. **Test**: Run `make test`
5. **Format**: Run `make fmt`
6. **Run**: Use `make dev` for development

### Code Generation Workflow

```bash
# After updating OpenAPI spec
make generate

# Check what was generated
git diff pkg/generated/

# Implement required interface methods
# The compiler will tell you what's missing
```

### Testing Strategy

- Unit tests for services and repositories
- Integration tests for handlers
- Use generated types for consistency

## Project Structure

```
go-newsletter/
├── api/
│   ├── openapi.yaml              # OpenAPI specification
│   └── oapi-config.yaml          # OpenAPI code generation config
├── cmd/
│   └── server/                   # Application entry points
├── documentation/
│   └── implementation-guide.md   # This file
├── internal/
│   ├── repository/               # Data access layer
│   ├── services/                 # Business logic layer
│   ├── server/                   # HTTP handlers
│   └── utils/                    # Shared utilities
├── pkg/
│   └── generated/                # Generated OpenAPI code
├── Makefile                      # Build automation
├── main.go                       # Application entry point
├── go.mod                        # Go modules
└── README.md                     # Project overview
```

## Best Practices

### 1. OpenAPI First
- Always define endpoints in OpenAPI spec first
- Run code generation before implementing
- Use generated types throughout the application

### 2. Dependency Injection
- Keep dependencies explicit and injected
- Use interfaces for testability
- Initialize all dependencies in main()

### 3. Error Handling
- Use structured logging with slog
- Return appropriate HTTP status codes
- Provide meaningful error messages

### 4. Context Usage
- Pass context through all layers
- Use context for cancellation and timeouts
- Extract user information from context in handlers

### 5. Database Access
- Use repositories for data access
- Handle pgx.ErrNoRows appropriately
- Use transactions for complex operations

### 6. Code Organization
- Keep handlers thin, business logic in services
- Use consistent naming conventions
- Group related functionality in packages

## Common Patterns

### Request Validation

```go
func (s *APIServer) CreateSomething(w http.ResponseWriter, r *http.Request) {
    var req generated.CreateRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
        return
    }
    
    // Validate required fields
    if req.Name == "" {
        utils.WriteErrorResponse(w, http.StatusBadRequest, "Name is required")
        return
    }
    
    // Continue with business logic...
}
```

### Path Parameter Extraction

```go
func (s *APIServer) GetSomethingById(w http.ResponseWriter, r *http.Request, id string) {
    // id is automatically extracted from path by generated router
    
    result, err := s.service.GetByID(r.Context(), id)
    if err != nil {
        // Handle error
        return
    }
    
    if result == nil {
        utils.WriteErrorResponse(w, http.StatusNotFound, "Not found")
        return
    }
    
    utils.WriteJSONResponse(w, http.StatusOK, result)
}
```

### Authorization Checks

```go
func (s *APIServer) UpdateNewsletter(w http.ResponseWriter, r *http.Request, newsletterId string) {
    userID := r.Context().Value("userID").(string)
    
    // Check ownership
    newsletter, err := s.newsletterService.GetByID(r.Context(), newsletterId)
    if err != nil {
        utils.WriteErrorResponse(w, http.StatusInternalServerError, "Internal error")
        return
    }
    
    if newsletter.EditorId != userID {
        utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied")
        return
    }
    
    // Continue with update logic...
}
```