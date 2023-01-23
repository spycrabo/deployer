package main

import (
	"bc-deployer/internal/config"
	"bc-deployer/internal/listener"
	"bc-deployer/internal/listener/validator"
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

	repo, err := runner_repo.NewRepo(conf.Runner)
	if err != nil {
		output(fmt.Sprintf("Error initializing repo: \n%v", err))
		return
	}

	validators, err := validator.NewValidators(conf.Validators)
	if err != nil {
		output(fmt.Sprintf("Error initializing validators: \n%v", err))
		return
	}

	api := listener.NewApi(conf.Server.Port, repo, validators)

	output(fmt.Sprintf("Starting API on port %d", conf.Server.Port))

	err = api.Run(conf.Webhooks)
	if err != nil {
		output(fmt.Sprintf("Error starting api: \n%v", err))
		return
	}
}

func output(msg string) {
	now := time.Now()
	fmt.Printf("[%s] %s\n", now.Format("2006-01-02 15:04:05"), msg)
}
