package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

const (
	EnvLocal = "local"
	EnvProd  = "prod"
)

type Config struct {
	Env        string
	HTTPServer *HTTPServer `mapstructure:"http_server"`
	DB         *DB         `mapstructure:"db"`
}

type HTTPServer struct {
	Address     string        `mapstructure:"address"`
	Timeout     time.Duration `mapstructure:"timeout"`
	IdleTimeout time.Duration `mapstructure:"idle_timeout"`
}
type DB struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	Username string
	Password string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	env, err := getEnv("APP_ENV")
	if err != nil {
		return nil, fmt.Errorf("failed to get env variable: %w", err)
	}

	if env != EnvLocal && env != EnvProd {
		return nil, fmt.Errorf("%q is not a valid env", env)
	}

	var cfg Config

	//	Set app environment
	cfg.Env = env

	//	Viper setup
	v := viper.New()
	v.SetConfigFile(filepath.Join("config", env+".yml"))

	if err = v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed reading config: %w", err)
	}

	//	Unmarshal config file's data into Config struct
	if err = v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode into struct: %w", err)
	}

	//	Get and set DB_USERNAME & DB_PASSWORD
	dbUser, err := getEnv("DB_USERNAME")
	if err != nil {
		return nil, fmt.Errorf("failed to get env variable: %w", err)
	}

	dbPassword, err := getEnv("DB_PASSWORD")
	if err != nil {
		return nil, fmt.Errorf("failed to get env variable: %w", err)
	}

	cfg.DB.Username = dbUser
	cfg.DB.Password = dbPassword

	return &cfg, nil
}

// getEnv is a helper to easily get values
// from .env file and return error on failure
func getEnv(key string) (string, error) {
	v := os.Getenv(key)

	if strings.TrimSpace(v) == "" {
		return "", fmt.Errorf("%q is not set in env file", key)
	}

	return v, nil
}
