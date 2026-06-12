package config

import (
	"errors"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

const DefaultGreetingMessage = "Я робот, и Паша заставил меня каждое утро желать вам охуенного дня!\n\nОбнял, покружил, на место поставил! 😘"

func Load(path string) (*Config, error) {
	v := viper.New()

	setDefaults(v)

	v.SetEnvPrefix("BOT")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	bindEnv(v)

	if path != "" {
		v.SetConfigFile(path)
		if err := v.ReadInConfig(); err != nil {
			var notFound viper.ConfigFileNotFoundError
			if !errors.As(err, &notFound) && !os.IsNotExist(err) {
				return nil, err
			}
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return nil, err
	}

	cfg.Server.AllowedUpdates = normalizeAllowedUpdates(cfg.Server.AllowedUpdates)

	return &cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("server.poll_timeout", 30)
	v.SetDefault("server.allowed_updates", []string{"message", "my_chat_member"})
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "postgres")
	v.SetDefault("database.name", "khmurchik_bot")
	v.SetDefault("database.ssl_mode", "disable")
	v.SetDefault("database.max_open_conns", 25)
	v.SetDefault("database.max_idle_conns", 2)
	v.SetDefault("database.conn_max_lifetime", "5m")
	v.SetDefault("defaults.timezone", "Europe/Minsk")
	v.SetDefault("defaults.greeting_time", "10:00")
	v.SetDefault("defaults.greeting_enabled", true)
	v.SetDefault("defaults.greeting_message", DefaultGreetingMessage)
	v.SetDefault("log.level", "info")
}

func bindEnv(v *viper.Viper) {
	_ = v.BindEnv("bot.token", "BOT_TOKEN", "BOT_BOT_TOKEN")
	_ = v.BindEnv("bot.target_chat_id", "BOT_TARGET_CHAT_ID", "BOT_BOT_TARGET_CHAT_ID")
	_ = v.BindEnv("server.poll_timeout", "BOT_POLL_TIMEOUT", "BOT_SERVER_POLL_TIMEOUT")
	_ = v.BindEnv("server.allowed_updates", "BOT_ALLOWED_UPDATES", "BOT_SERVER_ALLOWED_UPDATES")
	_ = v.BindEnv("database.host", "BOT_DATABASE_HOST")
	_ = v.BindEnv("database.port", "BOT_DATABASE_PORT")
	_ = v.BindEnv("database.user", "BOT_DATABASE_USER")
	_ = v.BindEnv("database.password", "BOT_DATABASE_PASSWORD")
	_ = v.BindEnv("database.name", "BOT_DATABASE_NAME")
	_ = v.BindEnv("database.ssl_mode", "BOT_DATABASE_SSL_MODE")
	_ = v.BindEnv("database.max_open_conns", "BOT_DATABASE_MAX_OPEN_CONNS")
	_ = v.BindEnv("database.max_idle_conns", "BOT_DATABASE_MAX_IDLE_CONNS")
	_ = v.BindEnv("database.conn_max_lifetime", "BOT_DATABASE_CONN_MAX_LIFETIME")
	_ = v.BindEnv("defaults.timezone", "BOT_DEFAULT_TIMEZONE", "BOT_SCHEDULER_TIMEZONE")
	_ = v.BindEnv("defaults.greeting_time", "BOT_DEFAULT_GREETING_TIME")
	_ = v.BindEnv("defaults.greeting_enabled", "BOT_DEFAULT_GREETING_ENABLED")
	_ = v.BindEnv("defaults.greeting_message", "BOT_DEFAULT_GREETING_MESSAGE")
	_ = v.BindEnv("log.level", "BOT_LOG_LEVEL")
}

func normalizeAllowedUpdates(updates []string) []string {
	if len(updates) == 0 {
		return []string{"message", "my_chat_member"}
	}

	var normalized []string
	for _, item := range updates {
		for _, part := range strings.Split(item, ",") {
			part = strings.TrimSpace(part)
			if part != "" {
				normalized = append(normalized, part)
			}
		}
	}
	return normalized
}
