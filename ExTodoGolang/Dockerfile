FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download 2>/dev/null || go mod tidy

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o server ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates sqlite

WORKDIR /app

COPY --from=builder /app/server .
COPY web/ web/

ENV HTTP_ADDR=:8080 \
    DB_PATH=todo.db \
    ENV=production

EXPOSE 8080

CMD ["./server"]
