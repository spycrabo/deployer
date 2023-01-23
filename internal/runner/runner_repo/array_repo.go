package runner_repo

import (
	"errors"
	"github.com/google/uuid"
)

type ArrayRepo struct {
	tasks           map[uuid.UUID]RepoTask
	processingTasks map[uuid.UUID]RepoTask
}

func NewArrayRepo() *ArrayRepo {
	return &ArrayRepo{
		tasks:           map[uuid.UUID]RepoTask{},
		processingTasks: map[uuid.UUID]RepoTask{},
	}
}

func (r ArrayRepo) PushTask(taskName string, vars map[string]string) error {
	r.tasks[uuid.New()] = RepoTask{
		TaskName:  taskName,
		Variables: vars,
	}
	return nil
}

func (r ArrayRepo) GetPendingTasks() (map[uuid.UUID]RepoTask, error) {
	return r.tasks, nil
}

func (r ArrayRepo) GetNextTask() (*uuid.UUID, *RepoTask, error) {
	for k, v := range r.tasks {
		return &k, &v, nil
	}
	return nil, nil, nil
}

func (r ArrayRepo) StartTask(uuid uuid.UUID) error {
	t, ok := r.tasks[uuid]
	if !ok {
		return errors.New("task not found")
	}
	_, ok = r.processingTasks[uuid]
	if ok {
		return errors.New("task already started")
	}
	delete(r.tasks, uuid)
	r.processingTasks[uuid] = t
	return nil
}

func (r ArrayRepo) FinishTask(uuid uuid.UUID) error {
	_, ok := r.processingTasks[uuid]
	if !ok {
		return errors.New("task not found")
	}
	delete(r.processingTasks, uuid)
	return nil
}

func (r ArrayRepo) TotalRunningTasks() (uint, error) {
	return uint(len(r.processingTasks)), nil
}
