# План Реализации Проекта

Этот план разбит на последовательные фазы и шаги, предназначенные для выполнения ИИ-агентом.

## Фаза 0: Инициализация и настройка окружения

Цель: Подготовить структуру проекта, конфигурацию API Gateway и Docker Compose для всей системы.

**Шаг 0.1: Создание структуры каталогов**
Создать базовую структуру папок для проекта.
```bash
mkdir -p cmd/orders-service
mkdir -p cmd/payments-service
mkdir -p internal/orders/handler
mkdir -p internal/orders/service
mkdir -p internal/orders/storage
mkdir -p internal/orders/models
mkdir -p internal/orders/messaging
mkdir -p internal/payments/handler
mkdir -p internal/payments/service
mkdir -p internal/payments/storage
mkdir -p internal/payments/models
mkdir -p internal/payments/messaging
mkdir -p internal/config
mkdir -p api/openapi
mkdir -p gateway
```

**Шаг 0.2: Конфигурация KrakenD**
Создать файл `gateway/krakend.json`, который будет определять маршрутизацию запросов к микросервисам.
```json
{
  "version": 3,
  "$schema": "https://www.krakend.io/schema/v3.json",
  "port": 8080,
  "endpoints": [
    {
      "endpoint": "/orders", "method": "POST",
      "backend": [{"url_pattern": "/orders", "host": ["http://orders-service:8001"]}]
    },
    {
      "endpoint": "/orders", "method": "GET",
      "backend": [{"url_pattern": "/orders", "host": ["http://orders-service:8001"]}]
    },
    {
      "endpoint": "/orders/{id}", "method": "GET",
      "backend": [{"url_pattern": "/orders/{id}", "host": ["http://orders-service:8001"]}]
    },
    {
      "endpoint": "/payments/account", "method": "POST",
      "backend": [{"url_pattern": "/account", "host": ["http://payments-service:8002"]}]
    },
    {
      "endpoint": "/payments/deposit", "method": "POST",
      "backend": [{"url_pattern": "/deposit", "host": ["http://payments-service:8002"]}]
    },
    {
      "endpoint": "/payments/balance", "method": "GET",
      "backend": [{"url_pattern": "/balance", "host": ["http://payments-service:8002"]}]
    }
  ]
}
```

**Шаг 0.3: Создание `docker-compose.yml`**
Создать основной файл для оркестрации всех компонентов системы.
```yaml
# File: docker-compose.yml
version: '3.8'

services:
  zookeeper:
    image: confluentinc/cp-zookeeper:7.5.0
    container_name: zookeeper
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000

  kafka:
    image: confluentinc/cp-kafka:7.5.0
    container_name: kafka
    depends_on: [zookeeper]
    ports: ["9092:9092"]
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: 'zookeeper:2181'
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1

  db:
    image: postgres:15-alpine
    container_name: db
    ports: ["5431:5432"]
    environment:
        POSTGRES_USER: user
        POSTGRES_PASSWORD: password
        POSTGRES_DB: db
    volumes: ["orders_data:/var/lib/postgresql/data"]

  orders-service:
    build: { context: ., dockerfile: cmd/orders-service/Dockerfile }
    container_name: orders-service
    depends_on: 
        - db
        - kafka
    environment: { LISTEN_ADDR: ":8001", DATABASE_URL: "postgres://user:password@db:5432/orders_db?sslmode=disable", KAFKA_BROKERS: "kafka:9092" }

  payments-service:
    build: { context: ., dockerfile: cmd/payments-service/Dockerfile }
    container_name: payments-service
    depends_on: 
        - db
        - kafka
    environment:
        LISTEN_ADDR: ":8002"
        DATABASE_URL: "postgres://user:password@db:5432/payments_db?sslmode=disable"
        KAFKA_BROKERS: "kafka:9092"

  api-gateway:
    image: devopsfaith/krakend:latest
    container_name: api-gateway
    ports: ["8080:8080"]
    volumes: ["./gateway:/etc/krakend/"]
    depends_on: [orders-service, payments-service]
    command: ["run", "-d", "-c", "/etc/krakend/krakend.json"]

volumes:
  orders_data:
  payments_data:
```

**Шаг 0.4: Создание `.gitignore`**
```
# File: .gitignore
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test
*.out
.DS_Store
vendor/
```

## Фаза 1: Разработка `Payments Service` (Потребитель сообщений)

Цель: Реализовать сервис, отвечающий за счета и обработку платежей. Начинаем с него, так как он является конечным получателем в основном сценарии.

1.  **Создать `cmd/payments-service/Dockerfile`**.
2.  **Реализовать модели:** `internal/payments/models/account.go`.
3.  **Реализовать слой хранения:** `internal/payments/storage/postgres.go`. Он должен включать логику для работы с таблицами `accounts` и `inbox` (для паттерна Transactional Inbox).
4.  **Реализовать слой бизнес-логики:** `internal/payments/service/service.go`.
5.  **Реализовать обработчик Kafka:** `internal/payments/messaging/consumer.go`. Он будет получать сообщения о необходимости оплаты.
6.  **Реализовать обработчик HTTP:** `internal/payments/handler/http.go` для синхронных операций (создание счета, пополнение).
7.  **Собрать все вместе:** `cmd/payments-service/main.go` (инициализация конфига, БД, репозитория, сервиса, запуск Kafka-консьюмера в горутине и HTTP-сервера).

## Фаза 2: Разработка `Orders Service` (Инициатор процесса)

Цель: Реализовать сервис, который создает заказы и инициирует процесс оплаты.

1.  **Создать `cmd/orders-service/Dockerfile`**.
2.  **Реализовать модели:** `internal/orders/models/order.go`.
3.  **Реализовать слой хранения:** `internal/orders/storage/postgres.go`. Включает логику для `orders` и `outbox` (для паттерна Transactional Outbox).
4.  **Реализовать слой бизнес-логики:** `internal/orders/service/service.go`. Метод создания заказа должен в одной транзакции сохранять заказ и сообщение в `outbox`.
5.  **Реализовать Kafka-продюсер/поллер:** `internal/orders/messaging/producer.go`. Отдельная горутина будет опрашивать таблицу `outbox` и отправлять сообщения в Kafka.
6.  **Реализовать обработчик HTTP:** `internal/orders/handler/http.go` для создания и просмотра заказов.
7.  **Собрать все вместе:** `cmd/orders-service/main.go`.

## Фаза 3: Сквозное тестирование

Цель: Проверить работоспособность всего сценария через API Gateway.

1.  **Запустить всю систему:** `docker compose up --build`.
2.  **Выполнить тестовый сценарий с помощью `curl`:**
    -   Создать счет для user_id=1.
    -   Пополнить счет.
    -   Создать заказ для user_id=1.
    -   Проверить статус заказа (сначала NEW, потом FINISHED).
    -   Проверить баланс счета (должен уменьшиться).

## Фаза 4: Реализация WebSocket (для 9-10 баллов)

Цель: Добавить real-time уведомления об изменении статуса заказа.

1.  **В `Orders Service` добавить WebSocket хендлер.**
2.  **При подключении клиента** по WebSocket, сохранять соединение в мапу, где ключ — `order_id`.
3.  **Когда `Orders Service` получает из Kafka сообщение** об успешной оплате и меняет статус заказа, он должен найти соответствующее соединение в мапе и отправить по нему push-уведомление с новым статусом.

## Фаза 5: Финализация

1.  **Написать Unit-тесты**, чтобы достичь покрытия >65%.
2.  **Создать файлы OpenAPI/Swagger** в `api/openapi/` для описания API.
3.  **Написать `README.md`** с описанием проекта и инструкцией по запуску.
