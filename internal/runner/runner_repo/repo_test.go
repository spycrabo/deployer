package runner_repo

import (
	"bc-deployer/internal/config"
	"testing"
)

func TestRepo_New_DirCreated(t *testing.T) {
	// Arrange
	conf := &config.Config{
		Runner: config.Runner{
			Dir: &config.DirRunner{
				Path: "test",
			},
		},
	}

	// Act
	r, err := NewRepo(conf)

	// Assert
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := r.(*DirRepo); !ok {
		t.Fatal("wrong type")
	}
}

func TestRepo_New_NoRepo(t *testing.T) {
	// Arrange
	conf := &config.Config{
		Runner: config.Runner{},
	}

	// Act
	_, err := NewRepo(conf)

	// Assert
	if err == nil {
		t.Fatal("expected error")
	}
}
