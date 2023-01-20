package validator

import (
	"bc-deployer/internal/config"
	"testing"
)

func TestValidators_CreatingFromConfig_Success(t *testing.T) {
	// Arrange
	conf := config.Validators{
		Gitlab: &config.GitlabValidator{
			Token: "test",
		},
		Fake: &config.FakeValidator{
			ModeSuccess: true,
		},
	}

	// Act
	v := NewValidators(conf)

	// Assert
	if len(v) != 2 {
		t.Fatalf("expected 2 validators, got %d", len(v))
	}
}
