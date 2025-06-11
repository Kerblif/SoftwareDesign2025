FROM golang:1.24-alpine AS builder

RUN apk add --no-cache build-base

RUN go install fyne.io/tools/cmd/fyne@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN /go/bin/fyne package -os web -name website -icon hse.png --src ./cmd/client-app/

# Stage 2: Serve the application with Nginx
# This stage uses a lightweight Nginx image to serve the static content.
FROM nginx:alpine

# Copy the built website from the builder stage to the Nginx public HTML directory.
COPY --from=builder /app/wasm /usr/share/nginx/html

# Expose port 80 to allow traffic to the web server.
EXPOSE 80

# The default command for the nginx image is to start the server.
# We add '-g "daemon off;"' to run Nginx in the foreground, which is best practice for containers.
CMD ["nginx", "-g", "daemon off;"]
