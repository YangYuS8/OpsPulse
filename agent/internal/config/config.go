package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	NodeID           string        `yaml:"node_id"`
	ServerURL        string        `yaml:"server_url"`
	Token            string        `yaml:"token"`
	Interval         time.Duration `yaml:"interval"`
	DockerEnabled    bool          `yaml:"docker_enabled"`
	ServiceWhitelist []string      `yaml:"service_whitelist"`
}

func Load(path string) (Config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	if err := yaml.Unmarshal(content, &cfg); err != nil {
		return Config{}, err
	}
	if cfg.NodeID == "" || cfg.ServerURL == "" || cfg.Token == "" {
		return Config{}, fmt.Errorf("node_id, server_url and token are required")
	}
	if cfg.Interval == 0 {
		cfg.Interval = 15 * time.Second
	}
	return cfg, nil
}
