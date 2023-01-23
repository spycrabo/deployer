package validator

import (
	"bc-deployer/internal/config"
	"testing"
)

func TestValidators_CreatingFromConfig_Success(t *testing.T) {
	// Arrange
	conf := config.ValidatorsConfig{
		"test-1": {
			Driver: "header",
			Options: map[string]string{
				"header": "Token",
				"value":  "xxx",
			},
		},
		"test-2": {
			Driver: "fake",
			Options: map[string]string{
				"success": "true",
			},
		},
	}

	// Act
	v, err := NewValidators(conf)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(v) != 2 {
		t.Fatalf("expected 2 validators, got %d", len(v))
	}
}
