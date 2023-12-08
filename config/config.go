package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const configPath string = "./config/config.yml"

type Config struct {
	ZapLogger  ZapLogger  `yaml:"zap_logger"`
	Postgres   Postgres   `yaml:"postgres"`
	HTTPServer HTTPServer `json:"http_server"`
}

type ZapLogger struct {
	Test             bool     `yaml:"test"`
	Level            string   `yaml:"level"`
	OutputPaths      []string `yaml:"output_paths"`
	ErrorOutputPaths []string `yaml:"error_output_paths"`
}

type Postgres struct {
	Host     string `env:"POSTGRES_HOST" env-required:"true"`
	Port     int    `env:"POSTGRES_PORT" env-required:"true"`
	User     string `env:"POSTGRES_USER" env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	DBName   string `env:"POSTGRES_DB" env-required:"true"`
	SSLMode  string `yaml:"ssl_mode" env:"POSTGRES_SSLMODE" env-required:"true"`
	MaxConns int32  `yaml:"max_conns"`
	MinConns int32  `yaml:"min_conns"`
}

type HTTPServer struct {
	Host            string        `env:"HTTP_SERVER_HOST" env-required:"true"`
	Port            int           `env:"HTTP_SERVER_PORT" env-required:"true"`
	ReadTimeout     time.Duration `yaml:"read_timeout" env-default:"3s"`
	WriteTimeout    time.Duration `yaml:"write_timeout" env-default:"5s"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env-default:"5s"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return nil, err
	}

	err = cleanenv.UpdateEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("error updating env: %w", err)
	}

	return &cfg, nil
}
