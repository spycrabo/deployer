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
	}

	executor, err := createExecutor(conf)
	if err != nil {
		output(fmt.Sprintf("Error creating executor: \n%v", err))
	}

	_, err = executor.Run()
	if err != nil {
		output(fmt.Sprintf("Error running tasks: \n%v", err))
	}
}

func createExecutor(conf *config.Config) (*runner.Executor, error) {
	repo, err := runner_repo.NewRepo(conf)
	if err != nil {
		return nil, err
	}

	notifiers := notifier.NewNotifiers(conf.Notifiers)

	return runner.NewExecutor(repo, notifiers, conf.Tasks)
}

func output(msg string) {
	now := time.Now()
	fmt.Printf("[%s] %s\n", now.Format("2006-01-02 15:04:05"), msg)
}
