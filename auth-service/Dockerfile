# Используйте официальный образ Go в качестве базового
FROM golang:1.20 AS builder

# Установите рабочую директорию
WORKDIR /microservice

# Копируйте go.mod и go.sum для загрузки зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируйте весь исходный код в контейнер
COPY . .

# Компилируйте приложение для Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o employee-service ./cmd/main.go

# Используйте минимальный образ для запуска
FROM alpine:3.18

WORKDIR /microservice

# Обновите репозитории и установите зависимости
RUN apk update && apk --no-cache add ca-certificates

# Копируйте скомпилированное приложение из предыдущего этапа
COPY --from=builder /microservice/catalog-service .

# Установите права на выполнение
RUN chmod +x ./employee-service

# Укажите команду для запуска приложения
CMD ["./catalog-service"]
