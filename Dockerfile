# Use official Go image
FROM golang:1.25-alpine3.21

# Install git (needed for go mod download)
RUN apk add --no-cache git

# Set working directory inside the container
WORKDIR /app/restapiwithgin

# Copy go.mod and go.sum first (for dependency caching)
COPY restapiwithgin/go.mod restapiwithgin/go.sum ./

# Download Go dependencies
RUN go mod download

# Copy only the restapiwithgin source code
COPY restapiwithgin .

# Optional: list files to verify copy (debugging)
RUN echo "=== DEBUG: Listing /app/restapiwithgin contents ===" \
    && ls -la /app/restapiwithgin \
    && echo "=== End of /app/restapiwithgin listing ==="

# Build the binary (output to /app directory)
RUN go build -o /app/main-restapi .

# Set working directory back to /app
WORKDIR /app

# Expose port (if your app runs on 8080)
EXPOSE 8080

# Command to run the application
CMD ["./main-restapi"]
