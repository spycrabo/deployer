package main

import (
	"bc-deployer/internal/config"
	"bc-deployer/internal/notifier"
	"bc-deployer/internal/runner"
	"bc-deployer/internal/runner/runner_repo"
	"flag"
	"fmt"
	"time"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "config.yml", "Path to config yaml file")
	flag.Parse()

	conf, err := config.Read(configPath)
	if err != nil {
		output(fmt.Sprintf("Error reading config: \n%v", err))
		return
	}

	executor, err := createExecutor(conf)
	if err != nil {
		output(fmt.Sprintf("Error creating executor: \n%v", err))
		return
	}

	results, err := executor.Run()
	if err != nil {
		if err.Error() == "max concurrency reached" {
			output("Max concurrency reached, exiting")
			return
		}
		output(fmt.Sprintf("Error running tasks: \n%v", err))
		return
	}

	for _, res := range results {
		msg := fmt.Sprintf("Task %s", res.Task.TaskName)
		if res.Err != nil {
			msg = fmt.Sprintf("%s failed: \n%v", msg, res.Err)
		} else {
			msg = fmt.Sprintf("%s finished successfully", msg)
		}

		output(msg)
	}
}

func createExecutor(conf *config.Config) (*runner.Executor, error) {
	repo, err := runner_repo.NewRepo(conf.Runner)
	if err != nil {
		return nil, err
	}

	notifiers := notifier.NewNotifiers(conf.Notifiers)

	return runner.NewExecutor(repo, notifiers, conf.Tasks, conf.Runner.Concurrency)
}

func output(msg string) {
	now := time.Now()
	fmt.Printf("[%s] %s\n", now.Format("2006-01-02 15:04:05"), msg)
}
