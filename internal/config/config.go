package config

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Mode   Mode
	Host   string
	Port   string
	DBRoot string
	FS     FSClientConfig
	S3     S3ClientConfig
}

type S3ClientConfig struct {
	Bucket     string
	AccessKey  string
	SecretKey  string
	S3Endpoint string
	Region     string
}

type FSClientConfig struct {
	Root string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		slog.Warn("no .env file found, relying on real env vars")
	}

	m, err := strconv.Atoi(os.Getenv("MODE"))
	if err != nil {
		return nil, err
	}

	return &Config{
		Mode:   Mode(m),
		Host:   os.Getenv("HOST"),
		Port:   os.Getenv("PORT"),
		DBRoot: os.Getenv("DB_ROOT"),
		FS: FSClientConfig{
			Root: os.Getenv("ROOT"),
		},
		S3: S3ClientConfig{
			Bucket:     os.Getenv("BUCKET"),
			AccessKey:  os.Getenv("ACCESS_KEY"),
			SecretKey:  os.Getenv("SECRET_KEY"),
			S3Endpoint: os.Getenv("S3_ENDPOINT"),
			Region:     os.Getenv("REGION"),
		},
	}, nil
}
