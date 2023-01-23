package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Server     *Server            `yaml:"server"`
	Webhooks   map[string]Webhook `yaml:"webhooks"`
	Validators ValidatorsConfig   `yaml:"validators"`
	Runner     Runner             `yaml:"runner"`
	Tasks      map[string]Task    `yaml:"tasks"`
	Notifiers  Notifiers          `yaml:"notifiers"`
}

type Server struct {
	Port int `yaml:"port"`
}

type Webhook struct {
	Path       string            `yaml:"path"`
	Validators ValidatorsNames   `yaml:"validators"`
	When       map[string]string `yaml:"when"`
	Mapenv     map[string]string `yaml:"mapenv"`
	Tasks      []string          `yaml:"tasks"`
}

type ValidatorsNames []string

type ValidatorsConfig map[string]ValidatorConfig

type ValidatorConfig struct {
	Driver  string                 `yaml:"driver"`
	Options ValidatorOptionsConfig `yaml:"options"`
}

type ValidatorOptionsConfig map[string]string

type FakeValidator struct {
	ModeSuccess bool
}

type Runner struct {
	Concurrency uint       `default:"1" yaml:"concurrency"`
	Dir         *DirRunner `yaml:"dir"`
}

type DirRunner struct {
	Path string `yaml:"path"`
}

type Task struct {
	Dir           *string       `yaml:"dir"`
	Command       string        `yaml:"command"`
	Notifications Notifications `yaml:"notifications"`
}

type Notifications struct {
	Success *string
	Failure *string
}

type Notifiers struct {
	Telegram *TelegramNotifier `yaml:"telegram"`
}

type TelegramNotifier struct {
	BotToken string                 `yaml:"botToken"`
	ChatId   int64                  `yaml:"chatId"`
	Params   map[string]interface{} `yaml:"params"`
}

func Read(path string) (*Config, error) {
	config := &Config{}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("Error: config file does not exist")
		return nil, errors.New("config file does not exist")
	}
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.New("could not read config file")
	}

	err = yaml.Unmarshal(f, config)
	if err != nil {
		return nil, errors.New("could not unmarshal config file")
	}

	return config, nil
}
