package models

type Registration struct {
	RegistrationId      int    `json:"registrationId,omitempty"`
	ActivityId          int    `json:"activityId"`
	Comments            string `json:"comments"`
	UserId              int    `json:"userId"`
	TeamId              int    `json:"teamId"`
	GroupId             int    `json:"groupId"`
	IsSwipeRegistration bool   `json:"isSwipeRegistration"`
	ShowComments        bool   `json:"showComments"`
}
