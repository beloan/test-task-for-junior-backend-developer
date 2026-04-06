# Quick Start Guide для разработки

## Предварительные требования

- Go 1.23+ (https://golang.org/dl/)
- Docker и Docker Compose
- PostgreSQL 14+ (опционально, если не использовать Docker)
- Git

## Быстрый запуск с Docker

```bash
# 1. Клонируем репозиторий
git clone https://github.com/medods/test-task-for-junior-backend-developer.git
cd test-task-for-junior-backend-developer

# 2. Собираем и запускаем контейнеры
docker compose up --build

# 3. Сервис доступен по адресу:
# API: http://localhost:8080/api/v1
# Swagger: http://localhost:8080/swagger/
```

## Запуск для разработки (без Docker)

```bash
# 1. Устанавливаем зависимости
go mod download

# 2. Запускаем PostgreSQL (например через Docker только БД)
docker run --name postgres-dev -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=taskservice -p 5432:5432 -d postgres:14

# 3. Ждем немного, затем применяем миграции вручную
psql -h localhost -U postgres -d taskservice -f migrations/0001_create_tasks.up.sql
psql -h localhost -U postgres -d taskservice -f migrations/0002_add_recurrence.up.sql

# 4. Запускаем приложение
go run cmd/api/main.go
```

## Структура проекта

```
.
├── cmd/
│   └── api/
│       └── main.go              # Точка входа приложения
├── internal/
│   ├── domain/                  # Бизнес-логика (entities, rules)
│   │   └── task/
│   │       ├── task.go          # Структура Task
│   │       ├── recurrence.go    # Типы и конфигурация периодичности
│   │       ├── recurrence_generator.go  # Генератор расписания
│   │       └── errors.go        # Ошибки домена
│   ├── usecase/                 # Бизнес логика (use cases)
│   │   └── task/
│   │       ├── service.go       # Основной сервис
│   │       ├── ports.go         # Интерфейсы
│   │       └── errors.go        # Ошибки usecase
│   ├── repository/              # Слой доступа к данным
│   │   └── postgres/
│   │       └── task_repository.go
│   ├── infrastructure/          # Инфраструктура (БД, HTTP)
│   │   └── postgres/
│   │       └── pool.go          # Подключение к БД
│   └── transport/               # HTTP слой
│       └── http/
│           ├── handlers/        # HTTP handlers
│           ├── router.go        # Маршруты
│           └── docs/            # Swagger документация
├── migrations/                  # SQL миграции
├── docker-compose.yml           # Docker конфигурация
├── Dockerfile                   # Docker образ
├── go.mod, go.sum              # Go модули
├── README.md                    # Документация
└── IMPLEMENTATION_NOTES.md      # Подробности реализации
```

## Разработка

### Добавление новых типов периодичности

1. Добавить константу в `domain/task/recurrence.go`:
```go
const RecurrenceTypeCustom RecurrenceType = "custom"
```

2. Добавить валидацию в `RecurrenceConfig.Validate()`:
```go
case RecurrenceTypeCustom:
    // проверяем валидность параметров
```

3. Добавить метод генерирования в `RecurrenceGenerator`:
```go
func (rg *RecurrenceGenerator) generateCustom(...) []time.Time {
    // реализация
}
```

4. Обновить switch в `GenerateNextOccurrences()`:
```go
case RecurrenceTypeCustom:
    occurrences = rg.generateCustom(...)
```

###添加новых endpointов

1. Добавить метод в `TaskHandler` (`internal/transport/http/handlers/task_handler.go`)
2. Добавить маршрут в `router.go`
3. Обновить Swagger документацию

### Тестирование

```bash
# Запуск всех тестов
go test ./...

# Запуск с verbose output
go test -v ./...

# Запуск тестов конкретного пакета
go test ./internal/domain/task/...

# Запуск с покрытием
go test -cover ./...
```

## Полезные команды

```bash
# Форматирование кода
go fmt ./...

# Проверка ошибок
go vet ./...

# Линтинг (требует golangci-lint)
golangci-lint run

# Просмотр документации локально
go doc -http=:6060

# Генерирование Swagger документации (если используется)
swag init -g cmd/api/main.go
```

## Переменные окружения

```bash
# HTTP адрес и порт
HTTP_ADDR=:8080

# PostgreSQL DSN
DATABASE_DSN=postgres://postgres:postgres@localhost:5432/taskservice?sslmode=disable
```

## Отладка и логирование

Приложение использует `log/slog` для логирования. При необходимости можно обновить:
- Уровень логирования (INFO, DEBUG, ERROR)
- Формат вывода (TextHandler, JSONHandler)

## Проблемы и решения

### 1. "database is locked"
- Убедитесь что Docker контейнер PostgreSQL запущен
- Попробуйте пересоздать volume: `docker compose down -v`

### 2. "connection refused"
- Проверьте что PostgreSQL доступна на правильном адресе/порте
- Дождитесь полной инициализации контейнера (~5 сек)

### 3. Ошибки миграции
- Проверьте что файлы миграций в правильном месте
- Убедитесь что пермиссии на файлы правильные

## Примеры использования API

Смотрите `test_api.sh` для примеров curl команд.

```bash
# Разрешить выполнение скрипта
chmod +x test_api.sh

# Запустить примеры (после того как сервис запущен)
./test_api.sh
```

## Дополнительные ресурсы

- [Go documentation](https://golang.org/doc/)
- [PostgreSQL documentation](https://www.postgresql.org/docs/)
- [Docker documentation](https://docs.docker.com/)
- [Clean Architecture in Go](https://github.com/golang-standards/project-layout)

## Контроль качества кода

- Все публичные функции должны иметь документирующие комментарии
- Переменные должны иметь понятные имена
- Обработка ошибок обязательна
- Тесты для критичной бизнес-логики

## Развертывание в production

```bash
# Сборка образа
docker build -t task-service:latest .

# Запуск в production режиме
docker run -e DATABASE_DSN="..." -p 8080:8080 task-service:latest

# С использованием docker-compose
docker compose -f docker-compose.yml up -d
```
