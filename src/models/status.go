package models

import "time"

type Status struct {
	TransactionsAmount float64 `json:"transactionsAmount"`
	Target             float64 `json:"target"`
	PercentageValue    float64 `json:"percentageValue"`
	Currency           string  `json:"currency"`
}

type ClubMilestone struct {
	Id             int     `json:"id"`
	MilestoneValue float64 `json:"milestoneValue"`
	Date           time.Time
}
