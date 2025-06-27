FROM golang:1.24.2-alpine AS builder

RUN apk add --no-cache build-base

WORKDIR /app

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd/auth-service

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main .

COPY .env .

CMD ["/app/main"]
