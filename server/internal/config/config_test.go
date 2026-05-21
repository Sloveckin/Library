package config

import (
	"os"
	"testing"
	"time"
)

func TestMustLoad(t *testing.T) {
	env := map[string]string{
		"ENV":                      "test",
		"STORAGE_URL":              "postgres://localhost",
		"HTTP_SERVER_ADDRESS":      "localhost:8080",
		"HTTP_SERVER_TIMEOUT":      "10s",
		"HTTP_SERVER_IDLE_TIMEOUT": "5s",
	}

	for key, value := range env {
		if err := os.Setenv(key, value); err != nil {
			t.Fatalf("setenv %s: %v", key, err)
		}
	}
	t.Cleanup(func() {
		for key := range env {
			if err := os.Unsetenv(key); err != nil {
				t.Fatalf("unset %s: %v", key, err)
			}
		}
	})

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

}
