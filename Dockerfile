FROM golang:latest AS builder


WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build a static binary to avoid missing libraries
RUN go build -o /app/codeverse-code-ex-svc

# Stage 2: Use Debian instead of Alpine-based Docker image
FROM debian:latest

WORKDIR /app

# Install Docker CLI to interact with the Docker socket
RUN apt-get update && apt-get install -y docker.io

# Copy the Go binary from the builder stage
COPY --from=builder /app/codeverse-code-ex-svc /app/codeverse-code-ex-svc

# Ensure execution permission
RUN chmod +x /app/codeverse-code-ex-svc

# Run the application
CMD ["/app/codeverse-code-ex-svc"]

# CMD [ "tail", "-f", "/dev/null" ]
