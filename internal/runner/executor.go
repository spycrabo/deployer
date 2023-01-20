package runner

import (
	"bc-deployer/internal/config"
	"bc-deployer/internal/notifier"
	"bc-deployer/internal/runner/runner_repo"
	"errors"
	"github.com/google/uuid"
	"os/exec"
)

type Executor struct {
	repo      runner_repo.Repo
	notifiers notifier.Notifiers
	conf      conf
}

type TaskResult struct {
	Out []byte
	Err error
}

type conf struct {
	Tasks map[string]config.Task
}

func NewExecutor(repo runner_repo.Repo, notifiers notifier.Notifiers, tasks map[string]config.Task) (*Executor, error) {
	if tasks == nil {
		return nil, errors.New("empty config passed")
	}
	executor := &Executor{
		repo:      repo,
		notifiers: notifiers,
		conf: conf{
			Tasks: tasks,
		},
	}
	return executor, nil
}

func (e *Executor) Run() ([]TaskResult, error) {
	outputs := make([]TaskResult, 0)

	for {
		id, task, err := e.repo.GetNextTask()
		if id == nil && task == nil && err == nil {
			return outputs, nil
		}

		if err != nil {
			return outputs, err
		}
		if id == nil {
			return outputs, errors.New("got null task id")
		}
		if task == nil {
			return outputs, errors.New("got null task")
		}

		output, innerErr := e.runSingleTask(*id, *task)
		outputs = append(outputs, TaskResult{
			Out: output,
			Err: innerErr,
		})
	}
}

func (e *Executor) runSingleTask(id uuid.UUID, repoTask runner_repo.RepoTask) ([]byte, error) {
	err := e.repo.StartTask(id)
	if err != nil {
		return nil, err
	}

	taskConfig, ok := e.conf.Tasks[repoTask.TaskName]
	if !ok {
		err := e.repo.FinishTask(id)
		if err != nil {
			return nil, err
		}
		return nil, errors.New("unknown taskConfig passed: " + repoTask.TaskName)
	}

	cmd := exec.Command("/bin/bash", "-c", taskConfig.Command)
	if taskConfig.Dir != nil {
		cmd.Dir = *taskConfig.Dir
	}
	out, err := cmd.Output()
	if err != nil {
		e.notifyOnFailure(repoTask, taskConfig)
		return nil, err
	}

	err = e.repo.FinishTask(id)
	if err != nil {
		return nil, err
	}
	e.notifyOnSuccess(repoTask, taskConfig)
	return out, nil
}

func (e *Executor) notifyOnSuccess(t runner_repo.RepoTask, c config.Task) {
	msg := "Task " + t.TaskName + " finished successfully"
	if c.Notifications.Success != nil {
		msg = *c.Notifications.Success
	}
	e.notifiers.Notify(msg)
}

func (e *Executor) notifyOnFailure(repoTask runner_repo.RepoTask, task config.Task) {
	msg := "Task " + repoTask.TaskName + " failed"
	if task.Notifications.Success != nil {
		msg = *task.Notifications.Success
	}
	e.notifiers.Notify(msg)
}
