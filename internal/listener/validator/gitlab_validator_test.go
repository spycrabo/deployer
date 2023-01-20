package validator

import (
	"net/http"
	"testing"
)

func TestGitlabValidator_Success(t *testing.T) {
	// Arrange
	secret := "secret"
	validator := NewGitlabValidator(secret)
	r := &http.Request{
		Header: http.Header{
			"X-Gitlab-Token": []string{secret},
		},
	}

	// Act
	ok := validator.Validate(r)

	// Assert
	if !ok {
		t.Error("valid token returned false")
	}
}

func TestGitlabValidator_WrongTokenError(t *testing.T) {
	// Arrange
	validSecret := "validSecret"
	wrongSecret := "wrongSecret"
	validator := NewGitlabValidator(validSecret)
	r := &http.Request{
		Header: http.Header{
			"X-Gitlab-Token": []string{wrongSecret},
		},
	}

	// Act
	ok := validator.Validate(r)

	// Assert
	if ok {
		t.Error("wrong token returned true")
	}
}

func TestGitlabValidator_NoTokenError(t *testing.T) {
	// Arrange
	validSecret := "validSecret"
	validator := NewGitlabValidator(validSecret)
	r := &http.Request{}

	// Act
	ok := validator.Validate(r)

	// Assert
	if ok {
		t.Error("no token returned true")
	}
}
