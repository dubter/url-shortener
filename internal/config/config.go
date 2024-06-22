package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"os"
	"time"
)

type Config struct {
	Env           string `env:"ENV" env-default:"local"`
	Postgres      PostgresConfig
	Redis         RedisConfig
	Http          HTTPConfig
	TemplatesPath string `env:"TEMPLATES_PATH" env-required:"true"`
}

type HTTPConfig struct {
	Port            string        `env:"PORT"  env-default:"8080"`
	ReadTimeout     time.Duration `env:"READ_TIMEOUT" env-default:"10s"`
	WriteTimeout    time.Duration `env:"WRITE_TIMEOUT" env-default:"10s"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" env-default:"10s"`
	Limiter         Limiter
}

type Limiter struct {
	RPS   int           `env:"LIMITER_RPS" env-default:"10"`
	Burst int           `env:"LIMITER_BURST" env-default:"20"`
	TTL   time.Duration `env:"LIMITER_TTL" env-default:"10m"`
}

type PostgresConfig struct {
	PostgresURL string `env:"POSTGRES_URL" env-required:"true"`
}

type RedisConfig struct {
	Hosts    []string `env:"REDIS_HOSTS" yaml:"hosts" env-required:"true"`
	Password string   `env:"REDIS_PASSWORD" env-required:"true"`
}

func LoadConfig() (*Config, error) {
	envPath := fetchPath()

	if envPath == "" {
		return nil, fmt.Errorf("'.env' file path is empty")
	}

	if err := godotenv.Load(envPath); err != nil {
		return nil, fmt.Errorf("no .env file found")
	}

	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("can not read config: %w", err)
	}

	return &cfg, nil
}

func fetchPath() string {
	var envPath string

	flag.StringVar(&envPath, "env", "", "path to '.env' file")
	flag.Parse()

	if envPath == "" {
		envPath = os.Getenv("ENV_PATH")
	}

	return envPath
}
