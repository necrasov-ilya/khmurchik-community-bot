# khmurchik-community-bot — Техническая документация

## Оглавление

1. [Обзор проекта](#обзор-проекта)
2. [Архитектура](#архитектура)
3. [Структура проекта](#структура-проекта)
4. [Технологический стек](#технологический-стек)
5. [Конфигурация](#конфигурация)
6. [База данных](#база-данных)
7. [Модульная система](#модульная-система)
8. [Активные модули](#активные-модули)
9. [Команды бота](#команды-бота)
10. [Жизненный цикл приложения](#жизненный-цикл-приложения)
11. [Развёртывание](#развёртывание)
12. [Разработка](#разработка)

---

## Обзор проекта

**khmurchik-community-bot** — это многомодульный Telegram-бот на Go, разработанный для сообщества Павла Хмурчика (техничный чат). Бот обеспечивает ежедневные приветственные сообщения и инструменты модерации чата.

### Ключевые возможности

- Ежедневное приветственное сообщение в заданный чат (по расписанию cron)
- Модерация: `/mute`, `/ban`, `/kick`, `/unmute`
- Персистентное хранение данных пользователей и чатов в PostgreSQL
- Логирование действий модерации
- Проверка прав администратора через Telegram API
- Конфигурация через YAML + environment variables

---

## Архитектура

Проект следует чистой модульной архитектуре с разделением на слои:

```
┌─────────────────────────────────────────────────────────┐
│                    cmd/bot/main.go                      │
│                   (Точка входа)                         │
├─────────────────────────────────────────────────────────┤
│                   internal/app/app.go                   │
│              (Оркестрация жизненного цикла)              │
├──────────┬──────────┬───────────┬───────────────────────┤
│  bot/    │  modules │  di/      │  modulesystem/        │
│  (роутер)│  (модули)│  (DI)     │  (интерфейс модулей)  │
├──────────┴──────────┴───────────┴───────────────────────┤
│  handlers/  middleware/  repository/  logger/  config/   │
│  (хендлеры) (проверки)  (данные)     (логи)  (настройки)│
├─────────────────────────────────────────────────────────┤
│                   internal/database/                    │
│                 (pgxpool → PostgreSQL)                   │
└─────────────────────────────────────────────────────────┘
```

### Принцип работы

1. `main.go` загружает конфигурацию, создаёт логгер и пул соединений БД
2. Создается `App` — центральный оркестратор
3. Регистрируются модули (greeting, moderation)
4. Запускается polling-цикл Telegram API
5. Каждое сообщение проходит через router → handler → repository
6. При SIGINT/SIGTERM — graceful shutdown

---

## Структура проекта

```
khmurchik-community-bot/
├── cmd/bot/main.go                  # Точка входа приложения
├── internal/
│   ├── app/app.go                   # Жизненный цикл приложения
│   ├── bot/router.go                # Роутер команд Telegram-бота
│   ├── config/
│   │   ├── config.go                # Структуры конфигурации
│   │   └── loader.go                # Загрузка конфигурации (Viper)
│   ├── database/postgres.go         # Создание pgxpool
│   ├── di/container.go              # DI-контейнер (запланирован)
│   ├── handlers/
│   │   ├── start.go                 # Хендлер /start
│   │   └── unknown.go               # Фолбэк для неизвестных команд
│   ├── logger/logger.go             # Zap logger (production mode)
│   ├── middleware/admin_check.go    # Проверка прав администратора
│   ├── modules/
│   │   ├── greeting/                # Модуль приветствий
│   │   └── moderation/              # Модуль модерации
│   ├── modulesystem/
│   │   ├── module.go                # Интерфейс модуля
│   │   └── registry.go              # Реестр модулей
│   ├── repository/
│   │   ├── users.go                 # UserRepository
│   │   └── chats.go                 # ChatRepository
│   └── timeutil/
│       ├── duration.go              # Парсер длительности (s/m/h/d)
│       └── timezone.go              # Загрузчик таймзон
├── configs/config.example.yaml      # Пример конфигурации
├── migrations/                      # SQL-миграции
│   ├── 0001_create_users.*.sql
│   ├── 0002_create_chats.*.sql
│   └── 0003_create_moderation_logs.*.sql
├── scripts/
│   ├── migrate_up.sh                # Применение миграций
│   └── migrate_down.sh              # Откат миграций
├── Dockerfile                       # Multi-stage build
├── docker-compose.yml               # Docker Compose (bot + postgres)
├── Makefile                         # Build/run/lint/migrate
├── go.mod / go.sum                  # Go модуль
├── .golangci.yml                    # Настройки линтера
└── .env.example                     # Шаблон переменных окружения
```

---

## Технологический стек

| Компонент | Библиотека | Версия | Назначение |
|-----------|-----------|--------|-----------|
| Язык | Go | 1.25.0 | Основной язык |
| Telegram API | `go-telegram-bot-api/v5` | v5.5.1 | Клиент Telegram Bot API |
| PostgreSQL | `jackc/pgx/v5` | v5.10.0 | Драйвер + connection pool (pgxpool) |
| Конфигурация | `spf13/viper` | v1.21.0 | YAML + env vars |
| Валидация | `go-playground/validator/v10` | v10.30.3 | Валидация структур |
| Cron | `robfig/cron/v3` | v3.0.1 | Расписание задач |
| Логирование | `uber-go/zap` | v1.28.0 | Production-логгер |
| Миграции | `golang-migrate/migrate` | — | CLI для миграций БД |
| Линтер | `golangci-lint` | — | Статический анализ |

### Прямые зависимости

```
github.com/go-playground/validator/v10    v10.30.3
github.com/go-telegram-bot-api/telegram-bot-api/v5  v5.5.1
github.com/jackc/pgx/v5                   v5.10.0
github.com/robfig/cron/v3                 v3.0.1
github.com/spf13/viper                    v1.21.0
go.uber.org/zap                           v1.28.0
```

---

## Конфигурация

### Иерархия конфигурации

```
Config
├── BotConfig
│   ├── Token        string  (validate: "required")
│   └── TargetChatID int64   (validate: "required")
├── ServerConfig
│   └── PollTimeout  string  (validate: "required")
├── DatabaseConfig
│   ├── Host         string  (default: "localhost")
│   ├── Port         int     (default: 5432)
│   ├── User         string  (default: "postgres")
│   ├── Password     string  (default: "postgres")
│   ├── Name         string  (default: "khmurchik_bot")
│   ├── SSLMode      string  (default: "disable")
│   ├── MaxOpenConns int     (default: 25)
│   ├── MaxIdleConns int     (default: 10)
│   └── ConnMaxLifetime string (default: "5m")
└── SchedulerConfig
    └── Timezone     string  (default: "Europe/Minsk")
```

### Загрузка конфигурации

1. Viper читает `configs/config.yaml` (YAML)
2. Автоматически связывает переменные окружения с префиксом `BOT_`
3. Ключи YAML `bot.token` → `BOT_BOT_TOKEN` (`.` заменяется на `_`)
4. Валидация через `go-playground/validator`

### Mapping env vars

| YAML key | Environment variable |
|----------|---------------------|
| `bot.token` | `BOT_BOT_TOKEN` |
| `bot.target_chat_id` | `BOT_BOT_TARGET_CHAT_ID` |
| `server.poll_timeout` | `BOT_SERVER_POLL_TIMEOUT` |
| `database.host` | `BOT_DATABASE_HOST` |
| `database.port` | `BOT_DATABASE_PORT` |
| `database.user` | `BOT_DATABASE_USER` |
| `database.password` | `BOT_DATABASE_PASSWORD` |
| `database.name` | `BOT_DATABASE_NAME` |
| `database.ssl_mode` | `BOT_DATABASE_SSL_MODE` |
| `database.max_open_conns` | `BOT_DATABASE_MAX_OPEN_CONNS` |
| `database.max_idle_conns` | `BOT_DATABASE_MAX_IDLE_CONNS` |
| `database.conn_max_lifetime` | `BOT_DATABASE_CONN_MAX_LIFETIME` |
| `scheduler.timezone` | `BOT_SCHEDULER_TIMEZONE` |

### Пример конфигурации (`configs/config.example.yaml`)

```yaml
bot:
  token: ""
  target_chat_id: 0

server:
  poll_timeout: 30

database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "postgres"
  name: "khmurchik_bot"
  ssl_mode: "disable"
  max_open_conns: 25
  max_idle_conns: 10
  conn_max_lifetime: "5m"

scheduler:
  timezone: "Europe/Minsk"
```

---

## База данных

### Схема

#### Таблица `users`

| Колонка | Тип | Описание |
|---------|-----|----------|
| `telegram_id` | `BIGINT` | Primary key, ID пользователя Telegram |
| `username` | `TEXT` | Имя пользователя Telegram |
| `first_name` | `TEXT` | Имя |
| `last_name` | `TEXT` | Фамилия |
| `language_code` | `TEXT` | Код языка |
| `created_at` | `TIMESTAMP` | Дата создания (авто) |
| `updated_at` | `TIMESTAMP` | Дата обновления (авто) |

**Индексы:** `idx_users_telegram_id`

#### Таблица `chats`

| Колонка | Тип | Описание |
|---------|-----|----------|
| `id` | `BIGSERIAL` | Primary key (автоинкремент) |
| `telegram_id` | `BIGINT` | Уникальный ID чата Telegram |
| `title` | `TEXT` | Название чата |
| `chat_type` | `TEXT` | Тип чата (group/supergroup/channel), default: 'group' |
| `created_at` | `TIMESTAMP` | Дата создания (авто) |
| `updated_at` | `TIMESTAMP` | Дата обновления (авто) |

**Индексы:** `idx_chats_telegram_id`

#### Таблица `moderation_logs`

| Колонка | Тип | Описание |
|---------|-----|----------|
| `id` | `BIGSERIAL` | Primary key (автоинкремент) |
| `user_id` | `BIGINT` | ID пользователя, к которому применено действие |
| `admin_id` | `BIGINT` | ID администратора, выполнившего действие |
| `chat_id` | `BIGINT` | ID чата |
| `action_type` | `TEXT` | Тип: 'mute', 'ban', 'kick', 'unmute' (CHECK constraint) |
| `reason` | `TEXT` | Причина действия |
| `expires_at` | `TIMESTAMP` | Дата окончания мьюта (nullable) |
| `created_at` | `TIMESTAMP` | Дата создания (авто) |

**Индексы:** `idx_moderation_logs_user_id`, `idx_moderation_logs_chat_id`, `idx_moderation_logs_action_type`

### Миграции

| # | Файлы | Описание |
|---|-------|----------|
| 0001 | `0001_create_users.up.sql` / `.down.sql` | Создание таблицы `users` |
| 0002 | `0002_create_chats.up.sql` / `.down.sql` | Создание таблицы `chats` |
| 0003 | `0003_create_moderation_logs.up.sql` / `.down.sql` | Создание таблицы `moderation_logs` |

### Repository слой

**`UserRepository`**:
- `UpsertUser(ctx, user)` — вставка/обновление пользователя
- `FindByUsername(ctx, username)` — поиск telegram_id по username

**`ChatRepository`**:
- `UpsertChat(ctx, chat)` — вставка/обновление чата

**`ModerationLogRepository`** (внутри moderation модуля):
- `LogAction(ctx, log)` — асинхронная запись лога модерации

---

## Модульная система

### Интерфейс модуля

```go
type Module interface {
    Name() string                    // Имя модуля
    Version() string                 // Версия модуля
    Register(bot *botapi.TeleBot, router *Router)  // Регистрация команд
    OnStart(ctx context.Context) error  // Запуск модуля
    OnStop(ctx context.Context) error   // Остановка модуля
}
```

### Registry

```go
type Registry struct {
    modules []Module
}

func (r *Registry) Register(m Module)     // Регистрация модуля
func (r *Registry) RegisterAll(...Module) // Пакетная регистрация
func (r *Registry) StartAll(ctx) error    // Запуск всех модулей (порядок регистрации)
func (r *Registry) StopAll(ctx) error     // Остановка всех модулей (обратный порядок)
func (r *Registry) Modules() []Module     // Получить список модулей
```

### Принцип регистрации

1. `App` создаёт `Registry`
2. Каждый модуль регистрируется через `registry.Register(module)`
3. `registry.StartAll()` вызывает `OnStart()` для каждого модуля
4. При graceful shutdown `registry.StopAll()` вызывает `OnStop()` в **обратном порядке**

---

## Активные модули

### Модуль приветствий (`greeting`)

**Файлы:**
- `internal/modules/greeting/module.go` — реализация `Module`
- `internal/modules/greeting/scheduler.go` — cron-планировщик
- `internal/modules/greeting/repository.go` — отправка сообщения
- `internal/modules/greeting/messages.go` — шаблоны сообщений

**Функциональность:**
- Ежедневная отправка приветственного сообщения в `target_chat_id`
- Cron-расписание: `0 10 * * *` (10:00 AM по таймзоне конфигурации)
- Использует `robfig/cron/v3` для планирования
- Таймзона берётся из `scheduler.timezone` (default: `Europe/Minsk`)

**Жизненный цикл:**
- `OnStart()` → создаёт cron-планировщик, запускает goroutine с таской
- `OnStop()` → останавливает cron-планировщик

---

### Модуль модерации (`moderation`)

**Файлы:**
- `internal/modules/moderation/module.go` — реализация `Module`
- `internal/modules/moderation/commands.go` — парсеры аргументов команд
- `internal/modules/moderation/service.go` — сервис модерации
- `internal/modules/moderation/executor.go` — исполнители Telegram API
- `internal/modules/moderation/repository.go` — репозиторий логов
- `internal/modules/moderation/handlers.go` — хендлеры команд

**Функциональность:**
- `/mute` — замутить пользователя на заданное время
- `/ban` — забанить пользователя
- `/kick` — кикнуть пользователя
- `/unmute` — размутить пользователя

**Компоненты:**

| Компонент | Назначение |
|-----------|-----------|
| `AdminChecker` | Проверка прав администратора через `GetChatMember` |
| `Service` | Сервис-оркестратор, логирование действий |
| `Executor` | Исполнение Telegram API (RestrictChatMember, BanChatMember, Kick) |
| `Repository` | Запись логов в `moderation_logs` |
| `commands.go` | Парсинг аргументов: target, duration, reason |
| `timeutil.Duration` | Парсер длительности: `10s`, `5m`, `2h`, `1d` |

**Формат команд:**

```
/mute @username или ответ на сообщение [длительность] [причина]
/ban @username или ответ на сообщение [причина]
/kick @username или ответ на сообщение
/unmute @username или ответ на сообщение
```

**Аргументы:**
- **Target**: @username ИЛИ ответ на сообщение (reply)
- **Duration** (для mute): `10s`, `5m`, `2h`, `1d` (s=секунды, m=минуты, h=часы, d=дни)
- **Reason**: текстовая причина (опционально)

**Flow модерации:**

```
Команда → Проверка админа → Парсинг аргументов
    → Executor (Telegram API) → Service.logAction (БД) → Reply пользователю
```

---

## Команды бота

| Команда | Описание | Требует прав админа |
|---------|----------|---------------------|
| `/start` | Приветственное сообщение | Нет |
| `/mute [target] [duration] [reason]` | Замутить пользователя | Да |
| `/ban [target] [reason]` | Забанить пользователя | Да |
| `/kick [target]` | Кикнуть пользователя | Да |
| `/unmute [target]` | Размутить пользователя | Да |
| *(unknown)* | Фолбэк для неизвестных команд | Нет |

### Telegram API методы

| Метод | Использование |
|-------|-------------|
| `GetUpdatesChan` | Polling обновлений |
| `Send` | Отправка сообщений |
| `GetChatMember` | Проверка прав администратора |
| `RestrictChatMember` | Мьют/размут |
| `BanChatMember` | Бан |
| `Kick` (Ban + Unban) | Кик |

---

## Жизненный цикл приложения

### Startup

```
main.go
  ├── LoadConfig()          // Viper → YAML + env → validate
  ├── NewLogger()           // Zap production mode
  ├── NewPostgresPool()     // pgxpool из DatabaseConfig
  ├── NewApp()              // Создание App struct
  ├── a.Start(ctx)          // Запуск приложения
  │   ├── bot.NewWithConfig()   // Создание Telegram bot
  │   ├── router.Register()     // Регистрация /start + unknown handler
  │   ├── registry.Register()   // Регистрация greeting + moderation
  │   ├── registry.StartAll()   // OnStart для каждого модуля
  │   ├── a.receiveUpdates()    // Запуск polling goroutine
  │   └── upsertUser/upsertChat // На каждое входящее сообщение
  ├── signal.Notify()       // Подписка на SIGINT/SIGTERM
  └── <-sigChan             // Блокировка до сигнала
```

### Shutdown

```
SIGINT/SIGTERM received
  ├── cancel context        // Отмена всех goroutine
  └── registry.StopAll()    // OnStop в обратном порядке
      ├── moderation.OnStop()
      └── greeting.OnStop()
```

### Обработка сообщения

```
Telegram Update
  ├── router.Route()        // Поиск команды в map[string]Handler
  ├── handler(ctx, update)  // Вызов хендлера
  ├── upsertUser()          // Upsert пользователя в БД
  ├── upsertChat()          // Upsert чата в БД
  └── reply message         // Ответ пользователю
```

---

## Развёртывание

### Локальный запуск

```bash
# 1. Скопировать конфигурацию
cp configs/config.example.yaml configs/config.yaml
# Отредактировать configs/config.yaml

# 2. Запустить PostgreSQL
docker-compose up postgres -d

# 3. Применить миграции
make migrate-up
# или: bash scripts/migrate_up.sh

# 4. Запустить бота
make run
# или: go run ./cmd/bot/main.go
```

### Docker Compose

```yaml
services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: khmurchik_bot
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  bot:
    build: .
    env_file: .env
    depends_on:
      - postgres
    environment:
      - BOT_DATABASE_HOST=postgres
    volumes:
      - ./configs:/app/configs:ro
```

### Multi-stage Docker build

**Stage 1 — builder:**
- Go 1.25-alpine
- `go mod download`, `go mod verify`
- `go build -o /app/bin/bot ./cmd/bot`

**Stage 2 — runtime:**
- Alpine Linux
- Копирование бинарника из builder
- `CMD ["/app/bin/bot"]`

---

## Разработка

### Makefile команды

| Команда | Описание |
|---------|----------|
| `make run` | Запуск бота (`go run ./cmd/bot/main.go`) |
| `make build` | Сборка бинарника (`go build -o bin/bot ./cmd/bot`) |
| `make lint` | Запуск golangci-lint |
| `make test` | Запуск тестов (`go test -v ./...`) |
| `make migrate-up` | Применение миграций |
| `make migrate-down` | Откат последней миграции |
| `make docker-build` | Сборка Docker image |
| `make docker-run` | Запуск через docker-compose |
| `make docker-logs` | Просмотр логов docker-compose |

### Линтинг

Конфигурация: `.golangci.yml`

Включённые линтеры:
- `errcheck` — unchecked errors
- `gosimple` — simplification suggestions
- `govet` — analysis of Go code
- `ineffassign` — detect ineffectual assignments
- `staticcheck` — advanced static analysis
- `unused` — detect unused code

Таймаут: 5 минут

### Утилиты времени

**`timeutil.Duration`** — кастомный тип длительности:
- Поддерживаемые суффиксы: `s` (секунды), `m` (минуты), `h` (часы), `d` (дни)
- Используется парсером аргументов модерации

**`timeutil.LoadTimezone`** — загрузка таймзоны:
- Default: `Europe/Minsk`
- Берётся из `scheduler.timezone`

### DI Container (запланирован)

`internal/di/container.go` содержит структуру для Dependency Injection:

```go
type Container struct {
    Config config.Config
    DB     *pgxpool.Pool
    Logger *zap.Logger
    Bot    *botapi.TeleBot
}
```

На данный момент не используется — `app.go` создаёт зависимости вручную. Подготовлен для будущего использования.

### Планируемые модули

Согласно README, в планах:
- **antispam** — антиспам
- **welcome** — приветствие новых участников
- **stats** — статистика чата
- **polls** — опросы
- **reminders** — напоминания

---

## Примечания по безопасности

- `.env` и `configs/config.yaml` в `.gitignore` — секреты не коммитятся
- `.env.example` — шаблон без реальных значений
- Пароли PostgreSQL в `docker-compose.yml` захардкожены для локальной разработки
- Нет интеграции с внешними менеджерами секретов (Vault, AWS Secrets Manager)

---

## Примечания

- Тесты в проекте отсутствуют (0 тестовых файлов)
- Go-модуль: `github.com/evart2006/khmurchik-community-bot`
- Лицензия: GPLv3
- Минимальная версия Go: 1.25.0
