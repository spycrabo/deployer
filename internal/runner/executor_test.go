package runner

import (
	"bc-deployer/internal/config"
	"bc-deployer/internal/notifier"
	"bc-deployer/internal/runner/runner_repo"
	"testing"
)

func newTestExecutor() *Executor {
	testRepo := runner_repo.NewArrayRepo()

	executor, err := NewExecutor(testRepo, notifier.Notifiers{}, map[string]config.Task{})
	if err != nil {
		panic(err)
	}
	return executor
}

func TestExecutor_Run_Success(t *testing.T) {
	// Arrange
	executor := newTestExecutor()
	err := executor.repo.PushTask("test", map[string]string{})
	dir := "/tmp"
	successMessage := "success"
	failureMessage := "failure"
	expectedOutput := "test\n"
	executor.conf.Tasks["test"] = config.Task{
		Dir:     &dir,
		Command: "echo \"test\"",
		Notifications: config.Notifications{
			Success: &successMessage,
			Failure: &failureMessage,
		},
	}
	if err != nil {
		panic(err)
	}

	// Act
	results, err := executor.Run()

	// Assert
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Err != nil {
		t.Fatalf("task failed: %s", results[0].Err)
	}
	outString := string(results[0].Out)
	if outString != expectedOutput {
		t.Errorf("expected task output '%s', got '%s'", expectedOutput, outString)
	}
}
