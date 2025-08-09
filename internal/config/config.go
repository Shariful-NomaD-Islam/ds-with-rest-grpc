package config

import (
	"fmt"
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server  ServerConfig   `yaml:"server"`
	Workers []WorkerConfig `yaml:"workers"`
	GRPC    GRPCConfig     `yaml:"grpc"`
	Logging LoggingConfig  `yaml:"logging"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

type WorkerConfig struct {
	URL string `yaml:"url"`
	ID  string `yaml:"id"`
}

type GRPCConfig struct {
	Timeout    string `yaml:"timeout"`
	MaxRetries int    `yaml:"max_retries"`
}

type LoggingConfig struct {
	Level string `yaml:"level"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}

	if len(c.Workers) == 0 {
		return fmt.Errorf("at least one worker must be configured")
	}

	for i, worker := range c.Workers {
		if worker.URL == "" {
			return fmt.Errorf("worker %d: URL is required", i)
		}
		if worker.ID == "" {
			return fmt.Errorf("worker %d: ID is required", i)
		}
	}

	return nil
}

func (c *Config) GetGRPCTimeout() time.Duration {
	if c.GRPC.Timeout == "" {
		return 10 * time.Second
	}

	timeout, err := time.ParseDuration(c.GRPC.Timeout)
	if err != nil {
		return 10 * time.Second
	}

	return timeout
}

func (c *Config) GetServerAddress() string {
	if c.Server.Host == "" {
		return ":" + c.Server.Port
	}
	return c.Server.Host + ":" + c.Server.Port
}

func (c *Config) GetWorkerURLs() []string {
	urls := make([]string, len(c.Workers))
	for i, worker := range c.Workers {
		urls[i] = worker.URL
	}
	return urls
}
