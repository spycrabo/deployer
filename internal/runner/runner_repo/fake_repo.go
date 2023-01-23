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

func (f FakeRepo) PushTask(_ string, _ map[string]string) error {
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

func (f FakeRepo) StartTask(_ uuid.UUID) error {
	if f.ModeSuccess {
		return nil
	}
	return errors.New("fake error")
}

func (f FakeRepo) FinishTask(_ uuid.UUID) error {
	if f.ModeSuccess {
		return nil
	}
	return errors.New("fake error")
}

func (f FakeRepo) TotalRunningTasks() (uint, error) {
	if f.ModeSuccess {
		return 0, nil
	}
	return 0, errors.New("fake error")
}
