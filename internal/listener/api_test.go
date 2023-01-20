package listener

import (
	"bc-deployer/internal/config"
	"bc-deployer/internal/runner/runner_repo"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestApi() *Api {
	return NewApi(0, runner_repo.NewFakeRepo(true))
}

func newTestWebhook() (string, config.Webhook) {
	whName := "test"
	wh := config.Webhook{
		Path:       "/test",
		Validators: config.Validators{},
		Mapenv:     nil,
		Tasks:      nil,
	}
	return whName, wh
}

func TestApi_WebhookHandler_Success(t *testing.T) {
	// Arrange
	expected := "Webhook test executed"
	whName, wh := newTestWebhook()
	wh.Tasks = []string{"fake"}
	handler := newTestApi().createHandler(whName, wh)
	req := httptest.NewRequest(http.MethodPost, wh.Path, nil)
	w := httptest.NewRecorder()

	// Act
	handler(w, req)

	// Assert
	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != expected {
		t.Errorf("expected %s, got %s", expected, string(data))
	}
}

func TestApi_WebhookHandler_OnlyPost(t *testing.T) {
	// Arrange
	expectedStatus := http.StatusMethodNotAllowed
	whName, wh := newTestWebhook()
	handler := newTestApi().createHandler(whName, wh)
	methods := []string{http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodPatch}

	for _, method := range methods {
		// Act
		req := httptest.NewRequest(method, wh.Path, nil)
		w := httptest.NewRecorder()
		handler(w, req)

		// Assert
		res := w.Result()
		if res.StatusCode != expectedStatus {
			t.Errorf("expected status %d, got %d", expectedStatus, res.StatusCode)
		}
	}
}

func TestApi_WebhookHandler_ValidatorError(t *testing.T) {
	// Arrange
	expectedStatus := http.StatusUnauthorized
	expectedMessage := "Unauthenticated"
	whName, wh := newTestWebhook()
	wh.Validators = config.Validators{
		Fake: &config.FakeValidator{
			ModeSuccess: false,
		},
	}
	handler := newTestApi().createHandler(whName, wh)

	// Act
	req := httptest.NewRequest(http.MethodPost, wh.Path, nil)
	w := httptest.NewRecorder()
	handler(w, req)

	// Assert
	res := w.Result()
	if res.StatusCode != expectedStatus {
		t.Errorf("expected status %d, got %d", expectedStatus, res.StatusCode)
	}
	defer func() {
		_ = res.Body.Close()
	}()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != expectedMessage {
		t.Errorf("expected messag %s, got %s", expectedMessage, string(data))
	}
}

func TestApi_WebhookHandler_ServerError(t *testing.T) {
	// Arrange
	expectedStatus := http.StatusInternalServerError
	whName, wh := newTestWebhook()
	wh.Tasks = []string{"fake"}
	api := newTestApi()
	api.repo.(*runner_repo.FakeRepo).ModeSuccess = false
	handler := api.createHandler(whName, wh)

	// Act
	req := httptest.NewRequest(http.MethodPost, wh.Path, nil)
	w := httptest.NewRecorder()
	handler(w, req)

	// Assert
	res := w.Result()
	if res.StatusCode != expectedStatus {
		t.Errorf("expected status %d, got %d", expectedStatus, res.StatusCode)
	}
}

func TestApi_DefaultHandler(t *testing.T) {
	// Arrange
	expectedStatus := http.StatusNotFound
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()

	// Act
	defaultHandler(w, req)

	// Assert
	res := w.Result()
	if res.StatusCode != expectedStatus {
		t.Errorf("expected status %d, got %d", expectedStatus, res.StatusCode)
	}
}
