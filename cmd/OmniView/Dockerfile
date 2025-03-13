# Stage 1: Build stage
FROM golang:1.24-bookworm AS build

# Set the working directory
WORKDIR /app

# Download Tailwind CSS
RUN curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/download/v4.0.14/tailwindcss-linux-arm64
RUN chmod u+x tailwindcss-linux-arm64
RUN mv tailwindcss-linux-arm64 tailwindcss

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

RUN ./tailwindcss -i "./internal/omniview/templates/custom.css" -o "./internal/omniview/templates/static/style.css"

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o OmniView ./cmd/OmniView

# Stage 2: Final stage
FROM alpine:edge

# Set the timezone and install CA certificates
RUN apk --no-cache add ca-certificates tzdata

# Create a non-root user and group
RUN addgroup -g 1001 -S appgroup && adduser -u 1001 -S appuser -G appgroup

# Set the working directory
WORKDIR /home/appuser/

# Copy the binary from the build stage
COPY --from=build /app/OmniView .

# Change ownership of the application binary to the non-root user
RUN chown 1001:1001 OmniView && chmod +x OmniView

# Set ENV variables
ENV PORT=80

EXPOSE 80

# Switch to the non-root user
USER 1001

# Set the entrypoint command
ENTRYPOINT ["./OmniView"]
