package runner_repo

import (
	"bc-deployer/internal/config"
	"github.com/google/uuid"
	"os"
	"testing"
)

func newTestDirRepo(t *testing.T) *DirRepo {
	path := t.TempDir()

	repo, err := NewDirRepo(&config.DirRunner{Path: path})
	if err != nil {
		panic(err)
	}
	return repo
}

func TestDirRepo_PushTask_Success(t *testing.T) {
	// Arrange
	repo := newTestDirRepo(t)
	taskName := "test"
	vars := map[string]string{}

	// Act
	err := repo.PushTask(taskName, vars)

	// Assert
	if err != nil {
		t.Error(err)
	}
	dirEntries, _ := os.ReadDir(repo.conf.Path)
	if len(dirEntries) != 1 {
		t.Errorf("expected 1 task, got %d", len(dirEntries))
	}
}

func TestDirRepo_GetPendingTasks_Success(t *testing.T) {
	// Arrange
	repo := newTestDirRepo(t)
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

func TestDirRepo_StartTask_Success(t *testing.T) {
	// Arrange
	repo := newTestDirRepo(t)
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
	dirEntries, _ := os.ReadDir(repo.conf.Path)
	if len(dirEntries) != 1 {
		t.Errorf("expected 1 task, got %d", len(dirEntries))
	}
}

func TestDirRepo_FinishTask_Success(t *testing.T) {
	// Arrange
	repo := newTestDirRepo(t)
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
	dirEntries, _ := os.ReadDir(repo.conf.Path)
	if len(dirEntries) != 0 {
		t.Errorf("expected 0 task, got %d", len(dirEntries))
	}
}
