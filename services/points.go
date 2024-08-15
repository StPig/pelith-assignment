package services

import (
	"fmt"
	"os"
	"pelith-assignment/database"
	"pelith-assignment/models"
)

func distributePoints(user *models.User, task *models.Task, points float64) error {
	trans, err := database.GetTx()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Distribute points failed: %v\n", err)
		return err
	}

	param := []interface{}{points, user.Address}
	err = database.TxExec(trans, "update users set points = points + $1 where address = $2", param)
	if err != nil {
		database.TxRollBack(trans)
		fmt.Fprintf(os.Stderr, "Update user points failed: %v\n", err)
		return err
	}

	param = []interface{}{user.ID, task.ID, points}
	err = database.TxExec(trans, "insert into complete_history (user_id, task_id, earn_points, create_time) values ($1, $2, $3, now())", param)
	if err != nil {
		database.TxRollBack(trans)
		fmt.Fprintf(os.Stderr, "Insert points history failed: %v\n", err)
		return err
	}

	user.OnboardingCompleted = true
	user.Points += int(points)

	database.TxCommit(trans)

	return nil
}

func GetPointsLeaderBoard() ([]models.User, error) {
	leaderBoard := []models.User{}

	rows, err := database.QueryRow("select id, address, points from users order by points desc limit 100", []interface{}{})
	if err != nil {
		return leaderBoard, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Address, &user.Points)
		if err != nil {
			return leaderBoard, err
		}
		leaderBoard = append(leaderBoard, user)
	}

	return leaderBoard, nil
}
