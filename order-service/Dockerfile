# Используйте официальный образ Go в качестве базового
FROM golang:1.20 AS builder

# Установите рабочую директорию
WORKDIR /microservice

# Копируйте go.mod и go.sum для загрузки зависимостей
COPY go.mod ./
RUN go mod download

# Копируйте весь исходный код в контейнер
COPY . .

# Компилируйте приложение для Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o order-service ./cmd/main.go

# Используйте минимальный образ для запуска
FROM alpine:latest

WORKDIR /microservice

# Копируйте скомпилированное приложение из предыдущего этапа
COPY --from=builder /microservice/order-service .

# Установите права на выполнение
RUN chmod +x ./order-service

# Укажите команду для запуска приложения
CMD ["./order-service"]
