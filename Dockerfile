# Dockerfile

# ---- Build Stage ----
FROM golang:1.22-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
# CGO_ENABLED=0 is important for a static build, GOOS=linux to specify the target OS
# -ldflags "-s -w" strips debug information and symbols, reducing binary size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /go-newsletter cmd/server/main.go

# ---- Final Stage ----
FROM alpine:latest

# Security non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /go-newsletter /app/go-newsletter

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
USER appuser
CMD ["/app/go-newsletter"] 