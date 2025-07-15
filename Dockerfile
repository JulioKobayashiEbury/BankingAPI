# Use a minimal Go base image
FROM golang:1.24.2-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /banking-api ./cmd

# Use a scratch image for the final, small image
FROM alpine/git:latest 
# Using alpine/git for curl, but scratch is fine if no curl needed
# If you don't need curl or other tools in the final image, use FROM scratch
# FROM scratch

# Set the working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /banking-api .

# Expose the port the proxy listens on
EXPOSE 25565

# Run the application
ENTRYPOINT ["/app/banking-api"]