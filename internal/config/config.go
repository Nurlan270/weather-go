package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"

	gpv "github.com/go-playground/validator/v10"
)

const (
	EnvLocal = "local"
	EnvProd  = "prod"
)

type Config struct {
	App         App         `mapstructure:"app"         validate:"required"`
	Session     Session     `mapstructure:"session"     validate:"required"`
	Cache       Cache       `mapstructure:"cache"       validate:"required"`
	HTTPServer  HTTPServer  `mapstructure:"http_server" validate:"required"`
	DB          DB          `mapstructure:"db"          validate:"required"`
	OpenWeather OpenWeather `mapstructure:"openweather" validate:"required"`
}

type App struct {
	Env string `mapstructure:"env" validate:"required"`
}

type Session struct {
	Name      string        `mapstructure:"name"       validate:"required"`
	ExpiresIn time.Duration `mapstructure:"expires_in" validate:"required"`
}

type Cache struct {
	ExpiresIn time.Duration `mapstructure:"expires_in" validate:"required"`
}

type HTTPServer struct {
	Address     string        `mapstructure:"address"      validate:"required"`
	Timeout     time.Duration `mapstructure:"timeout"      validate:"required"`
	IdleTimeout time.Duration `mapstructure:"idle_timeout" validate:"required"`
}
type DB struct {
	Host     string `mapstructure:"host"     validate:"required"`
	Port     string `mapstructure:"port"     validate:"required"`
	Name     string `mapstructure:"name"     validate:"required"`
	Username string `mapstructure:"username" validate:"required"`
	Password string `mapstructure:"password" validate:"required"`
}

type OpenWeather struct {
	ApiKey string `mapstructure:"api_key" validate:"required"`
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	env, err := getEnv("APP_ENV")
	if err != nil {
		return nil, err
	}

	if env != EnvLocal && env != EnvProd {
		return nil, fmt.Errorf("%q is not a valid env", env)
	}

	var cfg Config

	//	Set app environment
	cfg.App.Env = env

	//	Viper setup
	v := viper.New()
	v.SetConfigFile(filepath.Join("config", env+".yml"))
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	//	Bind variables from env file.
	//	Dot notation used, because Viper by default uses dot notation.
	//	It will then be treated as snake case as we set it above
	//	(e.g. db.username is treated as DB_USERNAME from env file)
	var keys = []string{
		"db.username", "db.password",
		"openweather.api_key",
	}

	if err = bindEnvs(v, keys...); err != nil {
		return nil, err
	}

	//	Read config file
	if err = v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed reading config: %w", err)
	}

	//	Unmarshal config file's data into Config struct
	if err = v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode into struct: %w", err)
	}

	//	Validate final config
	validate := gpv.New(gpv.WithRequiredStructEnabled())

	if err = validate.Struct(cfg); err != nil {
		return nil, fmt.Errorf("failed to validate config: %w", err)
	}

	return &cfg, nil
}

// getEnv is a helper function which gets values
// from env file. It returns env variable's value
// and error if its value is empty or not set at all.
func getEnv(key string) (string, error) {
	v, ok := os.LookupEnv(key)

	if ok && len(strings.TrimSpace(v)) != 0 {
		return v, nil
	}

	return "", fmt.Errorf("%q is empty or not set in .env file", key)
}

// bindEnvs is a wrapper around Viper's BindEnv.
// It's running BindEnv in a loop to easily bind
// multiple env variables at once.
func bindEnvs(v *viper.Viper, keys ...string) error {
	for _, k := range keys {
		if err := v.BindEnv(k); err != nil {
			return fmt.Errorf("failed to bind: %w", err)
		}
	}

	return nil
}
