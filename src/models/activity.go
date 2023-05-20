package models

import "time"

type Activity struct {
	ActivityId       int       `json:"activityId"`
	ActivityLocation string    `json:"activityLocation"`
	Start            time.Time `json:"start"`
	Finish           time.Time `json:"finish"`

	Registrations       int `json:"registration"`
	NeededRegistrations int `json:"neededRegistrations"`

	Name string `json:"name"`
}
