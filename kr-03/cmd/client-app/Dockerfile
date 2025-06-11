FROM golang:1.20 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o client-app ./cmd/client-app

FROM ubuntu:22.04
RUN apt-get update && apt-get install -y \
    xorg-dev \
    libgl1-mesa-dev \
    libegl1-mesa-dev \
    libgles2-mesa-dev \
    && rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/client-app /usr/local/bin/
EXPOSE 8080
CMD ["client-app", "-webgl"]