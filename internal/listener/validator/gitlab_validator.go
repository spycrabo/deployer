package validator

import "net/http"

type GitlabValidator struct {
	secret string
}

func NewGitlabValidator(secret string) *GitlabValidator {
	return &GitlabValidator{secret: secret}
}

func (g *GitlabValidator) Validate(r *http.Request) bool {
	return r.Header.Get("X-Gitlab-Token") == g.secret
}
