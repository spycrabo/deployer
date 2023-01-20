package runner_repo

import (
	"github.com/google/uuid"
	"testing"
)

func newTestArrayRepo(t *testing.T) *ArrayRepo {
	return NewArrayRepo()
}

func TestArrayRepo_PushTask_Success(t *testing.T) {
	// Arrange
	repo := newTestArrayRepo(t)
	taskName := "test"
	vars := map[string]string{}

	// Act
	err := repo.PushTask(taskName, vars)

	// Assert
	if err != nil {
		t.Error(err)
	}
	if len(repo.tasks) != 1 {
		t.Errorf("expected 1 task, got %d", len(repo.tasks))
	}
}

func TestArrayRepo_GetPendingTasks_Success(t *testing.T) {
	// Arrange
	repo := newTestArrayRepo(t)
	taskName := "test"
	vars := map[string]string{}
	_ = repo.PushTask(taskName, vars)

	// Act
	tasks, err := repo.GetPendingTasks()

	// Assert
	if err != nil {
		t.Error(err)
	}
	if len(tasks) != 1 {
		t.Errorf("expected 1 task, got %d", len(tasks))
	}
}

func TestArrayRepo_StartTask(t *testing.T) {
	// Arrange
	repo := newTestArrayRepo(t)
	taskName := "test"
	vars := map[string]string{}
	_ = repo.PushTask(taskName, vars)
	tasks, _ := repo.GetPendingTasks()
	var taskID uuid.UUID
	for id := range tasks {
		taskID = id
	}

	// Act
	err := repo.StartTask(taskID)

	// Assert
	if err != nil {
		t.Error(err)
	}
	if len(repo.tasks) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(repo.tasks))
	}
	if len(repo.processingTasks) != 1 {
		t.Errorf("expected 1 processing task, got %d", len(repo.processingTasks))
	}
}

func TestArrayRepo_FinishTask(t *testing.T) {
	// Arrange
	repo := newTestArrayRepo(t)
	taskName := "test"
	vars := map[string]string{}
	_ = repo.PushTask(taskName, vars)
	tasks, _ := repo.GetPendingTasks()
	var taskID uuid.UUID
	for id := range tasks {
		taskID = id
	}
	_ = repo.StartTask(taskID)

	// Act
	err := repo.FinishTask(taskID)

	// Assert
	if err != nil {
		t.Error(err)
	}
	if len(repo.tasks) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(repo.tasks))
	}
	if len(repo.processingTasks) != 0 {
		t.Errorf("expected 0 processing tasks, got %d", len(repo.processingTasks))
	}
}
