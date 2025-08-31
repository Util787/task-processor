package config

import (
	"os"
	"strconv"
	"time"
)

// Env variants
const (
	EnvLocal   = "local"
	EnvDev     = "dev"
	EnvProd    = "prod"
	defaultEnv = EnvLocal
)

// Http server default values
const (
	defaultHttpServerHost              = "localhost"
	defaultHttpServerPort              = 8080
	defaultHttpServerReadHeaderTimeout = 5 * time.Second
	defaultHttpServerWriteTimeout      = 5 * time.Second
	defaultHttpServerReadTimeout       = 5 * time.Second
)

// Process queue default values
const (
	defaultProcessQueueWorkers = 4
	defaultProcessQueueSize    = 64
)

type Config struct {
	Env string `env:"ENV"`
	HTTPServerConfig
	TaskProcessQueueConfig
}

type HTTPServerConfig struct {
	Host              string        `env:"HTTP_SERVER_HOST"`
	Port              int           `env:"HTTP_SERVER_PORT"`
	ReadHeaderTimeout time.Duration `env:"HTTP_SERVER_READ_HEADER_TIMEOUT"`
	WriteTimeout      time.Duration `env:"HTTP_SERVER_WRITE_TIMEOUT"`
	ReadTimeout       time.Duration `env:"HTTP_SERVER_READ_TIMEOUT"`
}

type TaskProcessQueueConfig struct {
	Workers   int `env:"WORKERS"`
	QueueSize int `env:"QUEUE_SIZE"`
}

func MustLoadConfig() *Config {
	cfg := &Config{}

	cfg.Env = os.Getenv("ENV")
	if cfg.Env == "" {
		cfg.Env = defaultEnv
	}
	if cfg.Env != EnvDev && cfg.Env != EnvLocal && cfg.Env != EnvProd {
		panic("invalid ENV variable, must be 'dev', 'prod' or 'local'")
	}

	cfg.HTTPServerConfig.Host = os.Getenv("HTTP_SERVER_HOST")
	if cfg.HTTPServerConfig.Host == "" {
		cfg.HTTPServerConfig.Host = defaultHttpServerHost
	}

	var err error
	cfg.HTTPServerConfig.Port, err = strconv.Atoi(os.Getenv("HTTP_SERVER_PORT"))
	if err != nil {
		cfg.HTTPServerConfig.Port = defaultHttpServerPort
	}

	cfg.HTTPServerConfig.ReadHeaderTimeout, err = time.ParseDuration(os.Getenv("HTTP_SERVER_READ_HEADER_TIMEOUT"))
	if err != nil {
		cfg.HTTPServerConfig.ReadHeaderTimeout = defaultHttpServerReadHeaderTimeout
	}

	cfg.HTTPServerConfig.WriteTimeout, err = time.ParseDuration(os.Getenv("HTTP_SERVER_WRITE_TIMEOUT"))
	if err != nil {
		cfg.HTTPServerConfig.WriteTimeout = defaultHttpServerWriteTimeout
	}

	cfg.HTTPServerConfig.ReadTimeout, err = time.ParseDuration(os.Getenv("HTTP_SERVER_READ_TIMEOUT"))
	if err != nil {
		cfg.HTTPServerConfig.ReadTimeout = defaultHttpServerReadTimeout
	}

	cfg.TaskProcessQueueConfig.Workers, err = strconv.Atoi(os.Getenv("WORKERS"))
	if err != nil {
		cfg.TaskProcessQueueConfig.Workers = defaultProcessQueueWorkers
	}

	cfg.TaskProcessQueueConfig.QueueSize, err = strconv.Atoi(os.Getenv("QUEUE_SIZE"))
	if err != nil {
		cfg.TaskProcessQueueConfig.QueueSize = defaultProcessQueueSize
	}

	return cfg
}
