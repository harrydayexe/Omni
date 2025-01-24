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
RUN CGO_ENABLED=0 GOOS=linux go build -o OmniRead ./cmd/OmniRead

# Stage 2: Final stage
FROM alpine:edge

# Set ENV variables
ENV PORT=80
ENV DATA_SOURCE_NAME="root:Password1!@tcp(localhost)/omni?parseTime=true"
ENV CONN_MAX_LIFETIME=3
ENV MAX_IDLE_CONNECTIONS=10
ENV MAX_OPEN_CONNECTIONS=10

# Set the working directory
WORKDIR /app

# Copy the binary from the build stage
COPY --from=build /app/OmniRead .

# Set the timezone and install CA certificates
RUN apk --no-cache add ca-certificates tzdata

EXPOSE 80

# Set the entrypoint command
ENTRYPOINT ["/app/OmniRead"]
