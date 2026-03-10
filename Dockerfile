FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOD= go build -o bot .

FROM alpine:latest
WORKDIR /Users/oli/Code/Cat-Bot
COPY --from=builder /app/bot .
CMD ["./bot"]


