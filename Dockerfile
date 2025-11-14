# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Копируем файлы модулей
COPY go.mod ./
RUN go mod download

# Копируем весь код и компилируем
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Run stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарник из builder stage
COPY --from=builder /app/main .

# Копируем статические файлы фронтенда (пока пустую папку)
COPY --from=builder /app/web ./web

EXPOSE 8080

CMD ["./main"]