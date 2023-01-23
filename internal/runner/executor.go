package runner

import (
	"bc-deployer/internal/config"
	"bc-deployer/internal/notifier"
	"bc-deployer/internal/runner/runner_repo"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"os/exec"
	"strings"
)

type Executor struct {
	repo        runner_repo.Repo
	notifiers   notifier.Notifiers
	conf        conf
	concurrency uint
}

type TaskResult struct {
	Task *runner_repo.RepoTask
	Out  []byte
	Err  error
}

type conf struct {
	Tasks map[string]config.Task
}

func NewExecutor(repo runner_repo.Repo, notifiers notifier.Notifiers, tasks map[string]config.Task, concurrency uint) (*Executor, error) {
	if tasks == nil {
		return nil, errors.New("empty config passed")
	}
	executor := &Executor{
		repo:      repo,
		notifiers: notifiers,
		conf: conf{
			Tasks: tasks,
		},
		concurrency: concurrency,
	}
	return executor, nil
}

func (e *Executor) Run() ([]TaskResult, error) {
	if e.concurrency > 0 {
		totalRunningTasks, err := e.repo.TotalRunningTasks()
		if err != nil {
			return nil, err
		}

		if totalRunningTasks >= e.concurrency {
			return nil, errors.New("max concurrency reached")
		}
	}

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
			Task: task,
			Out:  output,
			Err:  innerErr,
		})
		if innerErr == nil {
			e.notifyOnSuccess(*task, e.conf.Tasks[task.TaskName])
		} else {
			e.notifyOnFailure(*task, e.conf.Tasks[task.TaskName])
		}
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

	cmd := exec.Command("/bin/bash", "-c", buildCommand(repoTask, taskConfig))
	if taskConfig.Dir != nil {
		cmd.Dir = *taskConfig.Dir
	}
	out, err := cmd.Output()
	if err != nil {
		err2 := e.repo.FinishTask(id)
		if err2 != nil {
			return nil, errors.New(fmt.Sprintf("%v; %v", err, err2))
		}
		return nil, err
	}

	err = e.repo.FinishTask(id)
	if err != nil {
		return nil, err
	}
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
	if task.Notifications.Failure != nil {
		msg = *task.Notifications.Failure
	}
	e.notifiers.Notify(msg)
}

func buildCommand(task runner_repo.RepoTask, conf config.Task) string {
	cmd := ""
	vars := make([]string, 0)
	for k, v := range task.Variables {
		vars = append(vars, k+"=\""+v+"\"")
	}
	cmd += strings.Join(vars, " ")
	if len(cmd) > 0 {
		cmd += " && "
	}
	cmd += conf.Command
	return cmd
}
