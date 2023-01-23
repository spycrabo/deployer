package runner

import (
	"bc-deployer/internal/config"
	"bc-deployer/internal/notifier"
	"bc-deployer/internal/runner/runner_repo"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func newTestExecutor() *Executor {
	testRepo := runner_repo.NewArrayRepo()

	executor, err := NewExecutor(testRepo, notifier.Notifiers{}, map[string]config.Task{}, 1)
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

func TestExecutor_Run_WithEnv(t *testing.T) {
	// Arrange
	expectedOutput := "test_value\n"
	executor := newTestExecutor()
	err := executor.repo.PushTask("test", map[string]string{
		"TEST_VAR": "test_value",
	})
	successMessage := "success"
	failureMessage := "failure"
	executor.conf.Tasks["test"] = config.Task{
		Command: "echo \"$TEST_VAR\"",
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

func TestExecutor_Run_ConcurrencyLimit(t *testing.T) {
	// Arrange
	executor := newTestExecutor()
	cnt := 2
	for i := 0; i < cnt; i++ {
		err := executor.repo.PushTask("test", map[string]string{})
		if err != nil {
			t.Fatal(err)
		}
	}
	executor.conf.Tasks["test"] = config.Task{
		Command: "sleep 1",
	}

	// Act
	ch := make(chan error)
	errs := make([]error, cnt)

	var wg sync.WaitGroup
	for i := 0; i < cnt; i++ {
		wg.Add(1)
		go func(i int) {
			errs[i] = <-ch
		}(i)
		go func(executor *Executor, ch chan error, wg *sync.WaitGroup) {
			defer wg.Done()
			_, err := executor.Run()
			ch <- err
		}(executor, ch, &wg)
	}

	wg.Wait()
	close(ch)

	// Assert
	assert.Nil(t, errs[0])
	assert.Error(t, errs[1], "max concurrency reached")
}
