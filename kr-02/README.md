# Микросервисная архитектура анализатора текста

Этот проект реализует микросервисную архитектуру для анализа текстовых файлов, включая расчет статистики, обнаружение плагиата и генерацию облака слов.

## Архитектура

Система состоит из трех микросервисов:

1. **API Gateway** - Отвечает за маршрутизацию запросов к соответствующим сервисам
2. **File Storing Service** - Отвечает за хранение и получение файлов
3. **File Analysis Service** - Отвечает за анализ файлов и хранение результатов

## Предварительные требования

- Docker и Docker Compose
- Go 1.20 или новее
- Swag (для генерации документации Swagger)
  ```bash
  go install github.com/swaggo/swag/cmd/swag@latest
  ```

## Начало работы

### Клонирование репозитория

```bash
git clone <repository-url>
cd kr-02
```

### Генерация документации Swagger

```bash
make swagger
```

### Сборка и запуск сервисов

```bash
docker-compose up -d
```

Это запустит все сервисы:
- База данных PostgreSQL на порту 5432
- File Storing Service на порту 50051
- File Analysis Service на порту 50052
- API Gateway на порту 8080 (HTTP)
- Swagger UI на порту 8081

## Документация API

### Генерация документации Swagger

```bash
make swagger
```

Это сгенерирует документацию Swagger с использованием swag на основе аннотаций в коде.

### Доступ к Swagger UI

Затем откройте http://localhost:8081 в вашем браузере для просмотра документации API.

Вы также можете получить доступ к Swagger UI напрямую из API Gateway по адресу http://localhost:8080/swagger/index.html.

## Конечные точки API

### Загрузка файла

```
POST /api/v1/files
```

Запрос: multipart/form-data с полем файла с именем "file"

Пример использования curl:
```bash
curl -X POST -F "file=@example.txt" http://localhost:8080/api/v1/files
```

Ответ:
```json
{
  "file_id": "unique-file-id"
}
```

### Получение файла

```
GET /api/v1/files/{file_id}
```

Ответ: Бинарное содержимое файла с соответствующим заголовком Content-Disposition для скачивания

Пример использования curl:
```bash
curl -OJ http://localhost:8080/api/v1/files/{file_id}
```

### Анализ файла

```
POST /api/v1/analysis
```

Тело запроса:
```json
{
  "file_id": "unique-file-id",
  "generate_word_cloud": true
}
```

Ответ:
```json
{
  "paragraph_count": 5,
  "word_count": 100,
  "character_count": 500,
  "is_plagiarism": false,
  "similar_file_ids": [],
  "word_cloud_location": "word-cloud-location"
}
```

### Получение облака слов

```
GET /api/v1/wordcloud/{location}
```

Ответ: Изображение облака слов (бинарные данные) с Content-Type: image/png

Пример использования curl:
```bash
curl -o wordcloud.png http://localhost:8080/api/v1/wordcloud/{location}
```

## Разработка

### Структура проекта

```
kr-02/
├── build/                    # Dockerfiles
│   ├── api_gateway/
│   ├── file_analysis_service/
│   └── file_storing_service/
├── cmd/                      # Точки входа для каждого сервиса
│   ├── api_gateway/
│   ├── file_analysis_service/
│   └── file_storing_service/
├── configs/                  # Конфигурационные файлы
├── internal/                 # Внутренние пакеты
│   ├── pkg/                  # Общие пакеты
│   │   ├── api_gateway/      # Реализация API Gateway
│   │   │   ├── clients/      # Клиенты сервисов
│   │   │   ├── docs/         # Сгенерированная документация Swagger
│   │   │   └── handlers/     # HTTP обработчики
│   │   ├── file_analysis/    # Реализация File Analysis Service
│   │   └── file_storing/     # Реализация File Storing Service
│   └── proto/                # Сгенерированные proto файлы для gRPC сервисов
├── proto/                    # Proto определения для gRPC сервисов
├── scripts/                  # Вспомогательные скрипты
└── tests/                    # Тесты
```

### Запуск тестов

```bash
go test ./tests/...
```

### Запуск тестов с покрытием

Для запуска всех тестов и генерации отчета о покрытии кода:

```bash
make test-coverage
```

Эта команда:
1. Запускает все тесты в проекте
2. Генерирует отчет о покрытии кода в файле `coverage/coverage.out`
3. Создает HTML-версию отчета в файле `coverage/coverage.html`
4. Выводит в консоль сводку покрытия по функциям

Для просмотра HTML-отчета откройте файл `coverage/coverage.html` в браузере.

## Лицензия

Этот проект лицензирован под лицензией MIT - см. файл LICENSE для получения подробной информации.
