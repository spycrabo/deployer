package runner_repo

import (
	"bc-deployer/internal/config"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type DirRepo struct {
	conf config.DirRunner
}

const taskPrefix = ".task-"
const runningTaskSuffix = ".running"

func NewDirRepo(conf *config.DirRunner) (*DirRepo, error) {
	if conf == nil {
		return nil, errors.New("empty config passed")
	}
	runner := &DirRepo{
		conf: *conf,
	}
	return runner, nil
}

func (r *DirRepo) PushTask(taskName string, vars map[string]string) error {
	taskFileName := generateNewTaskFileName()

	task := RepoTask{
		TaskName:  taskName,
		Variables: vars,
	}

	yml, err := yaml.Marshal(task)
	if err != nil {
		return err
	}

	dir := strings.TrimSuffix(r.conf.Path, "/")
	path := fmt.Sprintf("%s/%s", dir, taskFileName)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	_, err = file.Write(yml)
	if err != nil {
		return err
	}
	_ = file.Close()
	return nil
}

func (r *DirRepo) GetPendingTasks() (map[uuid.UUID]RepoTask, error) {
	files, err := os.ReadDir(r.conf.Path)
	if err != nil {
		return nil, err
	}

	tasks := make(map[uuid.UUID]RepoTask)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.HasSuffix(file.Name(), runningTaskSuffix) {
			continue
		}
		if strings.HasPrefix(file.Name(), taskPrefix) {
			path := fmt.Sprintf("%s/%s", strings.TrimSuffix(r.conf.Path, "/"), file.Name())
			taskUuid, _ := uuid.Parse(strings.TrimPrefix(file.Name(), taskPrefix))
			var task RepoTask
			if _, err := os.Stat(path); os.IsNotExist(err) {
				return nil, errors.New("task file does not exist")
			}
			f, err := os.ReadFile(path)
			if err != nil {
				return nil, errors.New("could not read task file")
			}

			err = yaml.Unmarshal(f, &task)
			if err != nil {
				return nil, errors.New("could not unmarshal task file")
			}
			tasks[taskUuid] = task
		}
	}

	return tasks, nil
}

func (r *DirRepo) GetNextTask() (*uuid.UUID, *RepoTask, error) {
	pendingTasks, err := r.GetPendingTasks()
	if err != nil {
		return nil, nil, err
	}

	for id, task := range pendingTasks {
		return &id, &task, nil
	}

	return nil, nil, nil
}

func (r *DirRepo) StartTask(task uuid.UUID) error {
	taskFileName := getTaskFileName(task)
	dir := strings.TrimSuffix(r.conf.Path, "/")
	initialPath := fmt.Sprintf("%s/%s", dir, taskFileName)
	runningPath := fmt.Sprintf("%s.running", initialPath)
	return os.Rename(initialPath, runningPath)
}

func (r *DirRepo) FinishTask(task uuid.UUID) error {
	taskFileName := getTaskFileName(task)
	dir := strings.TrimSuffix(r.conf.Path, "/")
	runningPath := fmt.Sprintf("%s/%s.running", dir, taskFileName)
	return os.Remove(runningPath)
}

func (r *DirRepo) TotalRunningTasks() (uint, error) {
	files, err := os.ReadDir(r.conf.Path)
	if err != nil {
		return 0, err
	}

	var total uint
	total = 0
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.HasSuffix(file.Name(), runningTaskSuffix) {
			total++
		}
	}

	return total, nil
}

func getTaskFileName(taskUuid uuid.UUID) string {
	return taskPrefix + taskUuid.String()
}

func generateNewTaskFileName() string {
	id := uuid.New()
	return taskPrefix + id.String()
}
