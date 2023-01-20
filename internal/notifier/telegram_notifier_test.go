package notifier

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTelegramNotifier_Notify_Success(t *testing.T) {
	// Arrange
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	n := NewTelegramNotifier("token", 0, map[string]interface{}{
		"parse_mode": "HTML",
	})
	n.apiUrl = testServer.URL

	// Act
	err := n.Notify("test")

	// Assert
	if err != nil {
		t.Error(err)
	}
}

func TestTelegramNotifier_Notify_JsonError(t *testing.T) {
	// Arrange
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	n := NewTelegramNotifier("token", 0, map[string]interface{}{
		"key": math.Inf(1),
	})
	n.apiUrl = testServer.URL

	// Act
	err := n.Notify("test")

	// Assert
	if err == nil {
		t.Errorf("expected error, got none")
	}
	_, ok := err.(*json.UnsupportedValueError)
	if !ok {
		t.Errorf("expected UnsupportedValueError, got %v", err)
	}
}

func TestTelegramNotifier_Notify_HttpError(t *testing.T) {
	// Arrange
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Bad Request"))
	}))
	n := NewTelegramNotifier("token", 0, map[string]interface{}{
		"key": math.Inf(1),
	})
	n.apiUrl = testServer.URL

	// Act
	err := n.Notify("test")

	// Assert
	if err == nil {
		t.Errorf("expected error, got none")
	}
}
