package listener

import (
	"bc-deployer/internal/config"
	"bc-deployer/internal/listener/validator"
	"bc-deployer/internal/runner/runner_repo"
	"fmt"
	"net/http"
	"strings"
)

type Api struct {
	port int
	repo runner_repo.Repo
}

func NewApi(port int, repo runner_repo.Repo) *Api {
	api := &Api{
		port: port,
		repo: repo,
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

		mux.HandleFunc(path, a.createHandler(name, webhook))
	}
	err := http.ListenAndServe(fmt.Sprintf(":%d", a.port), mux)
	if err != nil {
		return err
	}
	return nil
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	respond(w, http.StatusNotFound, "Not found")
}

func (a *Api) createHandler(name string, webhook config.Webhook) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respond(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		vs := validator.NewValidators(webhook.Validators)

		for _, v := range vs {
			if !v.Validate(r) {
				respond(w, http.StatusUnauthorized, "Unauthenticated")
				return
			}
		}

		var vars map[string]string

		// todo get from json by path from mapenv

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
