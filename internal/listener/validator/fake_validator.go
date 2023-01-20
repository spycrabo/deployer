package validator

import "net/http"

type FakeValidator struct {
	ModeSuccess bool
}

func NewFakeValidator(modeSuccess bool) *FakeValidator {
	return &FakeValidator{
		ModeSuccess: modeSuccess,
	}
}

func (f *FakeValidator) Validate(r *http.Request) bool {
	return f.ModeSuccess
}
