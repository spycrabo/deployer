package runner_repo

import (
	"errors"
	"github.com/google/uuid"
)

type FakeRepo struct {
	ModeSuccess bool
}

func NewFakeRepo(modeSuccess bool) *FakeRepo {
	return &FakeRepo{
		modeSuccess,
	}
}

func (f FakeRepo) PushTask(taskName string, vars map[string]string) error {
	if f.ModeSuccess {
		return nil
	}
	return errors.New("fake error")
}

func (f FakeRepo) GetPendingTasks() (map[uuid.UUID]RepoTask, error) {
	if f.ModeSuccess {
		return nil, nil
	}
	return nil, errors.New("fake error")
}

func (f FakeRepo) GetNextTask() (*uuid.UUID, *RepoTask, error) {
	if f.ModeSuccess {
		return nil, nil, nil
	}
	return nil, nil, errors.New("fake error")
}

func (f FakeRepo) StartTask(uuid uuid.UUID) error {
	if f.ModeSuccess {
		return nil
	}
	return errors.New("fake error")
}

func (f FakeRepo) FinishTask(uuid uuid.UUID) error {
	if f.ModeSuccess {
		return nil
	}
	return errors.New("fake error")
}
