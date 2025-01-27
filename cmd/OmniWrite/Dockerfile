# Stage 1: Build stage
FROM golang:1.23-alpine AS build

# Set the working directory
WORKDIR /app

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o OmniWrite ./cmd/OmniWrite

# Stage 2: Final stage
FROM alpine:edge

# Set the timezone and install CA certificates
RUN apk --no-cache add ca-certificates tzdata

# Create a non-root user and group
RUN addgroup -g 1001 -S appgroup && adduser -u 1001 -S appuser -G appgroup

# Set the working directory
WORKDIR /home/appuser/

# Copy the binary from the build stage
COPY --from=build /app/OmniWrite .

# Change ownership of the application binary to the non-root user
RUN chown 1001:1001 OmniWrite && chmod +x OmniWrite

# Set ENV variables
ENV PORT=80
ENV DATA_SOURCE_NAME="root:Password1!@tcp(localhost)/omni?parseTime=true"
ENV CONN_MAX_LIFETIME=3
ENV MAX_IDLE_CONNECTIONS=10
ENV MAX_OPEN_CONNECTIONS=10

EXPOSE 80

# Switch to the non-root user
USER 1001

# Set the entrypoint command
ENTRYPOINT ["./OmniWrite"]
