package notifier

import "testing"

func TestFakeNotifier_Notify_Success(t *testing.T) {
	// Arrange
	n := &FakeNotifier{
		ModeSuccess: true,
	}

	// Act
	err := n.Notify("")

	// Assert
	if err != nil {
		t.Error("expected no error, got", err)
	}
}

func TestFakeNotifier_Notify_Error(t *testing.T) {
	// Arrange
	n := &FakeNotifier{
		ModeSuccess: false,
	}

	// Act
	err := n.Notify("")

	// Assert
	if err == nil {
		t.Error("expected error, got none")
	}
}
