package listener

import (
	"bc-deployer/internal/config"
	"bc-deployer/internal/listener/validator"
	"bc-deployer/internal/runner/runner_repo"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)

type Api struct {
	port       int
	repo       runner_repo.Repo
	validators validator.Validators
}

func NewApi(port int, repo runner_repo.Repo, validators validator.Validators) *Api {
	api := &Api{
		port:       port,
		repo:       repo,
		validators: validators,
	}
	return api
}

func (a *Api) Run(webhooks map[string]config.Webhook) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", defaultHandler)
	for name, webhook := range webhooks {
		path := webhook.Path
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}

		mux.HandleFunc(path, logHandler(a.createHandler(name, webhook)))
	}
	err := http.ListenAndServe(fmt.Sprintf(":%d", a.port), mux)
	if err != nil {
		return err
	}
	return nil
}

func defaultHandler(w http.ResponseWriter, _ *http.Request) {
	respond(w, http.StatusNotFound, "Not found")
}

func (a *Api) createHandler(name string, webhook config.Webhook) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respond(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		vs := a.validators.GetMany(webhook.Validators)

		for _, v := range vs {
			if !v.Validate(r) {
				respond(w, http.StatusUnauthorized, "Unauthenticated")
				return
			}
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			respond(w, http.StatusBadRequest, "Error reading request body")
			return
		}

		for path, expectedValue := range webhook.When {
			if strings.HasPrefix(path, "json.") {
				jsonPath := strings.TrimPrefix(path, "json.")

				actualValue := gjson.Get(string(body), jsonPath).String()

				if expectedValue != actualValue {
					respond(w, http.StatusOK, fmt.Sprintf("Expected %s to be %s, got %s. Skipping webhook", jsonPath, expectedValue, actualValue))
					return
				}
			}
		}

		vars := make(map[string]string)

		for path, varName := range webhook.Mapenv {
			if strings.HasPrefix(path, "json.") {
				jsonPath := strings.TrimPrefix(path, "json.")

				varValue := gjson.Get(string(body), jsonPath).String()

				vars[varName] = varValue
			}
		}

		for _, t := range webhook.Tasks {
			err := a.repo.PushTask(t, vars)
			if err != nil {
				respond(w, http.StatusInternalServerError, "ServerError")
				return
			}
		}

		respond(w, http.StatusOK, fmt.Sprintf("Webhook %s executed", name))
	}
}

func respond(w http.ResponseWriter, status int, content string) {
	w.WriteHeader(status)
	_, _ = w.Write([]byte(content))
}

func logHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		logStr := fmt.Sprintf("[%s] %s %s", now.Format("2006-01-02 15:04:05"), r.Method, r.URL.Path)
		rec := httptest.NewRecorder()
		fn(rec, r)
		logStr += fmt.Sprintf(" -> %d %s", rec.Code, rec.Body.String())
		_, _ = w.Write(rec.Body.Bytes())
		fmt.Println(logStr)
	}
}

func isJson(r *http.Request) bool {
	h, ok := r.Header["Content-Type"]
	if !ok {
		return false
	}
	for _, v := range h {
		if v == "application/json" {
			return true
		}
	}
	return false
}
