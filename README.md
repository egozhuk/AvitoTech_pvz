# AvitoTech PVZ

AvitoTech PVZ — микросервис для управления пунктами выдачи заказов (ПВЗ), приёмками и товарами. Реализован на Go, с использованием PostgreSQL, gRPC-ready архитектуры и поддержкой JWT-аутентификации.

⸻

Запуск проекта

1. Клонируй репозиторий
```
git clone https://github.com/your-org/avitotech-pvz.git
cd avitotech-pvz
```
2. Запусти базу и миграции
```
docker-compose up -d db
```
Используется PostgreSQL + Goose (в ./migrations)

3. Запусти сервис

```
go run ./cmd/
```

Сервис поднимется на localhost:8080

⸻

Архитектура

cmd/
└── pvz            # entrypoint

internal/
├── app            # инициализация всех зависимостей
├── controller     # HTTP-обработчики
├── domain         # бизнес-модели
├── repository     # интерфейсы и реализация (postgres)
├── service        # бизнес-логика
├── transport
│   └── middleware # JWT middleware
└── config         # конфигурация через переменные окружения

migrations/
└── *.sql          # goose миграции

⸻

Переменные окружения

Переменная	Значение по умолчанию
PORT	8080
DATABASE_URL	postgres://postgres:postgres@localhost:5432/...
JWT_SECRET	super-secret