package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTPServer     `yaml:"http_server"`
	PostgresConfig `yaml:"postgres_config"`
}

type HTTPServer struct {
	Host         string        `yaml:"host"`
	Port         string        `yaml:"port"`
	Timeout      time.Duration `yaml:"timeout"`
	Idle_timeout time.Duration `yaml:"idle_timeout"`
}

type PostgresConfig struct {
	Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port     string `yaml:"port" `
	Database string `yaml:"database"`
	Username string `env:"DB_USERNAME"`
	Password string `env:"DB_PASSWORD"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config %s", err)
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
