# Use an official lightweight Go parent image
FROM golang:alpine AS builder

# Set the current working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the binary with specific settings
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main .


# New stage, based on alpine
FROM alpine:latest

# Set the current working directory inside the container
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .
RUN chmod 755 main

# Expose port for the application
EXPOSE 8443

# Copy the web folder; required in order to render static files
COPY web web

# Run the binary
CMD ["./main"]