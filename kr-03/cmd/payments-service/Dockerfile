FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /payments-service ./cmd/payments-service

FROM alpine:latest

WORKDIR /

COPY --from=builder /payments-service /payments-service

EXPOSE 8002

CMD ["/payments-service"]