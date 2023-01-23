package validator

import "net/http"

type HeaderValidator struct {
	header string
	value  string
}

func NewHeaderValidator(header string, value string) *HeaderValidator {
	return &HeaderValidator{
		header: header,
		value:  value,
	}
}

func (h HeaderValidator) Validate(r *http.Request) bool {
	return r.Header.Get(h.header) == h.value
}
