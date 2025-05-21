# Go Newsletter Project Setup Guide

## Prerequisites
- Go 1.24 or higher
- Git
- Access to the Supabase project (Project ID: iiivolgfmqsxvlrggwsh)

## Initial Setup Steps

1. **Clone the Repository**
   ```bash
   git clone <repository-url>
   cd go-newsletter
   ```

2. **Install Dependencies**
   ```bash
   go mod download
   ```

3. **Environment Configuration**
   Create a `.env` file in the root directory with the following content:
   ```env
   # Server Configuration
   PORT=8080

   # Database Configuration (Supabase Pooler Connection)
   PGPASSWORD=<your-database-password>
   PGHOST=aws-0-us-east-2.pooler.supabase.com
   PGPORT=6543
   ```
   Replace `<your-database-password>` with the actual database password from Supabase.

4. **Verify Setup**
   Run the server to verify everything is working:
   ```bash
   go run cmd/server/main.go
   ```
   You should see output indicating successful database connection and server startup.

## Available Endpoints

- Health Check: `GET http://localhost:8080/health`
- List Profiles: `GET http://localhost:8080/api/v1/profiles`

## Project Structure
```
go-newsletter/
├── cmd/
│   └── server/
│       └── main.go       # Main application entry point
├── documentation/        # Project documentation
├── .env                 # Environment variables (not in git)
├── .gitignore
├── go.mod              # Go module definition
└── go.sum              # Go module checksums
```

## Common Issues

1. **Database Connection Issues**
   - Ensure your `.env` file exists and contains the correct credentials
   - Verify you have access to the Supabase project
   - Check that the password is correctly set without any surrounding quotes or spaces

2. **Server Already Running**
   - If you get a "port already in use" error, ensure no other instance of the server is running
   - You can change the PORT in the `.env` file if needed

## Development Workflow

1. Always pull the latest changes before starting work:
   ```bash
   git pull origin main
   ```

2. Create a new branch for your features:
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. Run the server locally while developing:
   ```bash
   go run cmd/server/main.go
   ```

## Additional Resources

- [Project Architecture Documentation](./architecture_1_supabase_explanation.md)
- [API Documentation](./project-kickoff/openapi.yaml)
- [Work Division](./project-kickoff/work_division_supabase_monolith.md) 