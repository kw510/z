# Use an official Golang runtime as a parent image
FROM --platform=$BUILDPLATFORM golang:1-alpine AS builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

# Set the working directory to /app
WORKDIR /app

# Download the go dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

# Build a static application binary
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o ./tmp/main

## Development stage, using air for hot reloading
FROM builder AS development
RUN go install github.com/air-verse/air@latest
CMD ["air", "-c", ".air.toml"]

## Production stage, using a static binary and scratch image
FROM scratch AS production
COPY --from=builder /app/tmp/main /app
CMD ["/app"]