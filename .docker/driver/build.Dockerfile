FROM golang:1.21.2-alpine

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY internal/driver ./internal/driver
COPY migrations/driver ./migrations/driver
COPY pkg ./pkg

COPY cmd/driver ./cmd/driver

WORKDIR /app/cmd/driver
RUN go build -o driver

