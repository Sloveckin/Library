package config

import (
	"os"
	"testing"
	"time"
)

func TestMustLoad(t *testing.T) {
	if err := os.Setenv("ENV", "test"); err != nil {
    	t.Fatalf("setenv ENV: %v", err)
	}
	if err := os.Setenv("STORAGE_URL", "postgres://localhost"); err != nil {
		t.Fatalf("setenv STORAGE_URL: %v", err)
	}
	if err := os.Setenv("HTTP_SERVER_ADDRESS", "localhost:8080"); err != nil {
		t.Fatalf("setenv HTTP_SERVER_ADDRESS: %v", err)
	}
	if err := os.Setenv("HTTP_SERVER_TIMEOUT", "10s"); err != nil {
		t.Fatalf("setenv HTTP_SERVER_TIMEOUT: %v", err)
	}
	if err := os.Setenv("HTTP_SERVER_IDLE_TIMEOUT", "5s"); err != nil {
		t.Fatalf("setenv HTTP_SERVER_IDLE_TIMEOUT: %v", err)
	}

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

	if err := os.Unsetenv("ENV"); err != nil {
    	t.Fatalf("unset ENV: %v", err)
	}
	if err := os.Unsetenv("STORAGE_URL"); err != nil {
		t.Fatalf("unset STORAGE_URL: %v", err)
	}
	if err := os.Unsetenv("HTTP_SERVER_ADDRESS"); err != nil {
		t.Fatalf("unset HTTP_SERVER_ADDRESS: %v", err)
	}
	if err := os.Unsetenv("HTTP_SERVER_TIMEOUT"); err != nil {
		t.Fatalf("unset HTTP_SERVER_TIMEOUT: %v", err)
	}
	if err := os.Unsetenv("HTTP_SERVER_IDLE_TIMEOUT"); err != nil {
		t.Fatalf("unset HTTP_SERVER_IDLE_TIMEOUT: %v", err)
	}
}
