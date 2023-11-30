FROM golang:1.21.2-alpine

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY internal/driver ./internal/driver
COPY cmd/driver ./cmd/driver
COPY migrations/driver ./migrations/driver
COPY pkg/driver ./pkg/driver

WORKDIR /app/cmd/driver
RUN go build -o app

