# khmurchik-community-bot

Многомодульный Telegram-бот для чатика Павла Хмурчика по технике.

## MVP-функции

- Ежедневное приветствие в 10:00 (Europe/Minsk)
- Модерация чата: /mute, /ban, /kick, /unmute

## Запуск

### Локально

```bash
# 1. Скопируйте конфиги
cp .env.example .env
cp configs/config.example.yaml configs/config.yaml

# 2. Заполните .env (BOT_BOT_TOKEN, BOT_BOT_TARGET_CHAT_ID и т.д.)

# 3. Поднимите PostgreSQL
docker-compose up -d postgres

# 4. Примените миграции
make migrate-up

# 5. Запустите бота
make run
```

### Docker

```bash
# Заполните .env, затем:
docker-compose up -d
```

## Env-переменные

| Переменная | Описание |
|---|---|
| BOT_BOT_TOKEN | Токен бота от @BotFather |
| BOT_BOT_TARGET_CHAT_ID | Chat ID для отправки приветствий и модерации |
| BOT_DATABASE_HOST | Адрес PostgreSQL |
| BOT_DATABASE_PORT | Порт PostgreSQL |
| BOT_DATABASE_USER | Пользователь БД |
| BOT_DATABASE_PASSWORD | Пароль БД |
| BOT_DATABASE_NAME | Имя БД |
| BOT_DATABASE_SSL_MODE | SSL режим |

## Команды модерации

Все команды — **reply-первый сценарий**. Укажите @username только если бот уже видел пользователя (сохранён в users).

| Команда | Описание |
|---|---|
| `/mute 1d "причина"` | Reply на сообщение → мут reply sender на 1 день |
| `/mute @user 1d "причина"` | Мут пользователя по username на 1 день |
| `/ban "причина"` | Reply на сообщение → бан reply sender |
| `/ban @user "причина"` | Бан пользователя по username |
| `/kick` | Reply на сообщение → кик reply sender |
| `/kick @user` | Кик пользователя по username |
| `/unmute` | Reply на сообщение → снять мут с reply sender |
| `/unmute @user` | Снять мут с пользователя по username |

Формат времени: `30m`, `3h`, `1d`, `7d` (s = секунды, m = минуты, h = часы, d = дни).

Проверка прав администратора — через Telegram API (GetChatMember).

## Архитектура

Бот модульный. Каждый модуль — в `internal/modules/`.

Будущие модули (planned): antispam, welcome, stats, polls, reminders.

## Makefile

| Команда | Описание |
|---|---|
| `make run` | Запустить бота локально |
| `make build` | Собрать бинарник (bin/bot) |
| `make test` | Запустить тесты |
| `make lint` | Запустить golangci-lint |
| `make migrate-up` | Применить все миграции |
| `make migrate-down` | Откатить одну миграцию |
| `make docker-up` | Поднять Docker-стек |
| `make docker-down` | Остановить Docker-стек |
