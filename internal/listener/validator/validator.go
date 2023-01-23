package validator

import (
	"bc-deployer/internal/config"
	"errors"
	"net/http"
)

type Validator interface {
	Validate(*http.Request) bool
}

type Validators map[string]Validator

func NewValidators(config config.ValidatorsConfig) (Validators, error) {
	validators := make(Validators, 0)

	for name, conf := range config {
		switch conf.Driver {
		case "header":
			header, ok := conf.Options["header"]
			if !ok {
				return nil, errors.New("header validator: no header specified")
			}
			value, ok := conf.Options["value"]
			if !ok {
				return nil, errors.New("header validator: no value specified")
			}
			validators[name] = NewHeaderValidator(header, value)
		case "fake":
			validators[name] = NewFakeValidator(conf.Options["success"] == "true")
		}
		// todo add more validators: json, query
	}

	return validators, nil
}

func (v Validators) Get(name string) *Validator {
	validator, ok := v[name]
	if !ok {
		return nil
	}
	return &validator
}

func (v Validators) GetMany(names []string) []Validator {
	validators := make([]Validator, 0)
	for _, name := range names {
		validator, ok := v[name]
		if !ok {
			continue
		}
		validators = append(validators, validator)
	}
	return validators
}
