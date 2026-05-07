package handler

import (
	"testing"
)

func TestOk(t *testing.T) {
	resp := Ok()
	if resp.Status != StatusOK {
		t.Errorf("Expected status %s, got %s", StatusOK, resp.Status)
	}
	if resp.Error != "" {
		t.Errorf("Expected empty error, got %s", resp.Error)
	}
}

func TestError(t *testing.T) {
	errMsg := "test error"
	resp := Error(errMsg)
	if resp.Status != StatusError {
		t.Errorf("Expected status %s, got %s", StatusError, resp.Status)
	}
	if resp.Error != errMsg {
		t.Errorf("Expected error %s, got %s", errMsg, resp.Error)
	}
}
