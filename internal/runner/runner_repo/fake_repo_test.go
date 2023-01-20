package runner_repo

import (
	"github.com/google/uuid"
	"testing"
)

func TestFakeRepo_StartTask_Success(t *testing.T) {
	runSuccessTest(t, func(r *FakeRepo) error {
		return r.StartTask(uuid.New())
	})
}

func TestFakeRepo_StartTask_Error(t *testing.T) {
	runFailureTest(t, func(r *FakeRepo) error {
		return r.StartTask(uuid.New())
	})
}

func TestFakeRepo_GetPendingTasks_Success(t *testing.T) {
	runSuccessTest(t, func(r *FakeRepo) error {
		_, err := r.GetPendingTasks()
		return err
	})
}

func TestFakeRepo_GetPendingTasks_Error(t *testing.T) {
	runFailureTest(t, func(r *FakeRepo) error {
		_, err := r.GetPendingTasks()
		return err
	})
}

func TestFakeRepo_PushTask_Success(t *testing.T) {
	runSuccessTest(t, func(r *FakeRepo) error {
		return r.PushTask("test", nil)
	})
}

func TestFakeRepo_PushTask_Error(t *testing.T) {
	runFailureTest(t, func(r *FakeRepo) error {
		return r.PushTask("test", nil)
	})
}

func TestFakeRepo_FinishTask_Success(t *testing.T) {
	runSuccessTest(t, func(r *FakeRepo) error {
		return r.FinishTask(uuid.New())
	})
}

func TestFakeRepo_FinishTask_Error(t *testing.T) {
	runFailureTest(t, func(r *FakeRepo) error {
		return r.FinishTask(uuid.New())
	})
}

func runSuccessTest(t *testing.T, f func(r *FakeRepo) error) {
	// Arrange
	r := NewFakeRepo(true)

	// Act
	err := f(r)

	// Assert
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func runFailureTest(t *testing.T, f func(r *FakeRepo) error) {
	// Arrange
	r := NewFakeRepo(false)

	// Act
	err := f(r)

	// Assert
	if err == nil {
		t.Errorf("expected error, got success")
	}
	if err != nil && err.Error() != "fake error" {
		t.Errorf("expected fake error, got %v", err)
	}
}
