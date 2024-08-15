package models

type User struct {
	ID                  int    `json:"id"`
	Address             string `json:"address"`
	OnboardingCompleted bool   `json:"onboardingCompleted"`
	Points              int    `json:"points"`
	Tasks               []Task `json:"tasks"`
}
