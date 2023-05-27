package telegram

import "github.com/samber/lo"

type Option interface {
	IsOption()
}

type editMessage struct{}

func (editMessage) IsOption() {}

var Edit = editMessage{}

func isEdit(opts []Option) bool {
	_, foundEdit := lo.Find(opts, func(p Option) bool {
		_, ok := p.(editMessage)
		return ok
	})
	return foundEdit
}
