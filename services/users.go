package services

import (
	"database/sql"
	"pelith-assignment/database"
	"pelith-assignment/models"
)

func CreateUser(address string) error {
	param := []interface{}{address}
	err := database.Exec("insert into users (address, create_time) values ($1, now()) on conflict (address) do nothing", param)
	if err != nil {
		return err
	}

	return nil
}

func getUserDetail(address string) (models.User, error) {
	user := models.User{}

	param := []interface{}{address}
	row := database.QuerySingleRow("select id, address, onboarding_completed, points from users where address = $1", param)
	err := row.Scan(&user.ID, &user.Address, &user.OnboardingCompleted, &user.Points)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, nil
		}
		return models.User{}, err
	}

	return user, nil
}
