package config

type Config struct {
	Bot       BotConfig       `mapstructure:"bot"`
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Scheduler SchedulerConfig `mapstructure:"scheduler"`
}

type BotConfig struct {
	Token        string `mapstructure:"token"        validate:"required"`
	TargetChatID int64  `mapstructure:"target_chat_id" validate:"required"`
}

type ServerConfig struct {
	PollTimeout int32 `mapstructure:"poll_timeout" validate:"required"`
}

type DatabaseConfig struct {
	Host            string `mapstructure:"host"            validate:"required"`
	Port            int    `mapstructure:"port"            validate:"required"`
	User            string `mapstructure:"user"            validate:"required"`
	Password        string `mapstructure:"password"        validate:"required"`
	Name            string `mapstructure:"name"            validate:"required"`
	SSLMode         string `mapstructure:"ssl_mode"        validate:"omitempty"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"  validate:"gte=0"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"  validate:"gte=0"`
	ConnMaxLifetime string `mapstructure:"conn_max_lifetime"`
}

type SchedulerConfig struct {
	Timezone string `mapstructure:"timezone" validate:"required"`
}
