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

## Best Practices

1. Keep the `internal/` code private and specific to this project
2. Use `pkg/` for code that might be reused by other projects
3. Maintain clear separation of concerns between handlers, models, and services
4. Keep configuration in the `configs/` directory
5. Document any new edge functions in the Supabase directory

## Notes
- All main application code should go into `internal/`
- Shared utilities should go into `pkg/shared/`
- Configuration files should be environment-specific in `configs/`
- Edge functions should be properly documented in their respective directories 