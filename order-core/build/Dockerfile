FROM golang:1.24.2-alpine AS builder

RUN apk add --no-cache build-base

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main ./cmd/order-core

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main .

CMD ["/app/main"]
