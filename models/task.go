package models

import "time"

const (
	OnboardingType = iota
	SharePoolType
)

type Task struct {
	ID          int       `json:"id"`
	Type        int       `json:"type"`
	IsCompleted bool      `json:"isCompleted"`
	EarnPoint   int       `json:"earnPoint"`
	TotalPoint  int       `json:"totalPoint"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
}
