FROM golang:1.21.2-alpine

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY internal/location ./internal/location
COPY migrations/location ./migrations/location

COPY cmd/location ./cmd/location

WORKDIR /app/cmd/location
RUN go build -o app

