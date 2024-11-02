# Use the official Go 1.21 image based on Alpine as the build stage
FROM golang:1.21-alpine as builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum to the working directory
COPY go.mod go.sum ./

# Install build dependencies
RUN apk add --no-cache gcc libc-dev

# Download dependencies
RUN go mod download

# Copy all source code to the working directory
COPY . .

# Compile the application
RUN go build -o generator .

# Use a lightweight Alpine runtime image
FROM alpine:3.20

# Set the working directory
WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates

# Copy the binary file from the build stage to the runtime image
COPY --from=builder /app/generator .

# Execute the application
CMD ["./generator"]
