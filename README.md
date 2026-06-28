# simple-moderation-community-bot

Telegram-бот для модерации community-чатов на Go.

Бот работает в любых `group`/`supergroup`, куда его добавили. Модерация всегда применяется к текущему чату команды, а список админов и права проверяются через Telegram Bot API.

## Возможности

- `/ban`, `/mute`, `/kick`, `/unmute` для админов
- `/balabol`, `/unbalabol` — внутренняя метка, запрещающая пользователю отправлять репорты
- `/report` — репорт на сообщение с пингом админов в этом же чате
- Настраиваемое утреннее приветствие для каждого чата
- Хранение чатов, пользователей, репортов, меток и moderation logs в PostgreSQL
- Конфигурация через env, YAML опционален

## Запуск

```bash
cp .env.example .env

docker-compose up -d postgres
make migrate-up
make run
```

Для Docker:

```bash
docker-compose up -d
```

## Env

| Переменная | Описание |
|---|---|
| `BOT_TOKEN` | Токен бота от `@BotFather` |
| `BOT_TARGET_CHAT_ID` | Optional legacy fallback; модерация работает в текущем чате |
| `BOT_POLL_TIMEOUT` | Long polling timeout, default `30` |
| `BOT_ALLOWED_UPDATES` | Например `message,my_chat_member` |
| `BOT_LOG_LEVEL` | `debug`, `info`, `warn`, `error` |
| `BOT_DATABASE_HOST` | Адрес PostgreSQL |
| `BOT_DATABASE_PORT` | Порт PostgreSQL |
| `BOT_DATABASE_USER` | Пользователь БД |
| `BOT_DATABASE_PASSWORD` | Пароль БД |
| `BOT_DATABASE_NAME` | Имя БД |
| `BOT_DATABASE_SSL_MODE` | SSL mode |
| `BOT_DEFAULT_TIMEZONE` | Default timezone для новых чатов, default `Europe/Minsk` |
| `BOT_DEFAULT_GREETING_TIME` | Default greeting time, default `10:00` |
| `BOT_DEFAULT_GREETING_ENABLED` | Включать приветствие для новых чатов |
| `BOT_DEFAULT_GREETING_MESSAGE` | Default greeting text |

Старые `BOT_BOT_TOKEN`, `BOT_BOT_TARGET_CHAT_ID`, `BOT_SERVER_POLL_TIMEOUT`, `BOT_SCHEDULER_TIMEZONE` остаются поддержанными для совместимости.

## Команды

Пользовательские:

| Команда | Описание |
|---|---|
| `/report причина` | Reply на сообщение → отправить репорт админам |

Админские:

| Команда | Описание |
|---|---|
| `/mute 1h причина` | Reply на сообщение → мут |
| `/ban причина` | Reply на сообщение → бан |
| `/kick` | Reply на сообщение → кик |
| `/unmute` | Reply на сообщение → снять мут |
| `/balabol причина` | Reply на сообщение → запретить пользователю отправлять репорты |
| `/unbalabol` | Reply на сообщение → снять метку |

Fallback `@username` поддержан, только если бот уже видел пользователя и сохранил его в БД.

Настройки приветствия:

| Команда | Описание |
|---|---|
| `/greeting_status` | Показать текущие настройки |
| `/greeting_on` | Включить приветствие |
| `/greeting_off` | Выключить приветствие |
| `/greeting_time 10:00 Europe/Minsk` | Настроить локальное время |
| `/greeting_text текст` | Настроить текст |

## Telegram-права

Для модерации бот должен быть администратором в чате и иметь право ограничивать пользователей (`can_restrict_members`). Команды также проверяют права администратора у вызывающего пользователя через `getChatMember`.

## Проверки

```bash
go test ./...
go vet ./...
go build ./cmd/bot
```
