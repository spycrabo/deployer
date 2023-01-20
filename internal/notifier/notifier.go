package notifier

import "bc-deployer/internal/config"

type Notifier interface {
	Notify(string) error
}

type Notifiers []Notifier

func NewNotifiers(conf config.Notifiers) Notifiers {
	notifiers := make([]Notifier, 0)

	if conf.Telegram != nil {
		n := NewTelegramNotifier(conf.Telegram.BotToken, conf.Telegram.ChatId, conf.Telegram.Params)
		notifiers = append(notifiers, n)
	}

	return notifiers
}

func (notifiers Notifiers) Notify(msg string) []error {
	errs := make([]error, 0)

	for _, n := range notifiers {
		err := n.Notify(msg)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
