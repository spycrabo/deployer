package notifier

import (
	"errors"
)

type FakeNotifier struct {
	ModeSuccess bool
}

func (n *FakeNotifier) Notify(msg string) error {
	if n.ModeSuccess {
		return nil
	}
	return errors.New("fake error")
}
