package models

import (
	"fmt"
	"time"
)

type Activity struct {
	ActivityId       int         `json:"activityId"`
	ActivityLocation string      `json:"activityLocation"`
	Start            MyshareTime `json:"start"`
	Finish           MyshareTime `json:"finish"`

	Registrations       int `json:"registrations"`
	NeededRegistrations int `json:"neededRegistrations"`

	RegistrationId     int `json:"registrationId"`
	RegistrationStatus int `json:"registrationStatus"`

	Name        string `json:"name"`
	Description string `json:"description"`

	Responsible ActivityResponsible `json:"responsible"`
}

type RegistrationStatus int

const (
	RegistrationRegistered    = 1
	RegistrationNotRegistered = 0
)

type ActivityResponsible struct {
	DisplayName string
}

type MyshareTime time.Time

func (m *MyshareTime) UnmarshalText(data []byte) error {
	fmt.Printf("parse time: %s\n", data)
	t, err := time.Parse(time.RFC3339, string(data))
	if err == nil {
		*m = MyshareTime(t)
		return nil
	}
	data = append(data, 'Z')
	t, err = time.Parse(time.RFC3339, string(data))
	if err == nil {
		*m = MyshareTime(t)
	}

	return err
}

func (m MyshareTime) Format(f string) string {
	return time.Time(m).Format(f)
}

func (m MyshareTime) After(t time.Time) bool {
	return time.Time(m).After(t)
}
