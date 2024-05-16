# Use a multi-stage build for smaller image size
FROM golang:1.22.3 AS builder

# Set the current working directory inside the container
WORKDIR /app

# Copy the Go modules and build files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app 
RUN CGO_ENABLED=0 GOOS=linux go build -o dns-proxy

# Use a lightweight Alpine image as the base image for the final container
FROM alpine:latest

RUN apk update && apk upgrade

# Create a non-root user and group for running the application
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set the working directory inside the container
WORKDIR /app

# Copy the executable from the builder stage into the final image
COPY --from=builder /app/dns-proxy .

# Set ownership and permissions for the executable
RUN chown appuser:appgroup /app/dns-proxy && \
    chmod 755 /app/dns-proxy

# Expose TCP and UDP ports
EXPOSE 53
EXPOSE 54

# Switch to non-root user
USER appuser

# Run the DNS proxy server when the container starts
CMD ["./dns-proxy"]
