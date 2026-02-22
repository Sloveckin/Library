package config

import (
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string     `yaml:"env"`
	StorageUrl string     `yaml:"storage_url" env-required:"true"`
	HttpServer HttpServer `yaml:"http_server"`
}

type HttpServer struct {
	Address     string        `yaml:"address"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

func MustLoad() *Config {
	config := &Config{}

	config.Env = os.Getenv("ENV")
	config.StorageUrl = os.Getenv("STORAGE_URL")
	config.HttpServer.Address = os.Getenv("HTTP_SERVER_ADDRESS")

	var err error
	config.HttpServer.Timeout, err = time.ParseDuration(os.Getenv("HTTP_SERVER_TIMEOUT"))
	if err != nil {
		log.Fatalf("Cannot parse HTTP_SERVER_TIMEOUT: %v", err)
	}
	config.HttpServer.IdleTimeout, err = time.ParseDuration(os.Getenv("HTTP_SERVER_IDLE_TIMEOUT"))
	if err != nil {
		log.Fatalf("Cannot parse HTTP_SERVER_IDLE_TIMEOUT: %v", err)
	}

	// Логирование
	log.Printf("Config loaded:")
	log.Printf("  Env: %s", config.Env)
	log.Printf("  StorageUrl: %s", config.StorageUrl)
	log.Printf("  HttpServer.Address: %s", config.HttpServer.Address)
	log.Printf("  HttpServer.Timeout: %v", config.HttpServer.Timeout)
	log.Printf("  HttpServer.IdleTimeout: %v", config.HttpServer.IdleTimeout)
	return config
}
