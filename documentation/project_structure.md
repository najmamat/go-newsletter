# Project Structure Documentation

## Overview
This document outlines the structure of the Go Newsletter project, explaining the purpose and contents of each directory.

## Directory Structure

```
go-newsletter/
├── cmd/                  # Main applications
│   └── server/           # The main application binary
│       └── main.go       # Entry point for the application
├── internal/             # Private application code
│   ├── handlers/         # HTTP handlers for processing requests
│   ├── models/           # Data models and database entities
│   ├── repository/       # Data access layer and database operations
│   └── services/         # Business logic implementation
├── pkg/                  # Library code that can be used by external applications
│   └── shared/           # Code that can be shared with other projects
├── supabase/             # Supabase-related code
│   ├── functions/        # Edge Functions
│   │   ├── function1/    # Individual edge function
│   │   └── function2/    # Individual edge function
│   └── config.toml       # Supabase configuration
└── configs/              # Configuration files

```

## Directory Explanations

### cmd/
Contains the main applications of the project. Each subdirectory represents a standalone executable program:
- `newsletter/`: The main application binary containing the entry point (`main.go`)

### internal/
Contains private application code that's specific to this project and shouldn't be imported by other projects:
- `handlers/`: HTTP request handlers implementing the API endpoints
- `models/`: Data structures and database models
- `repository/`: Data access layer implementing the repository pattern for database operations
- `services/`: Business logic implementation and core functionality

### pkg/
Contains library code that can be imported and used by external applications:
- `shared/`: Reusable components and utilities that could be useful for other projects

### supabase/
Contains all Supabase-related code and configuration:
- `functions/`: Edge Functions for serverless functionality
  - Each subdirectory represents a separate edge function
- `config.toml`: Configuration file for Supabase settings

### configs/
Contains configuration files for different environments (development, staging, production)

## Architecture Layers

The project follows clean architecture principles with clear separation of concerns:

1. **Presentation Layer** (`handlers/`): Handles HTTP requests and responses
2. **Business Logic Layer** (`services/`): Contains core business logic and use cases
3. **Data Access Layer** (`repository/`): Abstracts database operations and data persistence
4. **Data Models** (`models/`): Defines data structures and entities

## Best Practices

1. Keep the `internal/` code private and specific to this project
2. Use `pkg/` for code that might be reused by other projects
3. Maintain clear separation of concerns between handlers, services, repository, and models
4. Repository layer should implement interfaces defined in the service layer
5. Keep configuration in the `configs/` directory
6. Document any new edge functions in the Supabase directory
7. Follow dependency injection principles: handlers depend on services, services depend on repositories

## Notes
- All main application code should go into `internal/`
- Shared utilities should go into `pkg/shared/`
- Configuration files should be environment-specific in `configs/`
- Edge functions should be properly documented in their respective directories
- Repository interfaces should be defined in the service layer to maintain dependency inversion 