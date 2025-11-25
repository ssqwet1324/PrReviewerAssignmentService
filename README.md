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

### 2. Установите Docker и Docker Compose

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

## Бенчмарки API

| Эндпоинт                                    | Метод | Описание           | Средняя задержка | RPS      | Передача данных |
| ------------------------------------------- | ----- | ------------------ | ---------------- | -------- | --------------- |
| http://localhost:8080/team/add              | POST  | Добавление команды | 94.11ms          | 10000    | 1.94MB/s        |
| http://localhost:8080/team/get              | GET   | Получение команды  | 149.76ms         | 10000    | 4.34MB/s        |

### Подробные результаты

#### POST http://localhost:8080/team/add

``` wrk2 -t4 -c100 -d30s -R10000 -s ./postCreate "http://localhost:8080/team/add"
  4 threads and 100 connections
  Thread calibration: mean lat.: 1.573ms, rate sampling interval: 10ms
  Thread calibration: mean lat.: 1.348ms, rate sampling interval: 10ms
  Thread calibration: mean lat.: 1.565ms, rate sampling interval: 10ms
  Thread calibration: mean lat.: 1.567ms, rate sampling interval: 10ms
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    94.11ms  295.25ms   1.77s    90.87%
    Req/Sec     2.90k     1.97k   34.30k    97.06%
  296996 requests in 30.00s, 58.35MB read
  Non-2xx or 3xx responses: 296996
Requests/sec:   9899.67
Transfer/sec:      1.94MB
```

#### GET http://localhost:8080/team/get

``` wrk2 -t4 -c100 -d30s -R10000 "http://localhost:8080/team/get?team_name=juniors"

Running 30s test @ http://localhost:8080/team/get?team_name=juniors
  4 threads and 100 connections
  Thread calibration: mean lat.: 1.909ms, rate sampling interval: 10ms
  Thread calibration: mean lat.: 1.947ms, rate sampling interval: 10ms
  Thread calibration: mean lat.: 1.911ms, rate sampling interval: 10ms
  Thread calibration: mean lat.: 1.907ms, rate sampling interval: 10ms
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   149.76ms  397.05ms   1.93s    88.45%
    Req/Sec     2.93k     1.28k   11.44k    90.98%
  299388 requests in 30.00s, 109.07MB read
Requests/sec:   9979.13
Transfer/sec:      3.64MB
```

### Примеры скриптов для wrk

#### postCreate.lua (POST)

```
wrk.method = "POST"
wrk.body   = '{"pull_request_id":"pr-1001","pull_request_name":"Add search","author_id":"u1"}'
wrk.headers["Content-Type"] = "application/json"
```

---

| Эндпоинт                                       | Метод | Описание                       | Средняя задержка | RPS     | Передача данных |
| ---------------------------------------------- | ----- | ------------------------------ | ---------------- | ------- | --------------- |
| http://localhost:8080/users/setIsActive        | POST  | Изменение статуса пользователя | 17.59s           | 10000   | 3184.37KB/s     |
| http://localhost:8080/users/getReview          | GET   | Получение ревью пользователя   | 195.02ms         | 10000   | 1.50MB/s        |


### Подробные результаты

#### POST http://localhost:8080/users/setIsActive

``` wrk2 -t4 -c100 -d30s -R10000 -s ./postChange "http://localhost:8080/users/setIsActive"

  Running 30s test @ http://localhost:8080/users/setIsActive
  4 threads and 100 connections
  Thread calibration: mean lat.: 4271.807ms, rate sampling interval: 15548ms
  Thread calibration: mean lat.: 4253.044ms, rate sampling interval: 15515ms
  Thread calibration: mean lat.: 4263.256ms, rate sampling interval: 15605ms
  Thread calibration: mean lat.: 4264.301ms, rate sampling interval: 15556ms
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    17.59s     5.46s   26.39s    55.73%
    Req/Sec   269.75      0.83   271.00    100.00%
  36083 requests in 30.01s, 5.40MB read
  Socket errors: connect 0, read 0, write 0, timeout 1
Requests/sec:   1202.53
Transfer/sec:    184.37KB
```

#### GET http://localhost:8080/users/getReview

```  wrk2 -t4 -c100 -d30s -R10000 "http://localhost:8080/users/getReview?user_id=u1"

Running 30s test @ http://localhost:8080/users/getReview?user_id=u1
  4 threads and 100 connections
  Thread calibration: mean lat.: 1.661ms, rate sampling interval: 10ms
  Thread calibration: mean lat.: 2.035ms, rate sampling interval: 10ms
  Thread calibration: mean lat.: 1.624ms, rate sampling interval: 10ms
  Thread calibration: mean lat.: 2.041ms, rate sampling interval: 10ms
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   195.02ms  476.01ms   2.07s    87.12%
    Req/Sec     2.95k     1.26k    9.89k    86.78%
  294448 requests in 30.00s, 44.93MB read
Requests/sec:   9814.01
Transfer/sec:      1.50MB
```

### Примеры скриптов для wrk

#### postChange.lua (POST)

```
wrk.method = "POST"
wrk.body   = '{"pull_request_id":"pr-1001","pull_request_name":"Add search","author_id":"u1"}'
wrk.headers["Content-Type"] = "application/json"
```

---


## Проблемы с которыми столкнулся

1. Проблема подключения Goose к pgxPool
Проблема решилась тем, что нужно было создать pgx.ConnConfig и далее воспользоваться sql.DB.

2. Проблема создания самописных ошибок.
Я обычно обрабатывал ошибки через ctx.Json("тут http код ошибки", тут просто какая будет ошибка от кода)
Тут чтобы решить проблему пришлось создавать собтвенные ошибки и сравнивать их с теми, которые приходят от бизнес логики.
