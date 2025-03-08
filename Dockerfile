# syntax=docker/dockerfile:1
# A basic app to experiment with go
FROM golang:1.23-alpine AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Add build arguments for version and build date
ARG VERSION=0.1.0
ARG BUILD_DATE

# Set environment variables
ENV APP_VERSION=${VERSION}
ENV BUILD_DATE=${BUILD_DATE}

RUN CGO_ENABLED=0 GOOS=linux go build -o /server

EXPOSE 8080

CMD ["/server"]