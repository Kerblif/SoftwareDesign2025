FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /orders-service ./cmd/orders-service

FROM alpine:latest

WORKDIR /

COPY --from=builder /orders-service /orders-service

EXPOSE 8001

CMD ["/orders-service"]