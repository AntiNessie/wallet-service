# Wallet Service

Сервис управления балансом кошелька (тестовое задание)

## Стек

- Go 1.22
- PostgreSQL 15
- Docker 19.03+ (использовался Docker Toolbox)
- docker-compose 2.2 (ранняя версия из-за проблем с гипервизором и WSL2 на моей машине)

## Запуск

docker-compose up --build


Сервис будет доступен на http://localhost:8080.
При использовании Docker Toolbox — на http://192.168.99.100:8080.

## API

### Пополнение счёта

POST /api/v1/wallet
Content-Type: application/json

{"valletId":"550e8400-e29b-41d4-a716-446655440000","operationType":"DEPOSIT","amount":1000}

### Уменьшение счёта:

POST /api/v1/wallet
Content-Type: application/json

{"valletId":"550e8400-e29b-41d4-a716-446655440000","operationType":"WITHDRAW","amount":300}

### Просмотр баланса

GET /api/v1/wallets/550e8400-e29b-41d4-a716-446655440000


## Тесты

go test ./internal/service/

## Структура проекта

cmd/main.go - точка входа
internal/config/config.go - конфигурация из окружения
internal/handler/wallet.go - HTTP обработчики
internal/logger/logger.go - slog
internal/model/wallet.go - модели данных
internal/repository/ - слой работы с БД
internal/service/ - бизнес-логика
migrations/ - SQL миграции
config.env - переменные окружения
docker-compose.yml - docker-compose 2.2
Dockerfile - multi-stage сборка

## Конкурентность

Для обеспечения атомарности при 1000 RPS на один кошелёк используется UPDATE wallets SET balance = balance + $1 — баланс изменяется на стороне PostgreSQL без read-modify-write.
Это гарантирует что ни один запрос не будет потерян или обработан некорректно.

