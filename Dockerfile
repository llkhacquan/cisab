FROM golang:1.24-alpine as builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source from the current directory to the working directory inside the container
COPY . .

# Build the Go app
RUN go build -o main .

# Start a new stage from scratch
FROM alpine:latest

WORKDIR /app

# Copy config files
COPY config/ ./config/

# Copy the binary from builder
COPY --from=builder /app/main .

# Command to run
CMD ["./main"]
