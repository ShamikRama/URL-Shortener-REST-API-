package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml: "env" env-default: "local"`
	Storage    string `yaml: "storage_path" env-required:"true"`
	HTTPServer `yaml: "http_server" `
}

type HTTPServer struct {
	Address      string        `yaml: address env-default:"localhost:8080"`
	Timeout      time.Duration `yaml: timeout env-default:"4s"`
	Idle_timeout time.Duration `yaml: idle_timeout env-default:"60s"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// проверка на существование файла
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("File %s does not exist", configPath)
	}

	var cfg Config

	// читаем конфиг
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("cannot read config : %s", err)
	}

	return &cfg
}
