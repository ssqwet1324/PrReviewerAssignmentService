# PR Reviewer Assignment Service

Сервис для автоматического назначения ревьюверов на Pull Request'ы.

## Технологии

- **Go 1.25.3** - язык программирования
- **Gin** - веб-фреймворк для HTTP API
- **PostgreSQL** - база данных
- **pgx/v5** - драйвер для PostgreSQL
- **Zap** - структурированное логирование
- **Goose** - миграции базы данных
- **Docker & Docker Compose** - контейнеризация

## Установка зависимостей

### 1. Установите Go

Скачайте и установите Go версии 1.25.3 или выше:
- [Официальный сайт Go](https://golang.org/dl/)
- Или через пакетный менеджер вашей ОС

Проверьте установку:
```bash
go version
```

### 2. Установите Docker и Docker Compose

- **Docker**: [Скачать Docker](https://www.docker.com/get-started)
- **Docker Compose**: обычно устанавливается вместе с Docker

Проверьте установку:
```bash
docker --version
docker-compose --version
```

### 3. Установите Make (опционально)

- **Windows**: установите через [Chocolatey](https://chocolatey.org/) или используйте Git Bash
- **Linux/Mac**: обычно уже установлен

## Сборка и запуск

### Быстрый старт с Docker

1. Клонируйте репозиторий:
```bash
git clone <repository-url>
cd PrReviewerAssignmentService
```

2. Запустите проект:
```bash
make up
```

Или без Make:
```bash
docker-compose up -d
```

3. Сервис будет доступен по адресу: `http://localhost:8080`

4. Остановить проект:
```bash
make down
```

Или:
```bash
docker-compose down
```

### Локальный запуск (без Docker)

1. Установите зависимости:
```bash
go mod download
```

2. Настройте переменные окружения (создайте файл `.env`):
```env
DB_NAME=postgres
DB_USER=postgres
DB_PASSWORD=postgres
DB_HOST=localhost
DB_PORT=5432
```

3. Запустите PostgreSQL локально или используйте Docker только для БД:
```bash
docker-compose up -d postgres
```

4. Запустите приложение:
```bash
go run cmd/pr_reviewer_service/main.go
```

## Полезные команды

- `make up` - запустить проект в Docker
- `make down` - остановить проект
- `make restart` - перезапустить проект
- `make logs` - показать логи
- `make build` - пересобрать Docker образ
- `make rebuild` - пересобрать и запустить

## API Endpoints

- `POST /team/add` - создать команду
- `GET /team/get?team_name=<name>` - получить команду
- `POST /users/setIsActive` - изменить активность пользователя
- `GET /users/getReview?user_id=<id>` - получить PR'ы пользователя
- `POST /pullRequest/create` - создать PR
- `POST /pullRequest/merge` - замержить PR
- `POST /pullRequest/reassign` - переназначить ревьювера

## Структура проекта

```
.
├── cmd/pr_reviewer_service/  # Точка входа
├── internal/
│   ├── app/                  # Инициализация приложения
│   ├── config/               # Конфигурация
│   ├── entity/               # Модели данных
│   ├── handler/              # HTTP handlers
│   ├── repository/           # Работа с БД
│   └── usecase/              # Бизнес-логика
├── migrations/               # Миграции БД
├── docker-compose.yml        # Docker конфигурация
├── Dockerfile                # Образ приложения
└── Makefile                  # Команды для работы с проектом
```

## Проблемы с которыми столкнулся

1. Проблема подключения Goose к pgxPool
Проблема решилась тем, что нужно было создать pgx.ConnConfig и далее воспользоваться sql.DB.

2. Проблема создания самописных ошибок.
Я обычно обрабатывал ошибки через ctx.Json("тут http код ошибки", тут просто какая будет ошибка от кода)
Тут чтобы решить проблему пришлось создавать собтвенные ошибки и сравнивать их с теми, которые приходят от бизнес логики.