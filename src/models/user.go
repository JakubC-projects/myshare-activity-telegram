package models

import "golang.org/x/oauth2"

type User struct {
	ChatId        int64
	Token         *oauth2.Token
	DisplayName   string
	Org           *Org
	LastMessageId int
}
