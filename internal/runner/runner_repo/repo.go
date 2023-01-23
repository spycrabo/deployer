package runner_repo

import (
	"bc-deployer/internal/config"
	"errors"
	"github.com/google/uuid"
)

type Repo interface {
	PushTask(taskName string, vars map[string]string) error
	GetPendingTasks() (map[uuid.UUID]RepoTask, error)
	GetNextTask() (*uuid.UUID, *RepoTask, error)
	TotalRunningTasks() (uint, error)
	StartTask(uuid.UUID) error
	FinishTask(uuid.UUID) error
}

type RepoTask struct {
	TaskName  string
	Variables map[string]string
}

func NewRepo(conf config.Runner) (Repo, error) {
	if conf.Dir != nil {
		r, err := NewDirRepo(conf.Dir)
		if err != nil {
			return nil, err
		}
		return r, nil
	}

	return nil, errors.New("no repo specified")
}
