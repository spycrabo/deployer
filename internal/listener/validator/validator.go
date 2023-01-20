package validator

import (
	"bc-deployer/internal/config"
	"net/http"
)

type Validator interface {
	Validate(*http.Request) bool
}

type Validators []Validator

func NewValidators(conf config.Validators) Validators {
	validators := make(Validators, 0)
	if conf.Gitlab != nil {
		validators = append(validators, NewGitlabValidator(conf.Gitlab.Token))
	}
	if conf.Fake != nil {
		validators = append(validators, NewFakeValidator(conf.Fake.ModeSuccess))
	}
	return validators
}
