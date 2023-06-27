package models

import (
	"time"
)

type MyshareActivity struct {
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

type ContributionsActivity struct {
	Id     int `json:"id"`
	TeamId int `json:"teamID"`

	Name        string `json:"name"`
	Description string `json:"description"`

	Start   MyshareTime `json:"start"`
	Finish  MyshareTime `json:"finish"`
	Created MyshareTime `json:"created"`
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

func (m MyshareTime) MarshalText() ([]byte, error) {
	return time.Time(m).MarshalText()
}
func (m MyshareTime) Format(f string) string {
	return time.Time(m).Format(f)
}

func (m MyshareTime) After(t time.Time) bool {
	return time.Time(m).After(t)
}
