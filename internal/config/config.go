package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Mode   Mode           `toml:"mode"`
	Host   string         `toml:"host"`
	Port   string         `toml:"port"`
	Data   string         `toml:"data"`
	Backup string         `toml:"backup"`
	FS     FSClientConfig `toml:"fs"`
	S3     S3ClientConfig `toml:"s3"`
}

type FSClientConfig struct {
	Root string `toml:"root"`
}

type S3ClientConfig struct {
	Bucket     string `toml:"bucket"`
	AccessKey  string `toml:"access_key"`
	SecretKey  string `toml:"secret_key"`
	S3Endpoint string `toml:"s3_endpoint"`
	Region     string `toml:"region"`
}

func findConfigPath() (string, error) {
	if p := os.Getenv("GOCHIVE_CONFIG"); p != "" {
		return p, nil
	}

	candidates := []string{
		"/opt/gochive/config.toml",
		"./config.toml",
		filepath.Join(os.Getenv("HOME"), ".config/gochive/config.toml"),
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c, nil
		}
	}

	return "", fmt.Errorf("no config file found in known locations")
}

func LoadConfig() (*Config, error) {
	path, err := findConfigPath()
	if err != nil {
		return nil, err
	}

	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config %s: %w", path, err)
	}

	if k := os.Getenv("S3_ACCESS_KEY"); k != "" {
		cfg.S3.AccessKey = k
	}
	if k := os.Getenv("S3_SECRET_KEY"); k != "" {
		cfg.S3.SecretKey = k
	}

	return &cfg, nil
}
