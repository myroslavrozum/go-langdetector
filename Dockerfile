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

#Stage 2: pack binary
FROM gruebel/upx:latest AS packer
COPY --from=builder /app/go-langdetector go-langdetector
RUN upx --best --lzma  go-langdetector

#Stage 2: Create the minimal production runtime image
FROM scratch

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/webapp/assets webapp/assets
COPY --from=builder /app/webapp/templates webapp/templates

# Copy the compiled binary from the builder stage
COPY --from=builder /app/go-langdetector go-langdetector


# # Expose the port your web application listens on
EXPOSE 8080

ENV GIN_MODE=release

# # Run the binary
ENTRYPOINT ["/app/go-langdetector"]