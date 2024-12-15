# Сервис магазина продуктов

Этот репозиторий содержит прототип **Сервиса магазина продуктов**. Сервис предназначен для управления информацией о товарах, поддерживает интерфейсы REST и SOAP для операций CRUD. Реализован как архитектура микросервисов и использует Docker Compose для развертывания.

## Возможности
- **REST API** и **SOAP API** для выполнения операций CRUD с товарами.
- Поддержка добавления, получения, обновления и удаления товаров.
- Взаимодействие с базой данных PostgreSQL.
- Веб-клиент для удобного взаимодействия с API.
- Docker-контейнеры для упрощённой настройки и масштабирования.

---

## Структура проекта
```
.
├── catalog-service       # REST-сервис для управления товарами
│   ├── main.go           # Исходный код приложения на Go
│   └── Dockerfile        # Dockerfile для catalog-service
├── soap-service          # SOAP-сервис для управления товарами
│   ├── main.go           # SOAP-реализация на Go
│   └── Dockerfile        # Dockerfile для soap-service
├── proxy-service         # сервис для управления запросами и перенаправления в soap-service для управления товарами
│   ├── main.go           # SOAP-реализация на Go
│   └── Dockerfile        # Dockerfile для soap-service
├── web-client            # Простой веб-клиент на HTML
│   └── client.html       # Веб-клиент для взаимодействия с API
├── docker-compose.yml    # Docker Compose файл для оркестрации сервисов
└── README.md             # Этот файл
```

---

## Описание сервисов
### 1. Catalog Service (REST API)
Catalog Service предоставляет REST API для работы с товарами.
- Базовый URL: `http://localhost:8081`
- Эндпоинты:
  - `GET /products` - Получить список всех товаров.
  - `GET /products/{id}` - Получить товар по ID.
  - `POST /products` - Добавить новый товар.
  - `PUT /products/{id}` - Обновить существующий товар.
  - `DELETE /products/{id}` - Удалить товар.

### 2. SOAP Service
SOAP Service предоставляет SOAP API для работы с товарами.
- Базовый URL: `http://localhost:8082`
- Поддерживаемые действия:
  - `add` - Добавить новый товар.
  - `get` - Получить товар по ID.
  - `getAll` - Получить список всех товаров.
  - `update` - Обновить существующий товар.
  - `delete` - Удалить товар.

### 3. Web Client
Web Client — это простой интерфейс для взаимодействия с REST и SOAP API.
- Доступен по адресу: `http://localhost:8080`
- Поддерживает:
  - Добавление, получение, обновление и удаление товаров.
  - Переключение между REST и SOAP API.

---

## Предназначение сервисов
Этот прототип демонстрирует базовые возможности управления товарами для продуктового магазина. Такие сервисы могут быть использованы в:
- Интернет-магазинах для работы с каталогом товаров.
- Приложениях управления запасами в магазинах.
- Внедрении интеграции между внутренними системами учёта товаров.

---

## Требования
- Установленные **Docker** и **Docker Compose**.

---

## Запуск проекта
1. Клонируйте репозиторий:
   ```bash
   git clone <repository_url>
   cd <repository_folder>
   ```

2. Постройте и запустите сервисы с помощью Docker Compose:
   ```bash
   docker-compose up --build
   ```

3. Получите доступ к сервисам:
   - REST API: `http://localhost:8081`
   - SOAP API: `http://localhost:8082`
   - Веб-клиент: `http://localhost:8080`

---

## Переменные окружения
Сервисы используют следующие переменные окружения:
- `DB_HOST`: Хост базы данных (по умолчанию: `catalog-db`).
- `DB_PORT`: Порт базы данных (по умолчанию: `5432`).
- `DB_USER`: Пользователь базы данных (по умолчанию: `postgres`).
- `DB_PASSWORD`: Пароль базы данных (по умолчанию: `secret`).
- `DB_NAME`: Имя базы данных (по умолчанию: `catalog`).

---

## Примеры использования
### Пример REST API
#### Добавление товара:
```bash
curl -X POST -H "Content-Type: application/json"   -d '{"name": "Product 1", "price": 100.0}'   http://localhost:8081/products
```

#### Получение списка всех товаров:
```bash
curl http://localhost:8081/products
```

### Пример SOAP API
#### Добавление товара:
```xml
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
    <soapenv:Body>
        <ProductOperation>
            <action>add</action>
            <product>
                <name>Product 1</name>
                <price>100.0</price>
            </product>
        </ProductOperation>
    </soapenv:Body>
</soapenv:Envelope>
```

#### Получение списка всех товаров:
```xml
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
    <soapenv:Body>
        <ProductOperation>
            <action>getAll</action>
        </ProductOperation>
    </soapenv:Body>
</soapenv:Envelope>
```

---

## Возможности для улучшения
- Добавление аутентификации и авторизации.
- Реализация пагинации для больших объёмов данных.
- Улучшение обработки ошибок.
- Создание более интерактивного и удобного веб-клиента.

---

## Лицензия
Этот проект лицензируется на условиях MIT License.

---

## Авторы
- Смольский Валерий - Разработчик
