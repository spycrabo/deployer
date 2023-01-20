package notifier

import "testing"

func TestNotifiers_Notify_Success(t *testing.T) {
	// Arrange
	notifiers := Notifiers{
		&FakeNotifier{
			ModeSuccess: true,
		},
	}

	// Act
	errs := notifiers.Notify("")

	// Assert
	if len(errs) != 0 {
		t.Error("expected 0 errors, got", len(errs))
	}
}

func TestNotifiers_Notify_Error(t *testing.T) {
	// Arrange
	notifiers := Notifiers{
		&FakeNotifier{
			ModeSuccess: false,
		},
		&FakeNotifier{
			ModeSuccess: true,
		},
	}

	// Act
	errs := notifiers.Notify("")

	// Assert
	if len(errs) != 1 {
		t.Error("expected 1 error, got", len(errs))
	}
}
