package validator

import "testing"

func TestFakeValidator_Success(t *testing.T) {
	// Arrange
	v := NewFakeValidator(true)

	// Act
	result := v.Validate(nil)

	// Assert
	if !result {
		t.Error("expected true, got false")
	}
}

func TestFakeValidator_Failure(t *testing.T) {
	// Arrange
	v := NewFakeValidator(false)

	// Act
	result := v.Validate(nil)

	// Assert
	if result {
		t.Error("expected false, got true")
	}
}
