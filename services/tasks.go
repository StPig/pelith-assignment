package services

import (
	"fmt"
	"os"
	"pelith-assignment/database"
	"pelith-assignment/models"
	"time"
)

func checkOnboardingTask(user *models.User) error {
	task := models.Task{}

	param := []interface{}{models.OnboardingType}
	row := database.QuerySingleRow("select id, type, total_points, start_time, end_time from tasks where type = $1", param)
	err := row.Scan(&task.ID, &task.Type, &task.TotalPoint, &task.StartTime, &task.EndTime)
	if err != nil {
		return err
	}

	now := time.Now()
	if now.After(task.StartTime) && now.Before(task.EndTime) {
		totalSwap, err := GetTotalSwap(user.Address, task.StartTime, task.EndTime)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Get total swap failed: %v\n", err)
			return err
		}

		// 1000u 之後可以移去 db 或 config
		if totalSwap >= 1000 {
			err = distributePoints(user, &task, float64(task.TotalPoint))
			if err != nil {
				return err
			}

			param := []interface{}{user.Address}
			err = database.Exec("update users set onboarding_completed = true where address = $1", param)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Update user info failed: %v\n", err)
				return err
			}
		}
	}

	return nil
}

func checkSharePoolTask(user *models.User) error {
	// 須完成 onboarding task
	if !user.OnboardingCompleted {
		return nil
	}

	tasks := []models.Task{}

	param := []interface{}{user.ID, models.SharePoolType}
	rows, err := database.QueryRow("select t.id, type, case when earn_points is not null then true else false end, total_points, start_time, end_time from tasks t left join complete_history ch on t.id = ch.task_id and ch.user_id = $1 where type = $2", param)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Type, &task.IsCompleted, &task.TotalPoint, &task.StartTime, &task.EndTime)
		if err != nil {
			return err
		}
		tasks = append(tasks, task)
	}

	now := time.Now()
	for _, task := range tasks {
		if !task.IsCompleted && now.After(task.EndTime) {
			poolTotalSwap, err := GetTotalSwap("", task.StartTime, task.EndTime)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Get total swap failed: %v\n", err)
				return err
			}

			userTotalSwap, err := GetTotalSwap(user.Address, task.StartTime, task.EndTime)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Get total swap failed: %v\n", err)
				return err
			}

			if userTotalSwap > 0 {
				earnPoints := userTotalSwap / poolTotalSwap * float64(task.TotalPoint)

				err = distributePoints(user, &task, earnPoints)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func GetUserTaskStatus(address string) (models.User, error) {
	user, err := getUserDetail(address)
	if err != nil {
		return models.User{}, err
	}

	if !user.OnboardingCompleted {
		err = checkOnboardingTask(&user)
		if err != nil {
			return models.User{}, err
		}
	}

	if user.OnboardingCompleted {
		err := checkSharePoolTask(&user)
		if err != nil {
			return models.User{}, err
		}
	}

	param := []interface{}{user.ID}
	rows, err := database.QueryRow("select t.id, type, case when earn_points is not null then true else false end, COALESCE(earn_points, 0), total_points, start_time, end_time from tasks t left join complete_history ch on t.id = ch.task_id and ch.user_id = $1", param)
	if err != nil {
		return models.User{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Type, &task.IsCompleted, &task.EarnPoint, &task.TotalPoint, &task.StartTime, &task.EndTime)
		if err != nil {
			return models.User{}, err
		}
		user.Tasks = append(user.Tasks, task)
	}

	return user, nil
}

func GetUserPointsHistory(address string) (models.User, error) {
	user, err := getUserDetail(address)
	if err != nil {
		return models.User{}, err
	}

	if !user.OnboardingCompleted {
		err := checkOnboardingTask(&user)
		if err != nil {
			return models.User{}, err
		}
	}

	if user.OnboardingCompleted {
		err := checkSharePoolTask(&user)
		if err != nil {
			return models.User{}, err
		}
	}

	param := []interface{}{user.ID}
	rows, err := database.QueryRow("select task_id, type, true, earn_points, total_points, start_time, end_time from complete_history ch left join tasks t on ch.task_id = t.id where user_id = $1", param)
	if err != nil {
		return models.User{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Type, &task.IsCompleted, &task.EarnPoint, &task.TotalPoint, &task.StartTime, &task.EndTime)
		if err != nil {
			return models.User{}, err
		}
		user.Tasks = append(user.Tasks, task)
	}

	return user, nil
}
