package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	AppEnv           string `mapstructure:"APP_ENV" required:"true"`
	ServerAddress    string `mapstructure:"SERVER_ADDRESS" required:"true"`
	ContextTimeout   int    `mapstructure:"CONTEXT_TIMEOUT" required:"true"`
	RedisHost        string `mapstructure:"REDIS_HOST" required:"true"`
	RedisPort        int    `mapstructure:"REDIS_PORT" required:"true"`
	RedisPassword    string `mapstructure:"REDIS_PASSWORD"`
	EmailFrom        string `mapstructure:"EMAIL_FROM" required:"true"`
	EmailTo          string `mapstructure:"EMAIL_TO" required:"true"`
	EmailPassword    string `mapstructure:"EMAIL_PASSWORD" required:"true"`
	SMTPHost         string `mapstructure:"SMTP_HOST" required:"true"`
	SMTPPort         int    `mapstructure:"SMTP_PORT" required:"true"`
	CSVFilePath      string `mapstructure:"CSV_FILE_PATH" required:"true"`
	FakeEmail        bool   `mapstructure:"FAKE_EMAIL" required:"true"`
	RateLimit        int    `mapstructure:"RATE_LIMIT" required:"true"`
	RedisTimeoutSec  int    `mapstructure:"REDIS_TIMEOUT_SEC" required:"true"`
	CacheDurationSec int    `mapstructure:"CACHE_DURATION_SEC" required:"true"`
	DBHost           string `mapstructure:"DB_HOST" required:"true"`
	DBUser           string `mapstructure:"DB_USER" required:"true"`
	DBPassword       string `mapstructure:"DB_PASSWORD" required:"true"`
	DBName           string `mapstructure:"DB_NAME" required:"true"`
	DBPort           int    `mapstructure:"DB_PORT" required:"true"`
}

func NewEnv(envFile string) *Env {
	viper.SetConfigType("env")
	viper.SetConfigFile(envFile)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Can't find the file %s: %s", envFile, err)
	}

	viper.AutomaticEnv()
	viper.SetDefault("CONTEXT_TIMEOUT", 5)
	viper.SetDefault("REDIS_PORT", 6379)
	viper.SetDefault("SMTP_PORT", 587)
	viper.SetDefault("FAKE_EMAIL", false)
	viper.SetDefault("RATE_LIMIT", 1000)
	viper.SetDefault("REDIS_TIMEOUT_SEC", 5)
	viper.SetDefault("CACHE_DURATION_SEC", 600)
	viper.SetDefault("DB_PORT", 5432)

	var env Env
	if err := viper.Unmarshal(&env); err != nil {
		log.Fatalf("Environment can't be loaded: %s", err)
	}

	if env.AppEnv == "development" {
		log.Println("The App is running in development environment")
	}

	return &env
}

func (e *Env) Validate() error {
	requiredFields := []string{
		"APP_ENV",
		"SERVER_ADDRESS",
		"REDIS_HOST",
		"EMAIL_FROM",
		"EMAIL_TO",
		"EMAIL_PASSWORD",
		"SMTP_HOST",
		"CSV_FILE_PATH",
		"DB_HOST",
		"DB_USER",
		"DB_PASSWORD",
		"DB_NAME",
	}
	for _, field := range requiredFields {
		if viper.GetString(field) == "" {
			return fmt.Errorf("required environment variable %s not set", field)
		}
	}
	return nil
}
