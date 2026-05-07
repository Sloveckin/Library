package config

import (
	"os"
	"testing"
	"time"
)

func TestMustLoad(t *testing.T) {
	os.Setenv("ENV", "test")
	os.Setenv("STORAGE_URL", "postgres://localhost")
	os.Setenv("HTTP_SERVER_ADDRESS", "localhost:8080")
	os.Setenv("HTTP_SERVER_TIMEOUT", "10s")
	os.Setenv("HTTP_SERVER_IDLE_TIMEOUT", "5s")

	config := MustLoad()

	if config.Env != "test" {
		t.Errorf("Expected env 'test', got %s", config.Env)
	}
	if config.StorageUrl != "postgres://localhost" {
		t.Errorf("Expected storage_url 'postgres://localhost', got %s", config.StorageUrl)
	}
	if config.HttpServer.Address != "localhost:8080" {
		t.Errorf("Expected address 'localhost:8080', got %s", config.HttpServer.Address)
	}
	if config.HttpServer.Timeout != 10*time.Second {
		t.Errorf("Expected timeout 10s, got %v", config.HttpServer.Timeout)
	}
	if config.HttpServer.IdleTimeout != 5*time.Second {
		t.Errorf("Expected idle timeout 5s, got %v", config.HttpServer.IdleTimeout)
	}

	os.Unsetenv("ENV")
	os.Unsetenv("STORAGE_URL")
	os.Unsetenv("HTTP_SERVER_ADDRESS")
	os.Unsetenv("HTTP_SERVER_TIMEOUT")
	os.Unsetenv("HTTP_SERVER_IDLE_TIMEOUT")
}
