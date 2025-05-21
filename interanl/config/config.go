package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env      string `yaml:"env" env-requierd:"true"`
	Storage  string `yaml:"storage_path" env-requierd:"true"`
	GRPC     `yaml:"grpc"`
	TokenTTL time.Duration `yaml:"token_ttl" env-default:"1h"`
}

type GRPC struct {
	Port    int `yaml:"port"`
	Timeout time.Duration
}

func MustLoad() *Config {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		panic("config path is empty")
	}

	return MustLoadByPath(path)
}

func MustLoadByPath(path string) *Config {

	if _, err := os.Stat(path); err != nil {
		panic("config dose not exists by this path: " + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("cannot read config:" + err.Error())
	}

	return &cfg
}
