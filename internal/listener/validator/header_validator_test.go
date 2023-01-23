package validator

import (
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestHeaderValidator_Success(t *testing.T) {
	// Arrange
	validator := NewHeaderValidator("X-Test", "test")
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Test", "test")

	// Act
	result := validator.Validate(req)

	// Assert
	assert.True(t, result)
}

func TestHeaderValidator_FailureByValue(t *testing.T) {
	// Arrange
	validator := NewHeaderValidator("X-Test", "test")
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Test", "wrong")

	// Act
	result := validator.Validate(req)

	// Assert
	assert.False(t, result)
}

func TestHeaderValidator_FailureByHeader(t *testing.T) {
	// Arrange
	validator := NewHeaderValidator("X-Test", "test")
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Not-Test", "test")

	// Act
	result := validator.Validate(req)

	// Assert
	assert.False(t, result)
}
