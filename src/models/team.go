package models

type Org struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	TeamId int    `json:"teamId"`
}
