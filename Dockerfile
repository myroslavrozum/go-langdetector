# Stage 1: Build the Go binary
FROM golang:1.26.4-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Install SSL certificates (required for making HTTPS requests)
RUN apk --no-cache add ca-certificates

# Copy dependency files first to leverage Docker layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Compile the application as a statically linked binary
# CGO_ENABLED=0 removes C library dependencies so it can run on 'scratch'
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o go-langdetector .

#Stage 2: Create the minimal production runtime image
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/go-langdetector /app/go-langdetector


# # Expose the port your web application listens on
EXPOSE 8080

# # Run the binary
ENTRYPOINT ["/app/go-langdetector"]